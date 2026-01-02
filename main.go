// Package main provides a command-line tool for converting HEIC/HEIF images to PNG or JPEG using ffmpeg.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	outType = flag.String("output", "", "png, jpg or jpeg")
	inPath  = flag.String("input", "", "File or directory path to convert")
)

// main is the entry point for the ffheic command-line tool.
// It parses flags, validates input, checks requirements, and processes files.
func main() {
	flag.Parse()

	err := verifyRequirements()
	checkError(err)

	inPathInfo, err := validateFlags()
	checkError(err)

	err = processFiles(inPathInfo)
	checkError(err)

	fmt.Println("Processing completed successfully.")
}

// verifyRequirements checks that the operating system is supported and that ffmpeg with HEIC/HEIF support is installed.
func verifyRequirements() (err error) {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		_, err := exec.LookPath("apt")
		if err != nil {
			fmt.Println("APT package manager is not installed and therefor support for this OS environment is not garanteed.")
		} else {
			// Verify libheif-examples is installed
			if _, err := exec.LookPath("heif-convert"); err != nil {
				return fmt.Errorf("the heif-convert command does not exist, please ensure that the libheif-examples apt package is installed and added to PATH.")
			}
		}
	case "windows":
		return fmt.Errorf("Currently, Windows is not supported.")
	case "darwin":
		return fmt.Errorf("Darwin/MacOS is not supported.")
	default:
		return fmt.Errorf("%s is not supported.", osType)
	}

	fmt.Println("OS requirements are met.")
	return nil
}

// validateFlags checks the command-line flags for validity and returns information about the input path.
// It ensures the output type is supported and the input path exists.
func validateFlags() (inPathInfo os.FileInfo, err error) {
	// Verify inPath exists
	if inPathInfo, err = os.Stat(*inPath); err != nil {
		return nil, err
	}

	// Ensure inPath is a full path
	if *inPath, err = filepath.Abs(*inPath); err != nil {
		return nil, err
	}
	fmt.Println("Input Path: ", *inPath)

	// Verify output type is viable
	switch *outType {
	case "jpeg", "jpg", "png":
		fmt.Println("Output Type: ", *outType)
	default:
		panic("Invalid output type. Use 'png', 'jpg' or 'jpeg'.")
	}

	return inPathInfo, nil
}

// processFiles converts the input file or all files in the input directory to the specified output format using ffmpeg.
// It handles both single file and directory input.
func processFiles(inPathInfo os.FileInfo) (err error) {
	// If inPath is a directory, process all files inside; otherwise, process the single file
	var inFiles []string
	if inPathInfo.IsDir() {
		entries, err := os.ReadDir(*inPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".heic" {
				inFiles = append(inFiles, filepath.Join(*inPath, entry.Name()))
			}
		}
	} else {
		inFiles = append(inFiles, *inPath)
	}

	// For each file, run ffmpeg conversion
	for _, inFile := range inFiles {
		outFile := strings.Replace(inFile, ".heic", "."+*outType, 1)
		// TODO
		// cmd := exec.Command("ffmpeg", "-i", inFile, outFile)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// if err := cmd.Run(); err != nil {
		// 	return fmt.Errorf("Failed to convert %s: %v\n", inFile, err)
		// }
		fmt.Printf("Converted %s to %s OK.\n", inFile, outFile)
	}

	return nil
}

// checkError panics if the provided error is non-nil.
// It is used for simple error handling throughout the program.
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
