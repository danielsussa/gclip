package main

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"io"
	"os"
)

type Controller interface{
	Reader() io.Reader
	Print(txt string)
	PrintGreen(txt string)
	Store(txt string)
	GetStored()string
}

// production controller

type consoleController struct {
	b bytes.Buffer
}

func (cc *consoleController) Reader() io.Reader {
	return os.Stdin
}

func (cc *consoleController) Print(txt string) {
	fmt.Print(txt)
}

func (cc *consoleController) PrintGreen(txt string) {
	color.Greenp(txt)
}

func (cc *consoleController) Store(txt string) {
	cc.b.WriteString(txt)
}

func (cc *consoleController) GetStored() string {
	return cc.b.String()
}

// fake controller

type customReader struct {
	pReader io.PipeReader
	in []string
}

func (c *customReader) Read(p []byte) (n int, err error) {
	if len(c.in) == 0 {
		return 0, io.EOF
	}
	lenR := len(c.in[0]) + 1
	txt := c.in[0] + "\n"
	copy(p, txt)

	c.in = c.in[1:]
	return lenR, nil
}

type fakeController struct {
	counter int
	b bytes.Buffer
	r *customReader
	p func(string, int)
	pg func(string, int)
	str func(string)
}

type ConfigFake struct {
	p func(string, int)
	pg func(string, int)
	str func(string)
}

func NewFakeController(in []string, conf ConfigFake) *fakeController{
	return &fakeController{
		pg: conf.pg,
		p: conf.p,
		str: conf.str,
		r: &customReader{in: in},
	}
}

func (cc *fakeController) Reader() io.Reader {
	return cc.r
}

func (cc *fakeController) Print(txt string) {
	cc.p(txt, cc.counter)
	cc.counter++
}

func (cc *fakeController) PrintGreen(txt string) {
	cc.pg(txt, cc.counter)
	cc.counter++
}

func (cc *fakeController) Store(txt string) {
	cc.b.WriteString(txt)
}

func (cc *fakeController) GetStored() string {
	cc.str(cc.b.String())
	return cc.b.String()
}