package main

import (
	"aurox_task/cmd/cli"
	"fmt"
)

func main() {
	txt := fmt.Sprintf("%v", "main()")
	fmt.Println(txt)
	cli.Run()
}
