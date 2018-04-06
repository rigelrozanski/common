package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteCmds(t *testing.T) {
	proc, outChan, errChan := GoExecute("democoind start")
	defer proc.Kill()
	out, err := Execute("sleep 5")
	require.NoError(t, err)
	out, err := Execute("democli account C2F2E199A0CE9C7809DDD0EFF774C7DCB4529C26")
	require.NoError(t, err)
}
