package modules

import (
	"strings"
)

// RenderBanner возвращает ASCII-арт для строки s, используя карту m.
func RenderBanner(s string, m map[rune][]string) string {
	var b strings.Builder

	// тот же реплейс, что и в main.go
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
		`\\`, `\`,
	)
	input := replacer.Replace(s)

	if input == "" {
		return ""
	}
	if allNewlines(input) {
		return strings.Repeat("\n", strings.Count(input, "\n"))
	}

	lines := strings.Split(input, "\n")
	const tabWidth, height = 4, 8

	for _, line := range lines {
		if line == "" {
			b.WriteByte('\n')
			continue
		}
		for row := 0; row < height; row++ {
			for _, ch := range line {
				switch ch {
				case '\t':
					for i := 0; i < tabWidth; i++ {
						b.WriteString(m[' '][row])
						b.WriteByte(' ')
					}
				default:
					if block, ok := m[ch]; ok {
						b.WriteString(block[row])
						b.WriteByte(' ')
					} else {
						space := len(m[' '][0])
						b.WriteString(strings.Repeat(" ", space))
						b.WriteByte(' ')
					}
				}
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}
