package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/notnil/chess"
	"gopkg.in/freeeve/uci.v1"
)

var positions = []string{
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
	"r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
	"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
	"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
}

const depth = 4

var eng *uci.Engine

func main() {
	var err error
	eng, err = uci.NewEngine("./stockfish")
	if err != nil {
		log.Fatal(err)
	}
	eng.SetOptions(uci.Options{
		Hash:   128,
		Ponder: false,
	})

	for _, pos := range positions {
		if ok := checkAgainstStockFish(pos, depth); ok {
			log.Println("PASS")
		}
		fmt.Println("\n\n")
	}
}

func checkAgainstStockFish(pos string, depth int) bool {
	log.Printf("Testing %s", pos)
	g, err := chess.FEN(pos)
	if err != nil {
		log.Fatal(err)
	}
	game := chess.NewGame(g)
	log.Println(game.Position().Board().Draw())
	n, myMoves := countMoves(game.Position(), depth)
	log.Printf("Got %d nodes at depth %d", n, depth)

	n2, engineMoves := engineCount(game.Position(), depth)
	log.Printf("Stockfish Got %d nodes at depth %d", n2, depth)

	if n == n2 {
		return true
	}

	log.Println("FAIL")
	foundHardProblem := false
	misCounts := []string{}
	// check for moves I missed
	for move, count := range engineMoves {
		count2, ok := myMoves[move]
		if !ok {
			foundHardProblem = true
			log.Printf("MISSED MOVE %s", move)
		} else if count != count2 {
			misCounts = append(misCounts, move)
		}
	}
	for move := range myMoves {
		if _, ok := engineMoves[move]; !ok {
			log.Printf("INVALID MOVE %s", move)
			foundHardProblem = true
		}
	}
	if foundHardProblem {
		return false
	}
	sort.Strings(misCounts)
	firstMisCount := misCounts[0]
	log.Printf("Count mismatch under %s. Me: %d - Engine: %d", firstMisCount, myMoves[firstMisCount], engineMoves[firstMisCount])
	log.Printf("Going Deeper")
	movObj, err := chess.LongAlgebraicNotation{}.Decode(game.Position(), firstMisCount)
	if err != nil {
		log.Fatal(err)
	}
	newPos := game.Position().Update(movObj)
	return checkAgainstStockFish(newPos.String(), depth-1)
}

func engineCount(pos *chess.Position, depth int) (int, map[string]int) {
	err := eng.SetFEN(pos.String())
	if err != nil {
		log.Fatal(err)
	}
	total, moves, err := eng.Perft(depth)
	if err != nil {
		log.Fatal(err)
	}
	return total, moves
}

func countMoves(pos *chess.Position, depth int) (int, map[string]int) {
	moves := pos.ValidMoves()
	ret := map[string]int{}
	total := 0
	for _, m := range moves {
		n := 1
		if depth > 1 {
			n, _ = countMoves(pos.Update(m), depth-1)
		}
		total += n
		ret[m.String()] = n
	}
	return total, ret
}

// u64 Perft(int depth)
// {
//     MOVE move_list[256];
//     int n_moves, i;
//     u64 nodes = 0;

//     n_moves = GenerateLegalMoves(move_list);
//     if (depth == 1) return n_moves;
//     for (i = 0; i < n_moves; i++) {
//         MakeMove(move_list[i]);
//         nodes += Perft(depth - 1);
//         UndoMove(move_list[i]);
//     }
//     return nodes;
// }
