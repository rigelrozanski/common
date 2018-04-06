package common

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Execute the command, return standard output and error
func ExecuteT(t *testing.T, command string) (out string) {

	//split command into command and args
	var outByte []byte
	split := strings.Split(command, " ")
	require.NotZero(t, len(split), "no command provided")
	cmd := exec.Command(split[:]...)
	bz, err := cmd.CombinedOutput()
	require.NoError(t, err)
	out = strings.Trim(string(bz), "\n") //trim any new lines
	return out
}

// Asynchronously execute the command, return standard output and error
func GoExecuteT(t *testing.T, command string) (process *exec.Process, out chan string) {
	//split command into command and args
	var outByte []byte
	split := strings.Split(command, " ")
	require.NotZero(t, len(split), "no command provided")
	cmd := exec.Command(split[:]...)
	go func() {
		bz, err = cmd.CombinedOutput()
		require.NoError(t, err)
	}()
	out = strings.Trim(string(bz), "\n") //trim any new lines
	return cmd.Process, stdOut
}

//___________________________________________________________________________________________-

// Execute the command, return standard output and error
func Execute(command string) (out string, err error) {
	//split command into command and args
	var outByte []byte
	split := strings.Split(command, " ")
	if len(split) == 0 {
		return nil, "", errors.New("no command provided")
	}
	cmd := exec.Command(split[:]...)
	bz, err := cmd.CombinedOutput()
	out = strings.Trim(string(bz), "\n") //trim any new lines
	return out, err
}

// Asynchronously execute the command, return standard output and error
func GoExecute(command string) (process *exec.Process, out chan string, err chan error) {
	//split command into command and args
	var outByte []byte
	var cmd *exec.Cmd
	split := strings.Split(command, " ")
	if len(split) == 0 {
		return nil, "", errors.New("no command provided")
	}
	cmd = exec.Command(split[:]...)
	go func() {
		bz, err = cmd.CombinedOutput()
	}()
	out = strings.Trim(string(bz), "\n") //trim any new lines
	return cmd.Process, stdOut, err
}
