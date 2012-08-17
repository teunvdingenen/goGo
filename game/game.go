package main

import (
	"fmt"
	"goGo/gtp"
)

type Board struct {
	s    []uint8
	size uint8
}

var komi float32
var board Board

var todo chan []uint8
var running bool = true

func main() {
	todo = make(chan []uint8)
	gtp.Start(todo)
	watchTodo()
}

func watchTodo() {
	for running {
		c := <-todo
		switch c[0] {
		case 0: //play
			board.Play(c[1], c[2], c[3])
			gtp.Respond("", true)
		case 1: //genmove
			fmt.Printf("Genmove for color: %d\n", c[1])
			gtp.Respond(gtp.FromXY(0, 0), true)
		case 2: //komi
			komi = float32(c[1]) * 0.5
			gtp.Respond("", true)
		case 3: //boardsize
			initBoard(uint16(c[1]))
			gtp.Respond("", true)
			board.prnt()
		case 4: //clearboard
			//TODO clear board
			gtp.Respond("", true)
		case 5: //quit
			//make quit function
			running = false
			gtp.Respond("", true)
		default:
			fmt.Printf("Something went wrong in watchToDo()")

		}
	}
}

func initBoard(size uint16) {
	board = Board{make([]uint8, size*size), uint8(size)}
}

func (b *Board) Play(color, x, y uint8) bool {
	isLegal := true
	//TODO check legality (Tesuki, Ko etc.)
	board.place(color, x, y)
	//TODO remove dead stones
	board.prnt()
	return isLegal
}

func (b *Board) place(c, x, y uint8) {
	if y >= b.size || x >= b.size || c > 2 {
		panic("Invalid board.place operation") //TODO don't panic
	}
	b.s[x+y*b.size] = c
}

func (b *Board) isEqual( c *Board ) bool {
    for i,v := range b.s {
        if v != c.s[i] {
            return false
        }
    }
    return true
}

func getAdjecent(x, y uint8) (xa, xb, ya, yb uint8) {
    xa = x - 1
	xb = x + 1
	ya = y - 1
	yb = y + 1
    return xa, xb, ya, yb
}

func (b *Board) hasFreedom(i uint8, xs, ys []uint8 ) bool {
    xa, xb, ya, yb := getAdjecent( xs[i], ys[i] )

    if xa < b.size && b.getColor(xa, ys[i]) == 0 {
        return true
    }
    if xb < b.size && b.getColor(xb, ys[i]) == 0 {
        return true
    }
    if ya < b.size && b.getColor(xs[i], ya) == 0 {
        return true
    }
    if yb < b.size && b.getColor(xs[i], yb) == 0 {
        return true
    }
    i += 1
    if len(xs) == int(i) {
        return false
    }
    return b.hasFreedom( i, xs, ys )
}

func (b *Board) getGroup(c, i uint8, xs, ys []uint8) bool {
    xa, xb, ya, yb := getAdjecent( xs[i], ys[i] )

	if xa < b.size && b.getColor(xa, ys[i]) == c && !isPresent(xa, ys[i], xs, ys) {
		xs = append(xs, xa)
		ys = append(ys, ys[i])
	}
	if xb < b.size && b.getColor(xb, ys[i]) == c && !isPresent(xb, ys[i], xs, ys) {
		xs = append(xs, xb)
		ys = append(ys, ys[i])
	}
	if ya < b.size && b.getColor(xs[i], ya) == c && !isPresent(xs[i], ya, xs, ys) {
		xs = append(xs, xs[i])
		ys = append(ys, ya)
	}
	if yb < b.size && b.getColor(xs[i], yb) == c && !isPresent(xs[i], yb, xs, ys) {
		xs = append(xs, xs[i])
		ys = append(ys, yb)
	}
	i += 1
	if len(xs) == int(i) {
		return true
    }
    return b.getGroup(c, i, xs, ys)
}

func isPresent(x, y uint8, xs, ys []uint8) bool {
	for i := 0; i < len(xs); i++ {
		if xs[i] == x && ys[i] == y {
			return true
		}
	}
	return false
}

func (b *Board) getColor(x, y uint8) uint8 {
	if y >= b.size || x >= b.size {
		panic("Invalid board.getColor operation") //TODO don't panic
	}
	return b.s[x+y*b.size]
}

func (b *Board) prnt() {
	var i uint8 = 1
	for _, c := range b.s {
		fmt.Printf("%d ", c)
		if i == b.size {
			fmt.Printf("\n")
			i = 1
		} else {
			i += 1
		}
	}
}
