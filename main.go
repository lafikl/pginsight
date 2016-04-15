package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}
	command := strings.ToLower(os.Args[1])
	insight := NewInsight()
	switch command {
	case "index":
		insight.Index(os.Args[2:])
		break
	case "disk":
		insight.Disk(os.Args[2:])
		break
	case "cache":
		insight.Cache(os.Args[2:])
		break
	case "queries":
		insight.Queries(os.Args[2:])
		break
	}
}
