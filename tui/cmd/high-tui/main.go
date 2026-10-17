// Command high-tui draws the generated windows. It does not parse HTML.
package main

import (
	"fmt"
	"os"

	"high-tui/internal/ui"
)

func main() {
	once := false
	for _, arg := range os.Args[1:] {
		if arg == "--once" {
			once = true
			continue
		}
		fmt.Fprintln(os.Stderr, "usage : high-tui [--once]")
		os.Exit(2)
	}
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if once {
		frame, err := ui.Screen(root, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println(frame)
		return
	}
	if err := ui.Run(root); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
