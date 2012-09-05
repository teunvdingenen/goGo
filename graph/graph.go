//The Graph package hold all functions of the MCTS algorithm
package graph

import (
    "crypto/rand"
    "log"
    "time"
)

var topVertex *Vertex
var currentVertex *Vertex

var komi float32

var logger *log.Logger

//The Vertex structure holds a boardState and a score up to that point in the game
type Vertex struct {
    whiteWins  uint16
    blackWins  uint16
    whiteScore uint8
    blackScore uint8
    boardState *Board
    turn       uint8

    plyDepth uint
    outEdges []*Edge
    inEdge   *Edge
}

//The Edge structure holds a play and links together two Vertices
type Edge struct {
    playX uint8
    playY uint8

    fromVertex *Vertex
    toVertex   *Vertex
}

//Initiate sets up the first Vertex and links the logfile to this package
func Initiate(boardSize uint8, log *log.Logger) {
    logger = log
    topVertex = new(Vertex)
    topVertex.boardState = new(Board)
    topVertex.boardState.Create(uint16(boardSize))
    topVertex.turn = 1
    topVertex.plyDepth = 0
    currentVertex = topVertex
}

//setKomi informes the graph package of the komi used in the match
func SetKomi(k float32) {
    komi = k
}

//The Reset function removes all vertices up till now so that a new game can be started
func Reset() {
    topVertex.boardState.Create(uint16(topVertex.boardState.size))
    currentVertex = topVertex
}

//GetMove is called when the program should generate a move. It takes either a 1 or a 2
//as argument which stands for the generation of a black or white move
func GetMove(c uint8) (x, y uint8) {
    createGraph()
    var bestMove *Edge = nil
    var biggestDiff uint16 = 0
    if c != currentVertex.turn {
        panic("Trying to getmove for wrong color")
    }
    if c == 1 {
        for _, v := range currentVertex.outEdges {
            if v.toVertex.blackWins > v.toVertex.whiteWins {
                if diff := v.toVertex.blackWins - v.toVertex.whiteWins; diff > biggestDiff {
                    bestMove = v
                    biggestDiff = diff
                }
            }
        }
    } else if c == 2 {
        for _, v := range currentVertex.outEdges {
            if v.toVertex.blackWins < v.toVertex.whiteWins {
                if diff := v.toVertex.whiteWins - v.toVertex.blackWins; diff > biggestDiff {
                    bestMove = v
                    biggestDiff = diff
                }
            }
        }
    } else {
        panic("Unknown color in GetMove(..)")
    }
    if bestMove == nil {
        return 255, 255
        //panic("Unable to get a move from graph")
    }
    currentVertex.outEdges = []*Edge{bestMove}

    board := new(Board)
    board.Create(uint16(currentVertex.boardState.size))
    copy(board.s, currentVertex.boardState.s)
    board.Play(c, bestMove.playX, bestMove.playY)

    currentVertex = bestMove.toVertex
    currentVertex.boardState = board
    logger.Printf("My Board is now:\n")
    logger.Printf(board.tostr())
    return bestMove.playX, bestMove.playY
}

//When a human does a play, the graph needs to be updated to the real state of the board.
//This function updates the currentVertex to that play, or creates it if it does not exist yet
func UpdateCurrentVertex(c, x, y uint8) {
    if c != currentVertex.turn {
        panic("Updating graph to wrong color's move")
    }
    var newCurrent *Vertex = nil
    for _, v := range currentVertex.outEdges {
        if v.playX == x && v.playY == y {
            currentVertex.outEdges = []*Edge{v}
            newCurrent = v.toVertex
            break
        }
    }
    if newCurrent == nil {
        newCurrent = new(Vertex)
        board := new(Board)
        board.Create(uint16(currentVertex.boardState.size))
        copy(board.s, currentVertex.boardState.s)

        newCurrent.boardState = board

        score, _ := newCurrent.boardState.Play(c, x, y)

        if c == 1 {
            newCurrent.blackScore += score
            newCurrent.turn = 2
        } else {
            newCurrent.whiteScore += score
            newCurrent.turn = 1
        }
        newCurrent.plyDepth = currentVertex.plyDepth + 1

        newEdge := new(Edge)
        newEdge.fromVertex = currentVertex
        newEdge.toVertex = newCurrent

        newEdge.playX = x
        newEdge.playY = y

        currentVertex.outEdges = []*Edge{newEdge}
        newCurrent.inEdge = newEdge
    }
    currentVertex = newCurrent
    logger.Printf("My Board is now:\n")
    logger.Printf(newCurrent.boardState.tostr())
}

//The function createGraph is run before a move is selected. 
func createGraph() {
    toDepth := uint(currentVertex.boardState.size * currentVertex.boardState.size) / 2 + currentVertex.plyDepth + 1
    doUntil := time.Now().Add(20 * time.Second)
    for time.Now().Before(doUntil) {
        _ = doRoutine(currentVertex, toDepth)
    }
}

//doRoutine expands the graph downwards to a certain plydepth. This function is called recursively
func doRoutine(fromVertex *Vertex, toDepth uint) bool {
    if toDepth == fromVertex.plyDepth {
        return true
    }

    var x, y uint8
    var score uint8
    var board *Board
    err := "error"

    board = new(Board)
    board.Create(uint16(currentVertex.boardState.size))
    copy(board.s, fromVertex.boardState.s)

    possibilities := len(board.GetEmpty())
    i := 0

    for err != "" {
        x, y = getRandomMove(board)
        for _, v := range fromVertex.outEdges {
            if x == v.playX && y == v.playY {
                _, _ = board.Play(fromVertex.turn, x, y)
                v.toVertex.boardState = board
                return doRoutine(v.toVertex, toDepth)
            }
        }

        score, err = board.Play(fromVertex.turn, x, y)

        if fromVertex.plyDepth > 2 {
            koCompare := fromVertex.inEdge.fromVertex.boardState
            if board.IsEqual(koCompare) {
                err = "KO"
            }
        }
        if err != "" {
            board.Remove(x, y)
            i += 1
        }
        if possibilities == i {
            return false
        }
    }

    newVertex := new(Vertex)

    newEdge := new(Edge)
    newEdge.playX = x
    newEdge.playY = y
    newEdge.fromVertex = fromVertex
    newEdge.toVertex = newVertex

    fromVertex.outEdges = append(fromVertex.outEdges, newEdge)
    newVertex.inEdge = newEdge

    newVertex.boardState = board
    if fromVertex.turn == 1 {
        newVertex.turn = 2
        newVertex.blackScore += score
    } else {
        newVertex.turn = 1
        newVertex.whiteScore += score
    }
    newVertex.plyDepth = fromVertex.plyDepth + 1

    if newVertex.plyDepth > uint(board.size*board.size/2) {
        endBlack, endWhite := scoreBoardOld(newVertex)
        endBlack += newVertex.blackScore
        endWhite += newVertex.whiteScore + uint8(komi-0.5)
        if endBlack > endWhite {
            vertex := newVertex
            for vertex != currentVertex {
                vertex.blackWins += 1
                vertex = vertex.inEdge.fromVertex
            }
        } else if endBlack <= endWhite {
            vertex := newVertex
            for vertex != currentVertex {
                vertex.whiteWins += 1
                vertex = vertex.inEdge.fromVertex
            }
        }
    }
    return doRoutine(newVertex, toDepth)
}

//Get random move implements the crypt/rand library to get a truely random move
func getRandomMove(b *Board) (x, y uint8) {
    empty := b.GetEmpty()
    bs := make([]byte, 8)
    _, _ = rand.Read(bs)
    i:=0
    for _,v := range bs {
        i += int(v)
    }
    i = i % len(empty)
    if i%2 == 0 { //index is x
        x = empty[i]
        y = empty[i+1]
    } else { //index is y
        y = empty[i]
        x = empty[i-1]
    }
    return x, y
}

//scoreBoard attempts to get the score of a game at vertex v. This is only done after
//a certain amount of moves have passed. This function is not called in the current configuration
func scoreBoard(v *Vertex) (scoreBlack, scoreWhite uint8) {
    scoreBlack = 0
    scoreWhite = 0
    processedMatrix := make([]uint8, v.boardState.size*v.boardState.size)
    b := v.boardState

    for i, v := range processedMatrix {
        if v == 1 {
            continue
        }
        x, y := calcXY(i, b.size)
        if b.GetColor(x, y) == 0 {
            xa, xb, ya, yb := GetAdjecent(x, y)
            adjecentColor := uint8(0)
            found2Colors := false
            if xa < b.size && b.GetColor(xa, y) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(xa, y)
                } else {
                    found2Colors = true
                }
            }
            if xb < b.size && b.GetColor(xb, y) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(xb, y)
                } else {
                    found2Colors = true
                }
            }
            if ya < b.size && b.GetColor(x, ya) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(x, ya)
                } else {
                    found2Colors = true
                }
            }
            if yb < b.size && b.GetColor(x, yb) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(x, yb)
                } else {
                    found2Colors = true
                }
            }
            if !found2Colors && adjecentColor == 1 {
                scoreBlack += 1
            } else if !found2Colors && adjecentColor == 2 {
                scoreWhite += 1
            }
        }
    }
    logger.Printf("Scored Board (b,w): (%d,%d)\n", scoreBlack, scoreWhite)
    logger.Printf(b.tostr())
    return scoreBlack, scoreWhite
}

//Translates an index to the boardcoordinates it stands for
func calcXY(index int, size uint8) (x, y uint8) {
    i := uint8(index)
    y = i / size
    x = i % size
    return x, y
}

//This function scores a board, only counting the amount of empty fields that are
//surrounding by a single color.
func scoreBoardOld(v *Vertex) (scoreBlack, scoreWhite uint8) {
    b := v.boardState
    empty := b.GetEmpty()

    for len(empty) > 0 {
        xs := []uint8{empty[0]}
        ys := []uint8{empty[1]}
        b.GetGroup(0, 0, xs, ys)

        var adjecentColor uint8 = 0
        found2Colors := false
        for i, x := range xs {
            xa, xb, ya, yb := GetAdjecent(x, ys[i])
            if xa < b.size && b.GetColor(xa, ys[i]) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(xa, ys[i])
                } else {
                    found2Colors = true
                    break
                }
            }
            if xb < b.size && b.GetColor(xb, ys[i]) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(xb, ys[i])
                } else {
                    found2Colors = true
                    break
                }
            }
            if ya < b.size && b.GetColor(x, ya) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(x, ya)
                } else {
                    found2Colors = true
                    break
                }
            }
            if yb < b.size && b.GetColor(x, yb) != adjecentColor {
                if adjecentColor == 0 {
                    adjecentColor = b.GetColor(x, yb)
                } else {
                    found2Colors = true
                    break
                }
            }
        }
        if !found2Colors && adjecentColor == 1 {
            scoreBlack += uint8(len(xs))
        } else if !found2Colors && adjecentColor == 2 {
            scoreWhite += uint8(len(xs))
        }

        for i := 0; i < len(empty); i += 2 {
            emptyx := empty[i]
            emptyy := empty[i+1]
            if IsPresentinGroup(emptyx, emptyy, xs, ys) {
                empty = append(empty[:i], empty[i+2:]...)
            }
        }
    }
    logger.Printf("Scored Board (b,w): (%d,%d)\n", scoreBlack, scoreWhite)
    logger.Printf(b.tostr())

    return scoreBlack, scoreWhite
}
