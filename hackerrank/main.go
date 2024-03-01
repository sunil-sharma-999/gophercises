package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {

	fmt.Println("Choose hackerrank problem solver function:")
	fmt.Printf("1. CamelCase\n2. Caesar Cipher\n\n")

	option := 0
	_, err := fmt.Scan(&option)
	if err != nil {
		log.Fatal("Invalid input")
	}

	switch option {
	case 1:
		countWordsInCamelCase()
	case 2:
		caesarCipher()
	default:
		log.Fatal("Invalid input")

	}

}

func countWordsInCamelCase() {
	var inputStr string
	fmt.Print("\nText: ")
	fmt.Scan(&inputStr)
	camelToe := regexp.MustCompile("[A-Z]")
	inputStr = camelToe.ReplaceAllString(inputStr, "_$0")
	fmt.Printf("Word count: %v \n", len(strings.Split(inputStr, "_")))
}

func caesarCipher() {
	var inputStr string
	fmt.Print("\nText: ")
	fmt.Scan(&inputStr)

	decode := 1
	fmt.Printf("\n1. Encode\n2. Decode\n")
	_, err := fmt.Scan(&decode)
	if err != nil {
		log.Fatal(err)
	}

	lowerAAscii := 97
	lowerZAscii := 122
	capitalAAscii := 65
	capitalZAscii := 90

	inputRune := []rune(inputStr)

	for i, char := range inputStr {
		asciiChar := int(char)

		if asciiChar >= lowerAAscii && asciiChar <= lowerZAscii {

			if decode == 2 {
				if asciiChar-3 < lowerAAscii {
					asciiChar = (lowerZAscii - lowerAAscii) + (asciiChar - 3)
				} else {
					asciiChar -= 3
				}
			} else {
				if asciiChar+3 > lowerZAscii {
					asciiChar = ((asciiChar + 3) % lowerZAscii) + (lowerAAscii - 1)
				} else {
					asciiChar += 3
				}
			}

		}
		if asciiChar >= capitalAAscii && asciiChar <= capitalZAscii {
			if decode == 2 {
				if asciiChar-3 < capitalAAscii {
					asciiChar = (capitalZAscii - capitalAAscii) + (asciiChar - 3)
				} else {
					asciiChar -= 3
				}
			} else {
				if asciiChar+3 > capitalZAscii {
					asciiChar = ((asciiChar + 3) % capitalZAscii) + (capitalAAscii - 1)
				} else {
					asciiChar += 3
				}
			}

		}
		inputRune[i] = rune(asciiChar)
	}
	fmt.Printf("\nCiphered Text: %v\n", string(inputRune))
}
