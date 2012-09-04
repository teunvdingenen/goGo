package main

import (
	"goGo/graph"
	"goGo/gtp"
	"log"
	"os"
)

var todo chan []uint8

var fileLog *log.Logger

var running bool = true
var hasPassed bool
var hasResigned bool

func main() {
	file, _ := os.Create("logfile")
	fileLog = log.New(file, "", 0)
	todo = make(chan []uint8, 32)
	gtp.Start(todo)
	watchTodo()
}

func watchTodo() {
	for running {
		c := <-todo
		switch c[0] {
		case 0: //play
			if c[2] == 254 && c[3] == 254 {
				hasPassed = true
				gtp.Respond("", true)
			} else if c[2] == 255 && c[3] == 255 {
				hasResigned = true
				gtp.Respond("", true)
			} else {
				graph.UpdateCurrentVertex(c[1], c[2], c[3])
				gtp.Respond("", true)
			}
		case 1: //genmove
			if hasPassed {
				gtp.Respond("pass", true)
			} else {
			    x, y := graph.GetMove(c[1])
			    if x == 255 && y == 255 {
				    gtp.Respond("resign", true)
                } else {
                    gtp.Respond(gtp.FromXY(x, y), true)
                }
			}
		case 2: //komi
			graph.SetKomi(float32(c[1]) * 0.5)
			gtp.Respond("", true)
		case 3: //boardsize
            hasPassed = false
            hasResigned = false
			graph.Initiate(c[1], fileLog)
			gtp.Respond("", true)
		case 4: //clearboard
            hasPassed = false
            hasResigned = false
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
