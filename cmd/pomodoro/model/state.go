package model

import "time"

// State is a read-only snapshot passed to the renderer.
type State struct {
	Phase         Phase
	Remaining     time.Duration
	Total         time.Duration
	Running       bool
	PomodoroCount int
	Notification  string
	Width         int
	Height        int
	// config fields needed for rendering
	PomodorosBeforeLongBreak int
	LongBreak                int
	WorkLabel                string
	ShortBreakLabel          string
	LongBreakLabel           string
}

// Renderer renders a State to a terminal string.
type Renderer interface {
	Render(State) string
}

// Presenter runs the application UI loop.
type Presenter interface {
	Run() error
}
