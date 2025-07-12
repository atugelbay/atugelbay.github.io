package main

import (
	"ascii-art/modules"
	"fmt"
	"os"
	"strings"
)

const bannerFile = "standard.txt"

func main() {
	replacer := strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
		`\\`, `\`,
	)

	inputRaw := strings.Join(os.Args[1:], " ")
	input := replacer.Replace(inputRaw)

	banners, err := modules.LoadBanner(bannerFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки баннера: %v\n", err)
		os.Exit(1)
	}

	modules.PrintBanner(input, banners)
}
