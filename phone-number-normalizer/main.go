package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func normalizePhoneNumber(phoneNumStr string) string {
	reg := regexp.MustCompile(`\D`)
	return reg.ReplaceAllString(strings.TrimSpace(phoneNumStr), "")
}

func main() {
	path := flag.String("file", "", "Input files containing phone numbers")
	flag.Parse()

	if *path == "" {
		log.Fatal("Provide file")
	}

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(normalizePhoneNumber(scanner.Text()))
	}
}
