package progressbar

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func printBar(fd int, rate int) {
	f := os.NewFile(uintptr(fd), "/dev/tty")

	termWidth, _, err := term.GetSize(fd)
	if err != nil {
		return
	}

	barWidth := termWidth - 7
	w := int(float64(rate) / 100.0 * float64(barWidth))
	if w < 0 {
		w = 0
	}
	if w > barWidth {
		w = barWidth
	}
	fmt.Fprintf(f, "\r[%s%s] %3d%%", strings.Repeat("=", w), strings.Repeat(" ", barWidth-w), rate)
}

func Display(fd int, x <-chan int) {
	printBar(fd, 0)
	for rate := range x {
		printBar(fd, rate)
	}
}
