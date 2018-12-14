Perftwalk is a little tool I made to find and diagnose move generation issues in the [github.com/notnill/chess](https://github.com/notnill/chess) package.

It compares perft results agains stockfish's results. When the total counts differ in a position, we know one of three things:

1. Our engine returns (illegal) moves that stockfish does not.
2. Our engine misses some moves that stockfish returns.
3. At least one move has a different node count than stockfish reports.

In the case of 1 or 2, we have found a reproducable bug in a single position:

```
2018/12/14 01:09:43 Testing rnbq1k1r/pp1Pbppp/2p5/8/2B5/P7/1PP1N1PP/RNBQK2n w KQ - 0 9
2018/12/14 01:09:43
 A B C D E F G H
8♜ ♞ ♝ ♛ - ♚ - ♜
7♟ ♟ - ♙ ♝ ♟ ♟ ♟
6- - ♟ - - - - -
5- - - - - - - -
4- - ♗ - - - - -
3♙ - - - - - - -
2- ♙ ♙ - ♘ - ♙ ♙
1♖ ♘ ♗ ♕ ♔ - - ♞

2018/12/14 01:09:43 Got 1226 nodes at depth 2
2018/12/14 01:09:43 Stockfish Got 1196 nodes at depth 2
2018/12/14 01:09:43 FAIL
2018/12/14 01:09:43 INVALID MOVE e1g1
```

In case 3, we know what path the inconsistency lies, so we can recurse until we find it.

For each test position, this program will either decide there are no problems, or drill down until it finds a specific inaccuracy in a leaf position.

## Running

`go run main.go`

The program expects to find a binary named `stockfish` in the working direcory (alongside main.go).

It uses a vendored copy of `gopkg.in/freeeve/uci.v1` that I altered to parse the perft command. It does not vendor anything else.

If you get it all passing and want to increase depth, there is a single variable to change.