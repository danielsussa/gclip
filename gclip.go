package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
)

// go install gclip.go
func main() {
	grepPtr := flag.String("grep", "", "select line that contains")
	flag.Parse()

	grepTxt := *grepPtr

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}

	// read from terminal
	endOfLine := make(chan bool)
	var b bytes.Buffer
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text() + "\n"

			if grepTxt == "" {
				fmt.Print(text)
				b.WriteString(text)
				continue
			}

			ok, err := regexp.MatchString(grepTxt, text)
			if err != nil {
				panic(err)
			}

			if ok {
				color.Greenp(text)
				b.WriteString(text)
			} else {
				fmt.Print(text)
			}
		}
		endOfLine <- true
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		break
	case <-endOfLine:
		break
	}

	allLines := b.String()
	if err := clipboard.WriteAll(allLines); err != nil {
		panic(err)
	}
}

func hasFlag(f *string) bool{
	return f != nil && *f != ""
}