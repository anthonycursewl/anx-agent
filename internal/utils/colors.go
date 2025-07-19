package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	AgentName     = color.New(color.FgHiCyan, color.Bold)
	UserInput     = color.New(color.FgHiGreen, color.Bold)
	AgentResponse = color.New(color.FgHiWhite)
	ErrorMsg      = color.New(color.FgHiRed)
	FileContent   = color.New(color.FgHiYellow)
	Command       = color.New(color.FgHiMagenta)
	Loading       = color.New(color.FgHiBlue, color.Bold)
)

type PixelFrame struct {
	Pixels [][]string
	Colors [][]*color.Color
}

func LoadingAnimation(done chan bool, message string) {
	frames := []PixelFrame{
		createPixelFrame(0),
		createPixelFrame(1),
		createPixelFrame(2),
		createPixelFrame(3),
	}

	i := 0
	frameCount := len(frames)

	const (
		clearLine     = "\033[2K\r"
		saveCursor    = "\033[s"
		restoreCursor = "\033[u"
	)

	fmt.Print(saveCursor + message + "\n")

	for {
		select {
		case <-done:
			fmt.Print(clearLine + restoreCursor + clearLine)
			return
		default:
			frame := frames[i]

			fmt.Print(restoreCursor + "\033[1B" + clearLine)
			printPixelFrame(frame)

			i = (i + 1) % frameCount
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func createPixelFrame(phase int) PixelFrame {
	pixels := [][]string{
		{"◢", "■", "◣"},
		{"■", " ", "■"},
		{"◥", "■", "◤"},
	}

	colors := make([][]*color.Color, 3)
	for i := range colors {
		colors[i] = make([]*color.Color, 3)
		for j := range colors[i] {
			c := color.New(color.FgHiBlack).Add(color.Bold)
			colors[i][j] = c
		}
	}

	rainbow := []*color.Color{
		color.New(color.FgHiCyan).Add(color.Bold),
		color.New(color.FgHiMagenta).Add(color.Bold),
		color.New(color.FgHiBlue).Add(color.Bold),
		color.New(color.FgHiGreen).Add(color.Bold),
	}

	switch phase {
	case 0:
		colors[0][0] = rainbow[0]
		colors[0][2] = rainbow[1]
		colors[2][0] = rainbow[2]
		colors[2][2] = rainbow[3]
	case 1:
		colors[0][1] = rainbow[1]
		colors[1][0] = rainbow[2]
		colors[1][2] = rainbow[0]
		colors[2][1] = rainbow[3]
	case 2:
		colors[1][1] = rainbow[0]
		colors[0][1] = rainbow[1]
		colors[1][0] = rainbow[2]
		colors[1][2] = rainbow[3]
	case 3:
		colors[0][0] = rainbow[0]
		colors[0][2] = rainbow[1]
		colors[2][0] = rainbow[2]
		colors[2][2] = rainbow[3]
		colors[1][1] = color.New(color.FgHiWhite).Add(color.Bold)
	}

	return PixelFrame{
		Pixels: pixels,
		Colors: colors,
	}
}

func printPixelFrame(frame PixelFrame) {
	fmt.Print("  ")
	for j := 0; j < 3; j++ {
		c := frame.Colors[0][j]
		fmt.Print(c.Sprint(frame.Pixels[0][j] + " "))
	}
	fmt.Println()

	fmt.Print("  ")
	for j := 0; j < 3; j++ {
		c := frame.Colors[1][j]
		fmt.Print(c.Sprint(frame.Pixels[1][j] + " "))
	}
	fmt.Println()

	fmt.Print("  ")
	for j := 0; j < 3; j++ {
		c := frame.Colors[2][j]
		fmt.Print(c.Sprint(frame.Pixels[2][j] + " "))
	}
}

func StartLoading(message string) func() {
	done := make(chan bool)
	go LoadingAnimation(done, message)
	return func() {
		time.Sleep(50 * time.Millisecond)
		done <- true
	}
}
