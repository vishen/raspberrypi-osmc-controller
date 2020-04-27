package main

import (
	"io"
	"log"
	"os"
	"time"
)

func main() {

	filenames := []string{
		"/dev/hidraw0",
		"/dev/hidraw1",
	}

	resultsChan := make(chan *FunctionState, 1)
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		go func(r io.Reader) {
			var prevFunctionState *FunctionState
			for {
				functionState, err := parseFunctionState(r)
				if err != nil && err != io.EOF {
					log.Printf("error: unable to parse function from %q: %v", filename, err)
					continue
				}
				if functionState == nil {
					time.Sleep(2 * time.Second)
					continue
				}
				if prevFunctionState != nil && functionState.Action == Action_UP {
					functionState.Function = prevFunctionState.Function
				}
				resultsChan <- functionState
				prevFunctionState = functionState
			}
		}(f)
	}

	for fs := range resultsChan {
		log.Printf("state=%#v\n", fs)
	}
}

type Action string

const (
	Action_DOWN Action = "DOWN"
	Action_UP   Action = "UP"
)

type Function string

const (
	Function_HOME        Function = "HOME"
	Function_ARROW_UP    Function = "ARROW_UP"
	Function_ARROW_LEFT  Function = "ARROW_LEFT"
	Function_ARROW_RIGHT Function = "ARROW_RIGHT"
	Function_ARROW_DOWN  Function = "ARROW_DOWN"
	Function_OK          Function = "OK"
	Function_INFO        Function = "INFO"
	Function_RETURN      Function = "RETURN"
	Function_PLAY_PAUSE  Function = "PLAY_PAUSE"
	Function_STOP        Function = "STOP"
	Function_PREV        Function = "PREV"
	Function_NEXT        Function = "NEXT"
	Function_LIST        Function = "LIST"
)

type FunctionState struct {
	Action    Action
	Function  Function
	Timestamp time.Time
}

func parseFunctionState(r io.Reader) (*FunctionState, error) {
	fs := &FunctionState{
		Action:    Action_DOWN,
		Timestamp: time.Now(),
	}
	b1 := make([]byte, 4)
	if _, err := r.Read(b1); err != nil {
		return fs, err
	}
	switch b1[0] {
	default:
		return nil, nil
	case 0x00:
		switch b1[2] {
		case 'J':
			fs.Function = Function_HOME
		case 'R':
			fs.Function = Function_ARROW_UP
		case 'P':
			fs.Function = Function_ARROW_LEFT
		case 'O':
			fs.Function = Function_ARROW_RIGHT
		case 'Q':
			fs.Function = Function_ARROW_DOWN
		case '(':
			fs.Function = Function_OK
		default:
			fs.Action = Action_UP
		}
	case 0x04:
		switch b1[1] {
		case '`':
			fs.Function = Function_INFO
		case '$':
			fs.Function = Function_RETURN
		case 0xcd:
			fs.Function = Function_PLAY_PAUSE
		case '&':
			fs.Function = Function_STOP
		case 0xb4:
			fs.Function = Function_PREV
		case 0xb3:
			fs.Function = Function_NEXT
		default:
			fs.Action = Action_UP
		}
	case 0x05:
		switch b1[1] {
		case 0x08:
			fs.Function = Function_LIST
		default:
			fs.Action = Action_UP
		}
	}
	return fs, nil
}
