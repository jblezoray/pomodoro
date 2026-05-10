package main

import (
	"os"

	"pomodoro/cmd/pomodoro/cli"
	"pomodoro/cmd/pomodoro/model"
	tuitea "pomodoro/cmd/pomodoro/tui-tea"
)

// newPresenter is the composition root: add flags here to swap the UI.
func newPresenter(cfg model.Config) model.Presenter {
	m := model.New(cfg)
	for _, arg := range os.Args[1:] {
		if arg == "--simple" {
			return cli.New(m)
		}
	}
	return tuitea.New(m)
}
