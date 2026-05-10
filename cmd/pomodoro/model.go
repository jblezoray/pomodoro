package main

import (
	"time"

	tuitea "pomodoro/cmd/pomodoro/tui-tea"

	tea "github.com/charmbracelet/bubbletea"
)

type phase int

const (
	phaseWork phase = iota
	phaseShortBreak
	phaseLongBreak
)

type tickMsg time.Time

type model struct {
	cfg           Config
	phase         phase
	remaining     time.Duration
	total         time.Duration
	running       bool
	pomodoroCount int
	width         int
	height        int
	notification  string
	renderer      tuitea.Renderer
}

func newModel(cfg Config, r tuitea.Renderer) model {
	total := time.Duration(cfg.WorkDuration) * time.Minute
	return model{
		cfg:       cfg,
		phase:     phaseWork,
		remaining: total,
		total:     total,
		renderer:  r,
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if !m.running {
			return m, nil
		}
		m.remaining -= time.Second
		if m.remaining <= 0 {
			m.remaining = 0
			if m.cfg.SoundEnabled {
				beep(m.cfg)
			}
			m.advancePhase()
			if m.autoStart() {
				m.running = true
				return m, doTick()
			}
			m.running = false
			return m, nil
		}
		return m, doTick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case " ":
			m.notification = ""
			if m.running {
				m.running = false
				return m, nil
			}
			m.running = true
			return m, doTick()

		case "s":
			m.notification = ""
			wasRunning := m.running
			m.advancePhase()
			if wasRunning || m.autoStart() {
				if !wasRunning {
					m.running = true
					return m, doTick()
				}
				m.running = true
			} else {
				m.running = false
			}
			return m, nil

		case "r":
			m.remaining = m.total
			m.notification = ""
			return m, nil

		case "t":
			if m.cfg.SoundEnabled {
				beep(m.cfg)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *model) advancePhase() {
	switch m.phase {
	case phaseWork:
		m.pomodoroCount++
		if m.pomodoroCount%m.cfg.PomodorosBeforeLongBreak == 0 {
			m.phase = phaseLongBreak
			m.total = time.Duration(m.cfg.LongBreak) * time.Minute
			m.notification = "Long break! Rest well."
		} else {
			m.phase = phaseShortBreak
			m.total = time.Duration(m.cfg.ShortBreak) * time.Minute
			m.notification = "Short break, you earned it!"
		}
	case phaseShortBreak, phaseLongBreak:
		m.phase = phaseWork
		m.total = time.Duration(m.cfg.WorkDuration) * time.Minute
		m.notification = "Back to work!"
	}
	m.remaining = m.total
}

func (m model) autoStart() bool {
	if m.phase == phaseWork {
		return m.cfg.AutoStartPomodoros
	}
	return m.cfg.AutoStartBreaks
}

func (m model) uiPhase() tuitea.Phase {
	switch m.phase {
	case phaseShortBreak:
		return tuitea.PhaseShortBreak
	case phaseLongBreak:
		return tuitea.PhaseLongBreak
	}
	return tuitea.PhaseWork
}

func (m model) state() tuitea.State {
	return tuitea.State{
		Phase:                    m.uiPhase(),
		Remaining:                m.remaining,
		Total:                    m.total,
		Running:                  m.running,
		PomodoroCount:            m.pomodoroCount,
		Notification:             m.notification,
		Width:                    m.width,
		Height:                   m.height,
		PomodorosBeforeLongBreak: m.cfg.PomodorosBeforeLongBreak,
		LongBreak:                m.cfg.LongBreak,
		WorkLabel:                m.cfg.WorkLabel,
		ShortBreakLabel:          m.cfg.ShortBreakLabel,
		LongBreakLabel:           m.cfg.LongBreakLabel,
	}
}

func (m model) View() string {
	return m.renderer.Render(m.state())
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
