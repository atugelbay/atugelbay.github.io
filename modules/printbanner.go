package modules

import (
	"fmt"
	"strings"
)

func PrintBanner(s string, m map[rune][]string) {
	if s == "" {
		return
	}

	if allNewlines(s) {
		count := strings.Count(s, "\n")
		for i := 0; i < count; i++ {
			fmt.Println()
		}
		return
	}

	lines := strings.Split(s, "\n")
	const tabWidth = 4

	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}

		for row := 0; row < 8; row++ {
			for _, ch := range line {
				switch ch {
				case '\t':
					for i := 0; i < tabWidth; i++ {
						fmt.Print(m[' '][row], " ")
					}
				default:
					if block, ok := m[ch]; ok {
						fmt.Print(block[row], " ")
					} else {
						space := len(m[' '][0])
						fmt.Print(strings.Repeat(" ", space), " ")
					}
				}
			}
			fmt.Println()
		}
	}
}

func allNewlines(s string) bool {
	for _, ch := range s {
		if ch != '\n' {
			return false
		}
	}
	return true
}
