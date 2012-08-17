
package main

import (
    "goGo/gtp"
    "fmt"
)

type Board struct {
    s []uint8
    size uint8
}

var komi float32
var board Board

var todo chan []uint8 
var running bool = true

func main( ) {
    todo = make( chan []uint8 )
    gtp.Start( todo )
//    go gtp.Listen( )
    watchTodo( )
}

func watchTodo( ) {
    for running {
        c := <-todo
        switch c[0] {
            case 0: //play
                board.Play( c[1], c[2], c[3] )
                gtp.Respond( "", true )
            case 1: //genmove
                fmt.Printf("Genmove for color: %d\n", c[1])
                gtp.Respond( gtp.FromXY(0,0), true )
            case 2: //komi
                komi = float32(c[1]) * 0.5
                gtp.Respond( "", true )
            case 3: //boardsize
                initBoard( uint16(c[1]) )
                gtp.Respond( "", true )
                board.prnt( )
            case 4: //clearboard
                //TODO clear board
                gtp.Respond( "", true )
            case 5: //quit
                //make quit function
                running = false
                gtp.Respond( "", true )
            default:
                fmt.Printf("Something went wrong in watchToDo()")

        }
    }
}

func initBoard( size uint16 ) {
    board = Board{ make([]uint8, size*size), uint8(size)}
}

func ( b *Board ) Play( color, x, y uint8 ) bool {
    isLegal := true
    //TODO check legality (Tesuki, Ko etc.)
    board.put( color, x, y )
    //TODO remove dead stones
    board.prnt( )
    return isLegal
}

func ( b *Board ) put( c, x, y uint8 ) {
    if y >= b.size || x >= b.size || c > 1 {
        panic( "Invalid board.put operation")
    }
    b.s[x+y*b.size] = c
}

func (b *Board) prnt( ) {
    var i uint8 = 1
    for _,c := range b.s {
        fmt.Printf("%d ", c )
        if i == b.size {
            fmt.Printf("\n")
            i = 1
        } else { i+=1 }
    }
}
