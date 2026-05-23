package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/imagemeta"
)

var sourceDir string
var destinationDir string

func main() {

	fmt.Print("\nLets move some photos!\n")
	fmt.Print("    📂  🚛  🏞️    \n\n")

	sourceDir = getSourceDir()
	destinationDir = getDestinationDir()
	fmt.Printf("source: %s\n", sourceDir)
	fmt.Printf("Destination: %s\n", destinationDir)

	// collect all images from source directory

	// create dictionary of [date-taken:[photo]]

	// for each image, get date-taken and insert into dictionary

	// for each key in dictionary, create folder in destination (if needed)
	// then for each key in dictionary, iterate though value array and copy photos into folder

}

func getDateTaken(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ex, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}

	dateTaken := ex.ExifIFD.DateTimeOriginal

	return fmt.Sprintf("%d-%d-%d", dateTaken.Year(), dateTaken.Month(), dateTaken.Day())
}

func getSourceDir() string {
	fmt.Println("Enter the source directory")
	sourceInput := collectInput()
	sourcePath, sPathErr := filepath.Abs(sourceInput)
	for !directoryExists(sourcePath) || sPathErr != nil {
		fmt.Println("Direcectory not found or isn't accessible, try again")
		sourceInput = collectInput()
		sourcePath, sPathErr = filepath.Abs(sourceInput)
	}

	return sourcePath
}

func getDestinationDir() string {
	fmt.Println("\nEnter the destination directory")
	destinationInput := collectInput()
	destinationPath, dPathErr := filepath.Abs(destinationInput)
	for !directoryExists(destinationPath) || dPathErr != nil {
		fmt.Println("Directory doesn't exist or isn't accessible, try again")
		destinationInput = collectInput()
		destinationPath, dPathErr = filepath.Abs(destinationInput)
	}

	return destinationPath
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
