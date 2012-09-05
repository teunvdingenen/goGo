/*Commands format: [id] command_name arguments\n
    id = optional int

   Required commands:                       TODO
    protocol_version                        return 2
    name                                    return "goGo"
    version                                 return version nr (1.0)
    known_commands [command_name, string]   return bool (known or not)
    list_commands                           return known commands, one per row
    quit                                    respond and close connection
    boardsize [size, int]                   clear board (reset) en set to size)
    clear_board                             clear board reset score etc.
    komi [komi, float]                      set komi 
    play [move]                             set the move 
    genmove [color]                         gen move for color, return move 
   Would like to do:
    fixed_handicap [nr_of_stones]           set handicap stones, return pos of st.
    place_free_handicap [nr_of_stones]      set to handi. to preffered pos, return
    set_free_handicap[vertices]             set handi. to given positions
    undo                                    revert to one state before, or error:
                                            "cannot undo"
    time_settings [time]                    set time settings

    vertex = one letter, one number ( eg B13, j11 ) or "pass" or "resign"
    color = white, w, black or b
    move = color vertex (seperated by space) eg "white h10", "B F5", "w pass"

    time = main_time byo_yomi_time byo_yomi_stones

    Responses:
    succes:     =[id] response \n\n
                =[id]\n\n
    failure:
                ?[id] error_message
        id may be omitted if and only if it was omited in command
*/
package gtp

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var commands = []string{"protocol_version", "name", "version", "known_commands", "list_commands", "quit", "boardsize", "clear_board", "komi", "play", "genmove"}
var received chan []uint8

//Listen is run in a seperate goroutine and continuesly listens to std.in for new commands
func Listen() {
	listener := bufio.NewReader(os.Stdin)
	for {
		input, _ := listener.ReadBytes('\n')
		add(string(input))
		//        fmt.Printf("Read: %s\n", input )
	}
}

//Translates a coordinate in string form (ie D3) to x y positions (3,2) for usage in a matrix
func ToXY(s string) (uint8, uint8) {
	s = strings.ToLower(s)
	if s == "pass" {
		return 254, 254
	} else if s == "resign" {
		return 255, 255
	}
	x := s[0]
	x = x - 'a'
	if x > 8 {
		x -= 1
	}
	yt, _ := strconv.ParseUint(s[1:], 10, 8)
	y := uint8(yt) - 1
	return x, y
}

//Translates x,y positions to a coordinate string like D3
func FromXY(x, y uint8) string {
	var xstr string
	if x < 8 {
		xstr = string(int(x + 'a'))
	} else {
		x += 1
		xstr = string(int(x + 'a'))
	}
	return xstr + strconv.Itoa(int(y)+1)
}

//Translates a color string (black, white) to 1, 2
func FromColorStr(s string) uint8 {
	var retValue uint8
	if s[0] == 'w' || s[0] == 'W' {
		retValue = 2
	} else if s[0] == 'b' || s[0] == 'B' {
		retValue = 1
	} else {
		panic("What color are you trying to play!?!")
	}
	return retValue
}

//If c is 1 returns b, if c is 2 returns w
func ToColorStr(c uint8) string {
	var retValue string
	if c == 2 {
		retValue = "w"
	} else if c == 1 {
		retValue = "b"
	}
	return retValue
}

//Adds a command received in std.in to the command channel
func add(s string) {
	arg := make([]uint8, 4)
	s = strings.Trim(s, "\n")
	sSlice := strings.Split(s, " ")

	switch sSlice[0] {
	case "protocol_version":
		arg[0] = 6
	case "name":
		arg[0] = 7
	case "version":
		arg[0] = 8
	case "list_commands":
		arg[0] = 9
	case "quit":
		arg[0] = 5
	case "clear_board":
		arg[0] = 4
	case "known_commands":
		arg[0] = 10
		known := 0
		for _, v := range commands {
			if strings.EqualFold(sSlice[1], v) {
				known = 1
				break
			}
		}
		arg[1] = uint8(known)
	case "genmove":
		arg[0] = 1
		arg[1] = FromColorStr(sSlice[1])
	case "boardsize":
		size, _ := strconv.ParseUint(sSlice[1], 10, 8)
		arg[0] = 3
		arg[1] = uint8(size)
	case "komi":
		komi, _ := strconv.ParseFloat(sSlice[1], 32)
		komi = komi / 0.5 // convert to number of half points
		komiInt := uint8(komi)
		arg[0] = 2
		arg[1] = komiInt
	case "play":
		arg[0] = 0
		arg[1] = FromColorStr(sSlice[1])
		arg[2], arg[3] = ToXY(sSlice[2])
	default:
		arg[0] = 255
	}
	received <- arg
}

//Starts the gtp engine
func Start(todo chan []uint8) int {
	received = todo
	go Listen()
	return 0
}

//Responds to the server according to the gtp protocol
func Respond(s string, succes bool) {
	var b byte
	if succes {
		b = '='
	} else {
		b = '?'
	}
	fmt.Printf("%c %s\n\n", b, s)
}

//Returns a string of all known commands
func ListCommands() string {
	return strings.Join(commands, "\n")
}
