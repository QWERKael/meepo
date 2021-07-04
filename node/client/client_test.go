package client

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"testing"
)

func TestParse(t *testing.T) {
	ch := make(chan []byte, 1)

	ch <- []byte("abc")
	close(ch)
	for {
		if c, ok := <-ch; ok == true {
			fmt.Printf("|| %s ||\n", c)
		} else {
			fmt.Printf("没有了\n")
			break
		}
	}


	cache, err := lru.New(2)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	cache.Add("1", "a")
	cache.Add("2", "b")
	cache.Add("3", "c")
	cache.Add("4", "d")
	fmt.Println(cache.Get("1"))
	fmt.Println(cache.Get("2"))
	fmt.Println(cache.Get("3"))
	fmt.Println(cache.Get("4"))
	cache.Add("4", "d1")
	fmt.Println(cache.Get("4"))
	cache.Add("4", "d2")
	fmt.Println(cache.Get("4"))


}
