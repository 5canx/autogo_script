package main

import (
	"fmt"
	"os"

	"app/proxy"
)

func main() {
	target := "京东"
	host := "baidu.com"

	if len(os.Args) > 1 {
		target = os.Args[1]
	}
	if len(os.Args) > 2 {
		host = os.Args[2]
	}

	success := proxy.Run(target, host)
	fmt.Println(success)
}
