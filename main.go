package main

import (
	"fmt"
	"os"

	"git.progwebtech.com/osiloke/gotrack/cmd"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
