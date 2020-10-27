package monigo

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func getStrGID() string {
	return fmt.Sprintf("*%d*", getGID())
}

func int2str(i int) string {
	return strconv.Itoa(i)
}

func str2int(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func str2int64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func round(x float64, precision int) float64 {
	s := fmt.Sprintf("%."+strconv.Itoa(precision)+"f", x)
	x, _ = strconv.ParseFloat(s, 64)
	return x
}
