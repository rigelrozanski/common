package common

import (
	"fmt"
	"testing"
	//"github.com/stretchr/testify/assert"
)

func TestExecuteCmds(t *testing.T) {
	go Execute("democoind start")

	cmds2 := []string{
		"sleep 5",
		"democli account C2F2E199A0CE9C7809DDD0EFF774C7DCB4529C26",
	}
	_, err := ExecuteCmds(cmds2, true, false)
	//assert.NoError(t, err)
	if err != nil {
		fmt.Printf("faifaifalflaf")
	}
}
