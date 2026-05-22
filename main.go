package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Print("\nLets move some photos!\n\n")

	fmt.Println("Enter the source directory")
	source := collectInput()
	for !directoryExists(source) {
		fmt.Println("Source direcectory not found or isn't accessible, try again")
		source = collectInput()
	}

	fmt.Println("\nEnter the destination directory")
	destination := collectInput()
	for !directoryExists(destination) {
		fmt.Println("Destination directory doesn't exist or isn't accessible, try again")
		destination = collectInput()
	}

	fmt.Printf("source: %s\n", source)
	fmt.Printf("Destination: %s\n", destination)
}

func directoryExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func collectInput() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("→ ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		return text
	}
}
