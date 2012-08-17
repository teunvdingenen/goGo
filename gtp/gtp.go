/* 
  Commands format: [id] command_name arguments\n
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

func Listen() {
	listener := bufio.NewReader(os.Stdin)
	for {
		input, _ := listener.ReadBytes('\n')
		go add(string(input))
		//        fmt.Printf("Read: %s\n", input )
	}
}

func ToXY(s string) (uint8, uint8) {
	s = strings.ToLower(s)
	x := s[0]
	x = x - 'a'
	yt, _ := strconv.ParseUint(s[1:], 10, 8)
	y := uint8(yt) - 1
	return x, y
}

func FromXY(x, y uint8) string {
	xstr := string(int(x + 'a'))
	return xstr + strconv.Itoa(int(y)+1)
}

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

func ToColorStr(c uint8) string {
	var retValue string
	if c == 2 {
		retValue = "w"
	} else if c == 1 {
		retValue = "b"
	}
	return retValue
}

func add(s string) {
	shouldAdd := false
	arg := make([]uint8, 4)
	s = strings.Trim(s, "\n")
	sSlice := strings.Split(s, " ")

	switch sSlice[0] {
	case commands[0]: //protocol
		Respond("2.0", true)
	case commands[1]: //name
		Respond("goGo", true)
	case commands[2]: //version
		Respond("0.1", true)
	case commands[4]: //list_commands
		Respond(strings.Join(commands, "\n"), true)
	case commands[5]: // quit
		shouldAdd = true
		arg[0] = 5
	case commands[7]: // clear_board
		shouldAdd = true
		arg[0] = 4
	case commands[3]: //known_commands
		for _, v := range commands {
			if strings.EqualFold(sSlice[1], v) {
				Respond("true", true)
				return
			}
		}
		Respond("false", true)
	case commands[10]: //genmove
		shouldAdd = true
		arg[0] = 1
		arg[1] = FromColorStr(sSlice[1])
	case commands[6]: // boardsize
		shouldAdd = true
		size, _ := strconv.ParseUint(sSlice[1], 10, 8)
		arg[0] = 3
		arg[1] = uint8(size)
	case commands[8]: //komi
		shouldAdd = true
		komi, _ := strconv.ParseFloat(sSlice[1], 32)
		komi = komi / 0.5 // convert to number of half points
		komiInt := uint8(komi)
		arg[0] = 2
		arg[1] = komiInt
	case commands[9]: //play
		shouldAdd = true
		arg[0] = 0
		arg[1] = FromColorStr(sSlice[1])
		arg[2], arg[3] = ToXY(sSlice[2])
	default:
		Respond("unknown command", false)
	}
	// get command type
	// settle arguments
	if shouldAdd {
		received <- arg
	}
}

func Start(todo chan []uint8) int {
	received = todo
	go Listen()
	return 0
}

func Respond(s string, succes bool) {
	var b byte
	if succes {
		b = '='
	} else {
		b = '?'
	}
	fmt.Printf("%c %s\n\n", b, s)
}
