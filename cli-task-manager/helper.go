package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func printHelp(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}
	os.Exit(0)
}

func isHelpCommand(args []string) bool {
	for _, arg := range args {
		if arg == "--help" {
			return true
		}
	}
	return false
}
