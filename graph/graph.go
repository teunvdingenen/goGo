package graph

import (
	"log"
	//	"math/rand"
	"crypto/rand"
	"time"
)

var topVertex *Vertex
var currentVertex *Vertex

var komi float32

var logger *log.Logger

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

type Edge struct {
	playX uint8
	playY uint8

	fromVertex *Vertex
	toVertex   *Vertex
}

func Initiate(boardSize uint8, log *log.Logger) {
	logger = log
	//	rand.Seed(time.Now().Unix())
	topVertex = new(Vertex)
	topVertex.boardState = new(Board)
	topVertex.boardState.Create(uint16(boardSize))
	topVertex.turn = 1
	topVertex.plyDepth = 0
	currentVertex = topVertex
}

func SetKomi(k float32) {
	komi = k
}

func Reset() {
	topVertex.boardState.Create(uint16(topVertex.boardState.size))
	currentVertex = topVertex
}

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

	board := currentVertex.boardState
	currentVertex.boardState = nil
	board.Play(c, bestMove.playX, bestMove.playY)

	currentVertex = bestMove.toVertex
	currentVertex.boardState = board
    logger.Printf("My Board is now:\n")
    logger.Printf(board.tostr())
	return bestMove.playX, bestMove.playY
}

func UpdateCurrentVertex(c, x, y uint8) {
	if c != currentVertex.turn {
		panic("Updating graph to wrong color's move")
	}
	var newCurrent *Vertex = nil
	for _, v := range currentVertex.outEdges {
		if v.playX == x && v.playY == y {
			currentVertex.outEdges = []*Edge{v}
			board := currentVertex.boardState
			currentVertex.boardState = nil
			board.Play(c, x, y)
			newCurrent = v.toVertex
			newCurrent.boardState = board
			break
		}
	}
	if newCurrent == nil {
		newCurrent = new(Vertex)
		board := currentVertex.boardState
		currentVertex.boardState = nil
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

func createGraph() {
	//	maxMoves := uint(currentVertex.boardState.size * currentVertex.boardState.size)
	doUntil := time.Now().Add(25 * time.Second)
	for time.Now().Before(doUntil) {
		toDepth := 80 + currentVertex.plyDepth + 1
		_ = doRoutine(currentVertex, toDepth)
	}
}

func doRoutine(fromVertex *Vertex, toDepth uint) bool {
	if toDepth == fromVertex.plyDepth {
		fromVertex.boardState = nil
		return true
	}

	var x, y uint8
	var score uint8
	var board *Board
	err := "error"

	if fromVertex == currentVertex {
		board = new(Board)
		board.Create(uint16(currentVertex.boardState.size))
		copy(board.s, fromVertex.boardState.s)
	} else {
		board = fromVertex.boardState
		fromVertex.boardState = nil
	}

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
			toCompare := fromVertex.inEdge.fromVertex.inEdge
			if x == toCompare.playX && y == toCompare.playY {
				err = "KO"
			}
		}
		if err != "" {
			board.Remove(x, y)
			i += 1
		}
		if possibilities == i {
			if fromVertex != currentVertex {
				fromVertex.boardState = nil
			}
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

func getRandomMove(b *Board) (x, y uint8) {
	empty := b.GetEmpty()
	//	i := rand.Intn(len(empty))
	bs := make([]byte, 8)
	_, _ = rand.Read(bs)
	var i int = 0
	for _, v := range bs {
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

func calcXY(index int, size uint8) (x, y uint8) {
	i := uint8(index)
	y = i / size
	x = i % size
	return x, y
}

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
	return scoreBlack, scoreWhite
}
