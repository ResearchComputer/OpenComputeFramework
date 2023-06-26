package server

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
)

func PrintWelcomeMessage() {
	myFigure := figure.NewFigure("Open Compute Framework", "", true)
	myFigure.Print()
	fmt.Println(">> Join Discord for Discussion: https://discord.gg/3BD3RzK2K2")
	fmt.Println(">> Documentation: https://ocf.autoai.org")
}
