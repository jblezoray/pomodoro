package model

import "time"

type Phase int

const (
	PhaseWork Phase = iota
	PhaseShortBreak
	PhaseLongBreak
)

type AppModel struct {
	cfg           Config
	phase         Phase
	remaining     time.Duration
	total         time.Duration
	running       bool
	pomodoroCount int
	notification  string
}

func New(cfg Config) *AppModel {
	total := time.Duration(cfg.WorkDuration) * time.Minute
	return &AppModel{cfg: cfg, phase: PhaseWork, remaining: total, total: total}
}

// Tick advances the timer by one second.
// Sets running=false if the phase ends; returns true in that case.
func (m *AppModel) Tick() bool {
	m.remaining -= time.Second
	if m.remaining <= 0 {
		m.remaining = 0
		m.advancePhase()
		m.running = false
		return true
	}
	return false
}

// TogglePause flips the running state. Returns true if now running.
func (m *AppModel) TogglePause() bool {
	m.notification = ""
	m.running = !m.running
	return m.running
}

// Skip jumps to the next phase. Returns true if a new tick is needed.
func (m *AppModel) Skip() bool {
	m.notification = ""
	wasRunning := m.running
	m.advancePhase()
	m.running = wasRunning || m.AutoStart()
	return !wasRunning && m.running
}

func (m *AppModel) Reset() {
	m.remaining = m.total
	m.notification = ""
}

// Start sets the model to running (used after an auto-start decision).
func (m *AppModel) Start() { m.running = true }

func (m *AppModel) AutoStart() bool {
	if m.phase == PhaseWork {
		return m.cfg.AutoStartPomodoros
	}
	return m.cfg.AutoStartBreaks
}

func (m *AppModel) advancePhase() {
	switch m.phase {
	case PhaseWork:
		m.pomodoroCount++
		if m.pomodoroCount%m.cfg.PomodorosBeforeLongBreak == 0 {
			m.phase = PhaseLongBreak
			m.total = time.Duration(m.cfg.LongBreak) * time.Minute
			m.notification = "Long break! Rest well."
		} else {
			m.phase = PhaseShortBreak
			m.total = time.Duration(m.cfg.ShortBreak) * time.Minute
			m.notification = "Short break, you earned it!"
		}
	case PhaseShortBreak, PhaseLongBreak:
		m.phase = PhaseWork
		m.total = time.Duration(m.cfg.WorkDuration) * time.Minute
		m.notification = "Back to work!"
	}
	m.remaining = m.total
}

// State accessors
func (m *AppModel) Phase() Phase             { return m.phase }
func (m *AppModel) Remaining() time.Duration { return m.remaining }
func (m *AppModel) Total() time.Duration     { return m.total }
func (m *AppModel) Running() bool            { return m.running }
func (m *AppModel) PomodoroCount() int       { return m.pomodoroCount }
func (m *AppModel) Notification() string     { return m.notification }
func (m *AppModel) Cfg() Config              { return m.cfg }
