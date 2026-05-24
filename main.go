package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/imagemeta"
	"github.com/urfave/cli/v3"
)

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".webp": true,
	".heic": true,
	".raw":  true,
	".svg":  true,
}

type ProgressState int

const (
	StateReady ProgressState = iota
	StateGathering
	StateCopying
	StateDone
)

var progressIndicator = map[ProgressState]string{
	StateReady:     "📂 ∙ ∙ ∙ ∙ ∙ ∙ ∙ 🦖 📷",
	StateGathering: "📂 ∙ ∙ ∙ ∙ ∙ ∙ 🦖 🏞️ 📷",
	StateCopying:   "📂 ∙ ∙ ∙ 🏞️ 🦖 ∙ ∙ ∙ 📷",
	StateDone:      "📂 🦖 ∙ ∙ ∙ ∙ ∙ ∙ ∙ 📷",
}

func main() {
	fmt.Print("\nLets move some photos!\n")
	fmt.Println(progressIndicator[StateReady])

	if len(os.Args) == 1 {
		sourceDir := getSourceDir("")
		destinationDir := getDestinationDir("")
		movePhotos(sourceDir, destinationDir)
	} else {
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "source",
					Aliases: []string{"s"},
					Value:   ".",
					Usage:   "directory where disorganized photos are located",
				},
				&cli.StringFlag{
					Name:    "destination",
					Aliases: []string{"d"},
					Value:   ".",
					Usage:   "directory where folders of sorted photos will be placed, photos will not overwrite existing files of the same name",
				},
			},
			Name:  "move",
			Usage: "Organize photos into folders by date taken with format yyyy-mm-dd",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				movePhotos(
					getSourceDir(cmd.String("source")),
					getDestinationDir(cmd.String("destination")),
				)
				return nil
			},
		}

		if err := cmd.Run(context.Background(), os.Args); err != nil {
			log.Fatal(err)
		}
	}

}

func movePhotos(source string, destination string) {
	fmt.Printf("\nGathering photos from /%s/...\n", filepath.Base(source))
	fmt.Println(progressIndicator[StateGathering])

	imageMap := createImageDateMap(source)
	if len(imageMap) == 0 {
		fmt.Println("\n[!] No photos found in source directory, aborting")
		return
	}
	fmt.Printf("\nCopying photos to /%s/...\n", filepath.Base(destination))
	fmt.Println(progressIndicator[StateCopying])
	imageCount, skippedCount, dirCount, newDirCount := writeImageMap(destination, imageMap)

	printSummary(imageCount, skippedCount, dirCount, newDirCount)
}

func printSummary(imageCount int, skippedCount int, dirCount int, newDirCount int) {
	if imageCount > 0 {
		fmt.Println("\nDone!")
		fmt.Println(progressIndicator[StateDone])
		var skippedString string
		if skippedCount > 0 {
			skippedString = fmt.Sprintf(" (%d existing photos skipped)", skippedCount)
		}
		fmt.Printf("Photos copied: %d%s", imageCount, skippedString)

		var newString string
		if dirCount-newDirCount > 0 {
			newString = fmt.Sprintf(" (%d directories already existed)", dirCount-newDirCount)
		}
		fmt.Printf("\nFolders created: %d%s\n", newDirCount, newString)
	} else {
		fmt.Printf("\n\nOpe, all source images are already present in destination, nothing to move!\n\n")
	}
}

func writeImageMap(destinationDir string, imageMap map[string][]string) (imageCount int, skippedCount, dirCount int, newDirCount int) {
	keys := make([]string, 0, dirCount)
	for k := range imageMap {
		keys = append(keys, k)
	}

	for _, dateTaken := range keys {
		if !directoryExists(filepath.Join(destinationDir, dateTaken)) {
			err := os.Mkdir(filepath.Join(destinationDir, dateTaken), 0755)
			if err != nil {
				panic(err)
			}
			newDirCount += 1
		}

		for _, imageSourcePath := range imageMap[dateTaken] {
			imageDestinationPath := filepath.Join(destinationDir, dateTaken, filepath.Base(imageSourcePath))

			sourceStat, err := os.Stat(imageSourcePath)
			if !sourceStat.Mode().IsRegular() || err != nil {
				fmt.Printf("\nfailed to touch source file: %s", filepath.Base(imageSourcePath))
				fmt.Println(err)
				continue
			}

			source, err := os.Open(imageSourcePath)
			if err != nil {
				fmt.Printf("\nfailed to touch source file: %s", filepath.Base(imageSourcePath))
				fmt.Println(err)
				continue
			}
			defer source.Close()

			destination, err := os.OpenFile(imageDestinationPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if errors.Is(err, os.ErrExist) {
				skippedCount++
				continue
			} else if err != nil {
				fmt.Printf("\nfailed to create destination for file: %s", filepath.Base(imageDestinationPath))
				fmt.Println(err)
				continue
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)

			if err != nil {
				fmt.Printf("\nFailed to copy file: %s", filepath.Base(imageSourcePath))
			}
			imageCount += 1
		}
	}

	return imageCount, skippedCount, len(imageMap), newDirCount
}

func createImageDateMap(directoryPath string) map[string][]string {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	imageMap := make(map[string][]string)
	for _, file := range files {
		imagePath := filepath.Join(directoryPath, file.Name())
		imgExt := strings.ToLower(filepath.Ext(imagePath))

		// if not image file, skip.
		if !imageExtensions[imgExt] {
			continue
		}

		dateTaken := getDateTaken(imagePath)

		imageMap[dateTaken] = append(imageMap[dateTaken], imagePath)
	}

	return imageMap
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

	dateTaken := ex.OriginalDate()

	if !dateTaken.IsZero() {
		return fmt.Sprintf("%d-%d-%d", dateTaken.Year(), dateTaken.Month(), dateTaken.Day())
	} else {
		return "unknown"
	}
}

func getSourceDir(sourceInput string) string {
	if sourceInput == "" {
		fmt.Println("\nEnter the source directory")
		sourceInput = collectInput()
	}
	sourcePath, sPathErr := filepath.Abs(sourceInput)
	for !directoryExists(sourcePath) || sPathErr != nil {
		fmt.Println("Direcectory not found or isn't accessible, try again")
		sourceInput = collectInput()
		sourcePath, sPathErr = filepath.Abs(sourceInput)
	}

	return sourcePath
}

func getDestinationDir(destinationInput string) string {
	if destinationInput == "" {
		fmt.Println("\nEnter the destination directory")
		destinationInput = collectInput()
	}
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
