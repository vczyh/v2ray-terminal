package main

import (
	"fmt"
	"testing"
)

func TestParseVMESS(t *testing.T) {
	v:=vmess{
		v:  "111",
		ps: "pssss",
	}
	fmt.Println(v)
}
