package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecuteCmds(t *testing.T) {
	out, err := Execute("basecoind init")
	require.NoError(t, err)
	var initRes map[string]string
	out = "{" + strings.SplitN(out, "{", 2)[1]
	_ = json.Unmarshal([]byte(out), &initRes)
	moneyKey := initRes["secret"]
	fmt.Println(moneyKey)

	wc, outChan, errChan, err := GoExecute("basecli keys add rigey --recover")
	require.NoError(t, err)
	time.Sleep(time.Second * 1)
	_, err = wc.Write([]byte("1234567890\n"))
	require.NoError(t, err)
	time.Sleep(time.Second * 1)
	_, err = wc.Write([]byte("1234567890\n"))
	require.NoError(t, err)
	require.NoError(t, err)
	time.Sleep(time.Second * 1)
	_, err = wc.Write([]byte(fmt.Sprintf(moneyKey + "\n"))) // XXX why dis enter no work
	require.NoError(t, err)
	out = <-outChan
	fmt.Println(out)
	err = <-errChan
	fmt.Println(err)
	wc.Close()

	//out, err := Execute(" init")

	//proc, outChan, errChan := GoExecute("basecoind start")
	//defer proc.Kill()
	//out, err := Execute("sleep 5")
	//require.NoError(t, err)
	//out, err := Execute("democli account C2F2E199A0CE9C7809DDD0EFF774C7DCB4529C26")
	//require.NoError(t, err)
}
