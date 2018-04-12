package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestExecuteCmds(t *testing.T) {
	out, err := Execute("basecoind init")
	//require.NoError(t, err)
	var initRes map[string]string
	out = "{" + strings.SplitN(out, "{", 2)[1]
	_ = json.Unmarshal([]byte(out), &initRes)
	moneyKey := initRes["secret"]
	fmt.Println(moneyKey)
	//time.Sleep(time.Second * 2)

	wc, _, _, err := GoExecute("basecli keys add hoot")
	fmt.Println(err)
	//require.NoError(t, err)

	time.Sleep(time.Second)
	_, err = wc.Write([]byte("1234567890\n"))
	fmt.Println(err)
	//require.NoError(t, err)

	fmt.Println(err)
	fmt.Println("hoot")

	time.Sleep(time.Second * 3)
	wc1, _, _, err := GoExecute("basecli keys add pig --recover")
	time.Sleep(time.Second)
	fmt.Println(err)

	//require.NoError(t, err)
	_, err = wc1.Write([]byte("1234567890\n"))
	fmt.Println(err)
	//require.NoError(t, err)
	time.Sleep(time.Second)
	_, err = wc1.Write([]byte(moneyKey + "\n")) // XXX why dis enter no work
	fmt.Println(err)
	//require.NoError(t, err)
	//out = <-outChan
	//fmt.Println(out)
	//err = <-errChan
	//fmt.Println(err)

	time.Sleep(time.Second)
	out, err = Execute("basecli keys show pig")
	//require.NoError(t, err)
	pigAddr := strings.TrimLeft(out, "pig\t")
	fmt.Println(pigAddr)

	//out, err := Execute(" init")

	wc2, _, _, _ := GoExecute("basecoind start")
	defer wc2.Close()
	time.Sleep(time.Second)
	fmt.Println(err)
	//require.NoError(t, err)
	out, err = Execute(fmt.Sprintf("basecli account %v", pigAddr))
	fmt.Println(out)
	//require.NoError(t, err)
	fmt.Println(err)
}
