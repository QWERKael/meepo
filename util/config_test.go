package util

import (
	"fmt"
	"log"
	"testing"
	"utility-go/config"
)

func TestParserFromByte(t *testing.T) {
	var configText = `
host: 127.0.0.1
listen: 4001
`
	conf := Conf{}
	err := config.ParserFromByte([]byte(configText), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%#v\n\n", conf)
}