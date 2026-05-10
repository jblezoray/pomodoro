package main

import (
	"pomodoro/cmd/pomodoro/cli"
	"pomodoro/cmd/pomodoro/model"
)

func newPresenter(cfg model.Config) model.Presenter {
	return cli.New(model.New(cfg))
}
