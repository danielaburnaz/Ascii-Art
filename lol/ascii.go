package lol

// Import required packages
import (
	"bufio"
	// "fmt"
	"log"
	"os"
	"regexp"
)

// setScanner repositions the file scanner to the appropriate position based on the ASCII character passed.
func setScanner(file *os.File, offsetMap map[int]int64, sign rune) {
	// Calculate the starting line for the given ASCII character in the file.
	calc_sign_start := (int(sign)-31)*9 - 7
	// Seek to the calculated position in the file.
	file.Seek(offsetMap[calc_sign_start], 0)
}

// printAscii prints the ASCII representation of the given text.
func printAscii(file *os.File, offsetMap map[int]int64, text_array []string) map[int]string {
	// Create a map to hold each line of the ASCII art.
	line := make(map[int]string)

	// Loop over each rune in the text string.
	for index, text := range text_array {
		for _, run := range text {
			// Initialize a new scanner.
			scanner := bufio.NewScanner(file)
			// Set the scanner's position based on the current rune.
			setScanner(file, offsetMap, run)
			// Scan and add the corresponding ASCII lines to the map.
			for i := index * 8; i < 8*(index+1); i++ {
				scanner.Scan()
				line[i] += scanner.Text()
			}
		}
		for i := index * 8; i < 8*(index+1); i++ {
			line[i] += "<br>"
		}
	}

	// Print the ASCII lines.
	// fmt.Println(line)
	return line
}

// check is a simple utility function to check for errors.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// lol
// main is the entry point of the program.
func AsciiArtTransform(text string, s string) map[int]string {
	// Retrieve command-line arguments and join them into a single string.

	// Check if the input text is just a newline or empty.
	if text == "\\n" {
		return make(map[int]string)
	} else if text == "" {
		return make(map[int]string)
	}

	files := map[string]bool{
		"standard.txt":   true,
		"shadow.txt":     true,
		"thinkertoy.txt": true,
	}

	if !files[s] {
		s = "standard.txt"
	}
	// Open the standard.txt file.
	file, err := os.Open(s)
	check(err)
	defer file.Close()

	// Initialize a scanner to read the file.
	scanner := bufio.NewScanner(file)

	// Initialize byte offset and line number variables.
	byteOffset := int64(0)
	lineNumber := 1
	// Create a map to hold line numbers and their corresponding byte offsets.
	offsetMap := make(map[int]int64)

	// Read the entire file line by line, populating the offsetMap.
	for scanner.Scan() {
		line := scanner.Text()
		offsetMap[lineNumber] = byteOffset
		byteOffset += int64(len(line) + 1)
		lineNumber++
	}

	// Split the input text by newline characters.
	result := regexp.MustCompile(`\n`).Split(text, -1)

	return printAscii(file, offsetMap, result)

}
