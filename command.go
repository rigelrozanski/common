package common

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Execute the command, return standard output and error
func Execute(command string) (stdOut string, err error) {
	//split command into command and args
	var outByte []byte
	split := strings.Split(command, " ")
	switch len(split) {
	case 0:
		return "", errors.New("no command provided")
	case 1:
		outByte, err = exec.Command(split[0]).CombinedOutput()
	default:
		outByte, err = exec.Command(split[0], split[1:]...).CombinedOutput()
	}
	stdOut = string(outByte)
	stdOut = strings.Trim(stdOut, "\n") //trim any new lines
	return
}

// Execute a bunch o' commands, return standard outputs and error
func ExecuteCmds(commands []string, print, haltOnError bool) (stdOuts []string, errs []error) {

	stdOuts = make([]string, len(commands))
	if haltOnError {
		errs = make([]error, 1)
	} else {
		errs = make([]error, len(commands))
	}

	errI := 0
	for i, command := range commands {

		stdOut, err := Execute(command)
		if err != nil {

			err = fmt.Errorf("Error in command \"%v\":\n%s", command, err)
			errs[errI] = err
			errI++

			if print {
				fmt.Println(err)
			}
			if haltOnError {
				return
			}
		}
		if print {
			fmt.Println(stdOut)
		}
		stdOuts[i] = stdOut
	}
	return
}
