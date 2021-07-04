package main

import (
	"fmt"
	"testing"
)

func TestNet(t *testing.T) {

	if rst, err := Net(nil, nil); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%s", rst)
	}
}
