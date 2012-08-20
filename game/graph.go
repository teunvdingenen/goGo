package main

import (
	"math/rand"
)

var vertices []*Vertex
var koCheck Board

type Vertex struct {
	whiteWins  int
	blackWins  int
	boardState Board
	turn       uint8

	plyDepth int
	outEdges []*Edge
	inEdges  []*Edge
}

type Edge struct {
	playX uint8
	playY uint8

	fromVertex *Vertex
	toVertex   *Vertex
}

func Start() {

}

func doRoutine(fromVertex *Vertex, toDepth int) bool {
	if toDepth == fromVertex.plyDepth {
		return true
	}
	x, y := getRandomMove(fromVertex.boardState)
	var newBoard Board
	copy(newBoard.s, fromVertex.boardState.s)
	//score, err := newBoard.Play(fromVertex.turn, x, y )
	_, err := newBoard.Play(fromVertex.turn, x, y)

	if err != "" { //try again
		doRoutine(fromVertex, toDepth)
	}

	//Ko check

	newVertex := new(Vertex)
	newEdge := new(Edge)
	newEdge.playX = x
	newEdge.playY = y
	newEdge.fromVertex = fromVertex
	newEdge.toVertex = newVertex

	fromVertex.outEdges = append(fromVertex.outEdges, newEdge)
	newVertex.inEdges = append(newVertex.inEdges, newEdge)

	newVertex.boardState = newBoard
	if fromVertex.turn == 1 {
		newVertex.turn = 0
	} else {
		newVertex.turn = 1
	}
	newVertex.plyDepth = fromVertex.plyDepth + 1

	//score
	//continue to new vertex
	return true
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
