package common

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

//// Execute the command, return standard output and error
//func ExecuteT(t *testing.T, command string) (out string) {

////split command into command and args
//var outByte []byte
//split := strings.Split(command, " ")
//require.NotZero(t, len(split), "no command provided")
//cmd := exec.Command(split[:]...)
//bz, err := cmd.CombinedOutput()
//require.NoError(t, err)
//out = strings.Trim(string(bz), "\n") //trim any new lines
//return out
//}

//// Asynchronously execute the command, return standard output and error
//func GoExecuteT(t *testing.T, command string) (process *exec.Process, out chan string) {
////split command into command and args
//var outByte []byte
//split := strings.Split(command, " ")
//require.NotZero(t, len(split), "no command provided")
//cmd := exec.Command(split[:]...)
//go func() {
//bz, err = cmd.CombinedOutput()
//require.NoError(t, err)
//}()
//out = strings.Trim(string(bz), "\n") //trim any new lines
//return cmd.Process, stdOut
//}

//___________________________________________________________________________________________-

// Execute the command, return standard output and error
func Execute(command string) (out string, err error) {
	//split command into command and args
	split := strings.Split(command, " ")

	var cmd *exec.Cmd
	switch len(split) {
	case 0:
		return "", errors.New("no command provided")
	case 1:
		cmd = exec.Command(split[0], split[:]...)
	default:
		cmd = exec.Command(split[0], split[1:]...)
	}
	bz, err := cmd.CombinedOutput()
	out = strings.Trim(string(bz), "\n") //trim any new lines
	return out, err
}

// Asynchronously execute the command, return standard output and error
func GoExecute(command string) (pipe io.WriteCloser, outChan chan string, errChan chan error, err error) {
	//split command into command and args
	split := strings.Split(command, " ")

	var cmd *exec.Cmd
	switch len(split) {
	case 0:
		err = errors.New("no command provided")
		return nil, outChan, errChan, err
	case 1:
		cmd = exec.Command(split[0], split[:]...)
	default:
		cmd = exec.Command(split[0], split[1:]...)
	}
	pipe, err = cmd.StdinPipe()
	if err != nil {
		err = errors.New("no command provided")
		return nil, outChan, errChan, err
	}
	go func() {
		bz, err := cmd.CombinedOutput()
		errChan <- err
		outChan <- strings.Trim(string(bz), "\n") //trim any new lines
	}()
	return pipe, outChan, errChan, nil
}
