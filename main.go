/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"dongel/cmd"
	"fmt"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There was an error:", err)
		}
	}()
	cmd.Execute()
}
