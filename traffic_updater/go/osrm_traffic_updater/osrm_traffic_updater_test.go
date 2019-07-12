package main

import (
	"testing"
)

func setChan2True(c chan<- bool) {
	c <- true
}

func setChan2False(c chan<- bool) {
	c <- false
}

func TestWait4AllPreconditions1(t *testing.T) {
	c1 := make(chan bool, 1)
	c2 := make(chan bool, 1)

	go setChan2True(c1)
	go setChan2True(c2)

	b1, b2 := wait4AllPreconditions(c1, c2)
	if !b1 || !b2 {
		t.Error("Testwait4AllPreconditions1 failed.\n")
	}
}

func TestWait4AllPreconditions2(t *testing.T) {
	c1 := make(chan bool, 1)
	c2 := make(chan bool, 1)

	go setChan2True(c1)
	go setChan2False(c2)

	b1, b2 := wait4AllPreconditions(c1, c2)
	if b1 && b2 {
		t.Error("Testwait4AllPreconditions1 failed.\n")
	}
}
