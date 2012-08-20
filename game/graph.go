package main

import (
	"math/rand"
)

var currentVertex *Vertex

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

func CreateGraph() {

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
		copy(newBoard.s, fromVertex.boardState.s)
		score, err = newBoard.Play(fromVertex.turn, x, y)

		toCompare := fromVertex.inEdge.fromVertex
		if newBoard.IsEqual(toCompare.boardState) {
			err = "KO"
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

	if newVertex.plyDepth > 200 {
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
		for vertex != nil {
			vertex.blackWins += 1
			vertex = vertex.inEdge.fromVertex
		}
	} else if scoreBlack < scoreWhite {
		vertex := v
		for vertex != nil {
			vertex.blackWins += 1
			vertex = vertex.inEdge.fromVertex
		}
	}
	return scoreBlack, scoreWhite
}
