package modules

import (
	"bufio"
	"os"
)

func LoadBanner(filePath string) (map[rune][]string, error) {
	const asciiStart = 32

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var blocks [][]string
	var current []string
	// group by 8 lines
	for _, l := range lines {
		if l == "" {
			if len(current) == 8 {
				blocks = append(blocks, current)
			}
			current = nil
		} else {
			current = append(current, l)
		}
	}
	if len(current) == 8 {
		blocks = append(blocks, current)
	}

	bannerMap := make(map[rune][]string)
	for i, block := range blocks {
		bannerMap[rune(asciiStart+i)] = block
	}
	return bannerMap, nil
}
