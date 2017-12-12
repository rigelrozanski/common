package common

import "os/exec"

// Execute - Execute the command, return standard output and error
func Execute(command string, args ...string) (stdOut string, err error) {
	var outByte []byte
	outByte, err = exec.Command(command, args...).Output()
	stdOut = string(outByte)
	return
}
