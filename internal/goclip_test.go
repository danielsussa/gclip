package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGClip(t *testing.T) {
	{
		conf := ConfigFake{
			p: func(s string, c int) {
				if c == 0 {
					assert.Equal(t, "hello\n", s)
				}
				if c == 1 {
					assert.Equal(t, "world\n", s)
				}
			},
			pg: func(s string, c int) {
				assert.Fail(t, "cannot pass here")
			},
			str: func(s string) {
				assert.Equal(t, "hello\nworld\n", s)
			},
		}
		goClip(NewFakeController([]string{"hello", "world"}, conf), "")
	}
	{
		conf := ConfigFake{
			p: func(s string, c int) {
				if c == 0 {
					assert.Equal(t, "hello\n", s)
				}
				if c == 1 {
					assert.Equal(t, "world\n", s)
				}
			},
			pg: func(s string, c int) {
				assert.Equal(t, "hello\n", s)
			},
			str: func(s string) {
				assert.Equal(t, "hello\n", s)
			},
		}
		goClip(NewFakeController([]string{"hello", "world"}, conf), "hello")
	}
	{
		conf := ConfigFake{
			p: func(s string, c int) {
				if c == 0 {
					assert.Equal(t, "hello\n", s)
				}
				if c == 1 {
					assert.Equal(t, "my\n", s)
				}
				if c == 2 {
					assert.Equal(t, "w0rld\n", s)
				}
			},
			pg: func(s string, c int) {
				assert.Equal(t, "w0rld\n", s)
			},
			str: func(s string) {
				assert.Equal(t, "w0rld\n", s)
			},
		}
		goClip(NewFakeController([]string{"hello", "my", "w0rld"}, conf), "[0-9]")
	}
}
