package main

import (
	"fmt"
	"testing"
	"utility-go/codec"
)

func TestRun(t *testing.T) {
	j := []string{"a", "b"}
	s, _ := codec.EncodeJson(j)
	fmt.Println(string(s))
}
