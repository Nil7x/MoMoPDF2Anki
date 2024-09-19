package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	pdf.DebugOn = true
	content, err := readPDF("your file path")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(content)

	word := parseContent(content)
	writeAnkiFile(word, "anki.txt")
}

func readPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		log.Println("open false")
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		log.Println("getPlainText false")
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func parseContent(content string) []map[string]string {
	var words []map[string]string

	for i, char := range content {
		if unicode.IsDigit(char) {
			for j := i + 1; j < len(content); j++ {
				if unicode.IsDigit(rune(content[j])) {
					if len(content[i:j]) < 2 {
						break
					}
					if strings.Contains(content[i:j], "墨墨") {
						break
					}

					card := content[i+1 : j]
					log.Println("match: ", card)

					word := make(map[string]string)
					bracket := 0
					speech := 0

					for c := 0; c < len(card); c++ {
						if card[c] == '[' {
							log.Println("word: ", card[:c])
							word["word"] = card[:c]
							bracket = c + 1
							continue
						}
						if card[c] == ']' && bracket != 0 {
							word["pronunciation"] = card[bracket:c]
							log.Println("pronunciation: ", card[bracket:c])
							speech = c + 1
							continue
						}
						if card[c] == '.' && speech != 0 {
							word["part_of_speech"] = card[speech:c]
							log.Println("part_of_speech: ", card[speech:c])

							log.Println("definition: ", card[c+2:])
							word["definition"] = card[c+2:]
							break
						}
					}
					words = append(words, word)
					break
				}
			}
		}
	}
	return words
}

func writeAnkiFile(words []map[string]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range words {
		front := fmt.Sprintf("%s[%s]%s.", word["word"], word["pronunciation"], word["part_of_speech"])
		back := word["definition"]
		line := fmt.Sprintf("%s\t%s\n", front, back)
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
