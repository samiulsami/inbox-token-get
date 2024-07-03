package main

import (
	"log"
	"testing"
)

func TestJmap(t *testing.T) {
	a, err := getJMAPToken()
	if err != nil {
		t.Error(err)
	}
	log.Println(a)
}
