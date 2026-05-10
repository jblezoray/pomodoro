package main

import (
	"pomodoro/cmd/pomodoro/model"
	tuitea "pomodoro/cmd/pomodoro/tui-tea"
)

// newPresenter is the composition root: swap this to change the entire UI.
func newPresenter(cfg model.Config) model.Presenter {
	return tuitea.New(model.New(cfg))
}
