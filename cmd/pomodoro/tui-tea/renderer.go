package tuitea

import "time"

// Phase represents the current timer phase.
type Phase int

const (
	PhaseWork Phase = iota
	PhaseShortBreak
	PhaseLongBreak
)

// State is a read-only snapshot of the model passed to the renderer.
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
