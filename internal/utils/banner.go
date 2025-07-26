package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var colors = []string{
	"\033[31m", "\033[32m", "\033[33m", "\033[34m",
	"\033[35m", "\033[36m", "\033[91m", "\033[92m",
	"\033[93m", "\033[94m", "\033[95m",
}

const reset = "\033[0m"

const banner = `
  _______   _   _                _  __  _   _
 |__   __| (_) | |              | |/ / (_) | |
    | |     _  | |        ___   | ' /   _  | |_
    | |    | | | |       / _ \  |  <   | | | __|
    | |    | | | |____  | (_) | | . \  | | | |_
    |_|    |_| |______|  \___/  |_|\_\ |_|  \__|
`

func PrintBanner() {
	lines := splitLines(banner)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, line := range lines {
		color := colors[rand.Intn(len(colors))]
		fmt.Printf("%s%s%s\n", color, line, reset)
		time.Sleep(80 * time.Millisecond)
	}

	fmt.Println()
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
