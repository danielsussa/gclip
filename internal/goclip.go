package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"golang.design/x/clipboard"
)

// go install goclip.go
func main() {
	grepPtr := flag.String("grep", "", "select line that contains")
	flag.Parse()

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}
	goClip(&consoleController{}, *grepPtr)
}

func goClip(controller Controller, grepTxt string) {
	// read from terminal
	endOfLine := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(controller.Reader())
		for scanner.Scan() {
			text := scanner.Text() + "\n"

			if grepTxt == "" {
				controller.Print(text)
				controller.Store(text)
				continue
			}

			ok, err := regexp.MatchString(grepTxt, text)
			if err != nil {
				panic(err)
			}

			if ok {
				controller.PrintGreen(text)
				controller.Store(text)
			} else {
				controller.Print(text)
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

	allLines := controller.GetStored()

	clipboard.Write(clipboard.FmtText, []byte(allLines))
}