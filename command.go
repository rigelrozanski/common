package common

import "os/exec"

// Execute - Execute the command, return standard output and error
func Execute(command string) (stdOut string, err error) {
	var outByte []byte
	outByte, err = exec.Command(command).Output()
	stdOut = string(outByte)
	return
}
