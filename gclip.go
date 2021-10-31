package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const red    = "\033[31m"
const reset  = "\033[0m"


// go install gclip.go
func main() {
	grepPtr := flag.String("grep", "", "select line that contains")
	flag.Parse()
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode() & os.ModeNamedPipe == 0 {
		return
	}

	// read from terminal
	endOfLine := make(chan bool)
	var b bytes.Buffer
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if grepPtr != nil && !strings.Contains(text, *grepPtr){
				continue
			}
			b.WriteString(text)
			b.WriteString("\n")
			fmt.Print(text + "\n")
		}
		endOfLine <- true
	}()


	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <- c:
		break
	case <- endOfLine:
		break
	}

	allLines := b.String()
	if err := clipboard.WriteAll(allLines); err != nil {
		panic(err)
	}
}