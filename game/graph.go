package main

import (
	"math/rand"
)

var topVertex *Vertex

var currentVertex *Vertex
//TODO make channel, shits not working
var expandGraph bool = false

type Vertex struct {
	whiteWins  int
	blackWins  int
	whiteScore uint8
	blackScore uint8
	boardState Board
	turn       uint8

	plyDepth int
	outEdges []*Edge
	inEdge   *Edge
}

type Edge struct {
	playX uint8
	playY uint8

	fromVertex *Vertex
	toVertex   *Vertex
}

func Start() {
    if expandGraph {
        Reset()
        expandGraph = true
    } else {
        expandGraph = true
    }
	createGraph()
}

func Stop() {
	expandGraph = false
}

func Initiate(boardSize uint8) {
	topVertex = new(Vertex)
	topVertex.boardState.Create(uint16(boardSize))
	topVertex.turn = 1
	topVertex.plyDepth = 0
	currentVertex = topVertex
}

func Reset() {
    Stop()
	topVertex = nil
	currentVertex = nil
	expandGraph = false
}

func GetMove(c uint8) (x, y uint8) {
	var bestMove *Edge = nil
	var mostWins int = 0
	if c != currentVertex.turn {
		panic("Trying to getmove for wrong color")
	}
	if c == 1 {
		for _, v := range currentVertex.outEdges {
			if v.toVertex.blackWins > mostWins {
				bestMove = v
				mostWins = v.toVertex.blackWins
			}
		}
	} else if c == 2 {
		for _, v := range currentVertex.outEdges {
			if v.toVertex.whiteWins > mostWins {
				bestMove = v
				mostWins = v.toVertex.whiteWins
			}
		}
	} else {
		panic("Unknown color in GetMove(..)")
	}
	if bestMove == nil {
		panic("Unable to get a move from graph")
	}

	currentVertex.outEdges = []*Edge{bestMove}
	currentVertex = bestMove.toVertex
	return bestMove.playX, bestMove.playY
}

func UpdateCurrentVertex(c, x, y uint8) {
	if c != currentVertex.turn {
		panic("Updating graph to wrong color's move")
	}
	var newCurrent *Vertex = nil
	for _, v := range currentVertex.outEdges {
		if v.playX == x && v.playY == y {
			newCurrent = v.toVertex
		}
	}
	if newCurrent == nil {
		Stop()
		newCurrent = new(Vertex)
        newCurrent.boardState.Create(uint16(currentVertex.boardState.size))
		copy(newCurrent.boardState.s, currentVertex.boardState.s)
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
	if !expandGraph {
		//Start()
	}
}

func createGraph() {
	maxMoves := int(currentVertex.boardState.size * currentVertex.boardState.size)
	for expandGraph {
		toDepth := maxMoves/2 + currentVertex.plyDepth + 1
		_ = doRoutine(currentVertex, toDepth)
	}
}

func doRoutine(fromVertex *Vertex, toDepth int) bool {
	if toDepth == fromVertex.plyDepth {
		return true
	}
	var score uint8
	var x, y uint8
	var newBoard Board
	err := "error"
	for err != "" {
		x, y = getRandomMove(fromVertex.boardState)
        newBoard.Create(uint16(currentVertex.boardState.size))
		copy(newBoard.s, fromVertex.boardState.s)
		score, err = newBoard.Play(fromVertex.turn, x, y)

        if fromVertex.inEdge != nil {
		    toCompare := fromVertex.inEdge.fromVertex
		    if newBoard.IsEqual(toCompare.boardState) {
			    err = "KO"
		    }
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

	newVertex.boardState = newBoard
	if fromVertex.turn == 1 {
		newVertex.turn = 2
		newVertex.blackScore += score
	} else {
		newVertex.turn = 1
		newVertex.whiteScore += score
	}
	newVertex.plyDepth = fromVertex.plyDepth + 1

	if newVertex.plyDepth > int(newBoard.size*newBoard.size/2) {
		_, _ = scoreBoard(newVertex)
	}
	//continue to new vertex
	return doRoutine(newVertex, toDepth)
}

func getRandomMove(b Board) (x, y uint8) {
	empty := b.GetEmpty()
	i := rand.Intn(len(empty))
	if i%2 == 0 { //index is x
		x = empty[i]
		y = empty[i+1]
	} else { //index is y
		y = empty[i]
		x = empty[i-1]
	}
	return x, y
}

func scoreBoard(v *Vertex) (scoreBlack, scoreWhite int) {
	b := v.boardState
	empty := b.GetEmpty()

	for len(empty) > 0 {
		xs := []uint8{empty[0]}
		ys := []uint8{empty[1]}
		b.GetGroup(0, 0, xs, ys)

		var adjecentColor uint8 = 0
		found2Colors := false
		for i, x := range xs {
			xa, xb, ya, yb := getAdjecent(x, ys[i])
			if xa < b.size && b.getColor(xa, ys[i]) != adjecentColor {
				if adjecentColor == 0 {
					adjecentColor = b.getColor(xa, ys[i])
				} else {
					found2Colors = true
					break
				}
			}
			if xb < b.size && b.getColor(xb, ys[i]) != adjecentColor {
				if adjecentColor == 0 {
					adjecentColor = b.getColor(xb, ys[i])
				} else {
					found2Colors = true
					break
				}
			}
			if ya < b.size && b.getColor(x, ya) != adjecentColor {
				if adjecentColor == 0 {
					adjecentColor = b.getColor(x, ya)
				} else {
					found2Colors = true
					break
				}
			}
			if yb < b.size && b.getColor(x, yb) != adjecentColor {
				if adjecentColor == 0 {
					adjecentColor = b.getColor(x, yb)
				} else {
					found2Colors = true
					break
				}
			}
		}
		if !found2Colors && adjecentColor == 1 {
			scoreBlack += len(xs)
		} else if !found2Colors && adjecentColor == 2 {
			scoreWhite += len(xs)
		}

		for i := 0; i < len(empty); i += 2 {
			emptyx := empty[i]
			emptyy := empty[i+1]
			if IsPresentinGroup(emptyx, emptyy, xs, ys) {
				empty = append(empty[:i], empty[i+2:]...)
			}
		}
	}
	if scoreBlack > scoreWhite {
		vertex := v
		for vertex != topVertex {
			vertex.blackWins += 1
			vertex = vertex.inEdge.fromVertex
		}
	} else if scoreBlack < scoreWhite {
		vertex := v
		for vertex != topVertex {
			vertex.blackWins += 1
			vertex = vertex.inEdge.fromVertex
		}
	}
	return scoreBlack, scoreWhite
}
