package debug

import (
	"fmt"
	"os"
)

func Log(format string, a ...any) {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, format+"\n", a...); err != nil {
		panic(err)
	}
}
