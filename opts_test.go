package smsc

import (
	"fmt"
	"testing"
)

var ValidPanicTests = []struct {
	Hours, Minutes int
}{
	{-1, 0},
	{0, -1},
	{-1, -1},
	{0, 0},
	{0, 60},
	{23, 60},
	{24, 01},
}

func TestValid_panics(t *testing.T) {
	for _, tt := range ValidPanicTests {
		t.Run(fmt.Sprintf("%v:%v", tt.Hours, tt.Minutes), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatal()
				}
			}()
			Valid(tt.Hours, tt.Minutes)
		})
	}
}

var ValidStringTests = []struct {
	Valid valid
	S     string
}{
	{valid{0, 1}, "00:01"},
	{valid{0, 30}, "00:30"},
	{valid{12, 0}, "12:00"},
	{valid{23, 59}, "23:59"},
}

func TestValid_String(t *testing.T) {
	for _, tt := range ValidStringTests {
		if s := tt.Valid.String(); s != tt.S {
			t.Errorf("want %q, got %q", tt.S, s)
		}
	}
}
