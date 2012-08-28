package main

import (
	"goGo/graph"
	"goGo/gtp"
)

var komi float32

var todo chan []uint8
var running bool = true

func main() {
	todo = make(chan []uint8, 10)
	gtp.Start(todo)
	watchTodo()
}

func watchTodo() {
	for running {
		c := <-todo
		switch c[0] {
		case 0: //play
			graph.UpdateCurrentVertex(c[1], c[2], c[3])
			gtp.Respond("", true)
		case 1: //genmove
			x, y := graph.GetMove(c[1])
			gtp.Respond(gtp.FromXY(x, y), true)
		case 2: //komi
			komi = float32(c[1]) * 0.5
			gtp.Respond("", true)
		case 3: //boardsize
			graph.Initiate(c[1])
			gtp.Respond("", true)
		case 4: //clearboard
			graph.Reset()
			gtp.Respond("", true)
		case 5: //quit
			//TODO make quit function
			running = false
			gtp.Respond("", true)
		case 6: //protocol
			gtp.Respond("2.0", true)
		case 7: //name
			gtp.Respond("goGo", true)
		case 8: //version
			gtp.Respond("0.1", true)
		case 9: //list_commands
			gtp.Respond(gtp.ListCommands(), true)
		case 10: //known_commands
			if c[1] == 0 {
				gtp.Respond("false", true)
			} else {
				gtp.Respond("true", true)
			}
		default:
			gtp.Respond("unknown command", false)
		}
	}
}
