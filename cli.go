package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Command struct is a representation of an isolated command executed by a user
type Command struct {
	Command  string `json:"Command"`
	Username string `json:"Username"`
	Amount   string `json:"Amount"`
	Stock    string `json:"Stock"`
	Filename string `json:"Filename"`
}

// FromStringToCommandStruct takes a line from the user command file as an input and returns a defined golang structure
func FromStringToCommandStruct(line string) (*Command, error) {
	/*
		The line variable value should have a similar format to this: '[1] ADD,oY01WVirLr,63511.53'
	*/
	line = strings.Split(line, " ")[1]
	commandVars := strings.Split(line, ",")
	cmd := commandVars[0]

	if cmd == "ADD" {
		return &Command{Command: cmd, Username: commandVars[1], Amount: commandVars[2]}, nil
	}

	if cmd == "COMMIT_BUY" || cmd == "CANCEL_BUY" || cmd == "COMMIT_SELL" || cmd == "CANCEL_SELL" || cmd == "DISPLAY_SUMMARY" {
		return &Command{Command: cmd, Username: commandVars[1]}, nil
	}

	if cmd == "BUY" || cmd == "SELL" || cmd == "SET_BUY_AMOUNT" || cmd == "SET_BUY_TRIGGER" || cmd == "SET_SELL_AMOUNT" || cmd == "SET_SELL_TRIGGER" {
		return &Command{Command: cmd, Username: commandVars[1], Stock: commandVars[2], Amount: commandVars[3]}, nil
	}

	if cmd == "QUOTE" || cmd == "CANCEL_SET_BUY" || cmd == "CANCEL_SET_SELL" {
		return &Command{Command: cmd, Username: commandVars[1], Stock: commandVars[2]}, nil
	}

	if cmd == "DUMPLOG" {
		if len(commandVars) == 3 {
			// case: DUMPLOG,userid,filename
			return &Command{Command: cmd, Username: commandVars[1], Filename: commandVars[2]}, nil
		} else {
			// case: DUMPLOG,filename
			return &Command{Command: cmd, Filename: commandVars[1]}, nil
		}
	}

	if cmd == "DISPLAY_SUMMARY" {
		return &Command{Command: cmd, Username: commandVars[1]}, nil
	}

	return nil, fmt.Errorf("unable to conver given line: %s into golang struct", line)
}

func checkError(e error, additionalMessage string) {
	if e != nil {
		fmt.Println(additionalMessage)
		panic(e)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please follow the following format: go run cmd.go <path_to_workload_file.txt>")
		panic("Unexpected number of arguments")
	}

	commandsFilePath := os.Args[1]
	data, err := os.ReadFile(filepath.Clean(commandsFilePath))
	checkError(err, "Error while reading file")

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		requestData, err := FromStringToCommandStruct(line)
		checkError(err, "")

		fmt.Printf("iteration: %d requestData: %#v\n", i+1, requestData)

		parsedJson, err := json.Marshal(requestData)
		checkError(err, "Couldn't parse golang struct to JSON")

		res, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(parsedJson))
		checkError(err, "Got error while doing a post request")

		fmt.Printf("Got response code: %v\n", res.StatusCode)
	}

}
