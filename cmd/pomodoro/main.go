package main

import (
	"fmt"
	"os"

	"pomodoro/cmd/pomodoro/model"
)

func main() {
	if err := newPresenter(model.Load()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
