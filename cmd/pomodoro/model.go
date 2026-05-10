package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type phase int

const (
	phaseWork phase = iota
	phaseShortBreak
	phaseLongBreak
)

// largeur intérieure fixe du panneau (hors bordure et padding)
const panelWidth = 44

type tickMsg time.Time

// styles statiques (créés une fois)
var (
	styleDim   = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	styleKey   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EEEEEE")).Background(lipgloss.Color("#2A2A2A")).Padding(0, 1)
	styleDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	styleNotif = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)

	styleAccentWork  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E05C5C")).Bold(true)
	styleAccentShort = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC994")).Bold(true)
	styleAccentLong  = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EA8E0")).Bold(true)

	colorWork  = lipgloss.Color("#E05C5C")
	colorShort = lipgloss.Color("#4EC994")
	colorLong  = lipgloss.Color("#4EA8E0")

	shortcutsBar = lipgloss.JoinHorizontal(lipgloss.Top,
		styleKey.Render("ESPACE"), styleDesc.Render(" Start/Pause  "),
		styleKey.Render("S"), styleDesc.Render(" Session suiv.  "),
		styleKey.Render("R"), styleDesc.Render(" Reset  "),
		styleKey.Render("T"), styleDesc.Render(" Test son  "),
		styleKey.Render("Q"), styleDesc.Render(" Quitter"),
	)
)

type model struct {
	cfg           Config
	phase         phase
	remaining     time.Duration
	total         time.Duration
	running       bool
	pomodoroCount int
	progress      progress.Model
	width         int
	height        int
	notification  string
}

func newModel(cfg Config) model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithWidth(panelWidth-4),
	)
	total := time.Duration(cfg.WorkDuration) * time.Minute
	return model{
		cfg:       cfg,
		phase:     phaseWork,
		remaining: total,
		total:     total,
		progress:  p,
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
			m.notification = "Grande pause ! Reposez-vous bien."
		} else {
			m.phase = phaseShortBreak
			m.total = time.Duration(m.cfg.ShortBreak) * time.Minute
			m.notification = "Pause courte méritée !"
		}
	case phaseShortBreak, phaseLongBreak:
		m.phase = phaseWork
		m.total = time.Duration(m.cfg.WorkDuration) * time.Minute
		m.notification = "Au travail !"
	}
	m.remaining = m.total
}

func (m model) autoStart() bool {
	if m.phase == phaseWork {
		return m.cfg.AutoStartPomodoros
	}
	return m.cfg.AutoStartBreaks
}

func (m model) accent() lipgloss.Style {
	switch m.phase {
	case phaseShortBreak:
		return styleAccentShort
	case phaseLongBreak:
		return styleAccentLong
	}
	return styleAccentWork
}

func (m model) color() lipgloss.Color {
	switch m.phase {
	case phaseShortBreak:
		return colorShort
	case phaseLongBreak:
		return colorLong
	}
	return colorWork
}

func (m model) phaseLabel() string {
	switch m.phase {
	case phaseWork:
		return m.cfg.WorkLabel
	case phaseShortBreak:
		return m.cfg.ShortBreakLabel
	case phaseLongBreak:
		return m.cfg.LongBreakLabel
	}
	return ""
}

// ── View ────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {
		return "Initialisation..."
	}

	ac := m.accent()
	col := m.color()

	cp := func(s string) string {
		return lipgloss.PlaceHorizontal(panelWidth, lipgloss.Center, s)
	}

	// ── Titre ────────────────────────────────────────────────────────────────
	title := ac.Render("🍅  POMODORO TIMER")

	// ── Phase + dots ─────────────────────────────────────────────────────────
	phaseIcon := "⚡"
	if m.phase != phaseWork {
		phaseIcon = "☕"
	}
	phaseLine := ac.Render(phaseIcon + "  " + strings.ToUpper(m.phaseLabel()))

	total := m.cfg.PomodorosBeforeLongBreak
	done := m.pomodoroCount % total
	var dotsBuf strings.Builder
	for i := 0; i < total; i++ {
		if i > 0 {
			dotsBuf.WriteString("  ")
		}
		if i < done {
			dotsBuf.WriteString(ac.Render("◆"))
		} else {
			dotsBuf.WriteString(styleDim.Render("◇"))
		}
	}
	cycleStr := styleDim.Render(fmt.Sprintf("   cycle %d", (m.pomodoroCount/total)+1))
	dotsLine := dotsBuf.String() + cycleStr

	// ── Horloge ──────────────────────────────────────────────────────────────
	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	bigRows := renderBigTime(fmt.Sprintf("%02d:%02d", mins, secs))

	// ── Barre de progression ──────────────────────────────────────────────────
	var pct float64
	if m.total > 0 {
		pct = float64(m.total-m.remaining) / float64(m.total)
	}
	progressBar := m.progress.ViewAs(pct)

	// ── Statut ────────────────────────────────────────────────────────────────
	var statusLine string
	if m.running {
		statusLine = ac.Render("▶  En cours")
	} else {
		statusLine = styleDim.Render("⏸  En pause")
	}

	// ── Info prochaine pause ──────────────────────────────────────────────────
	nextInfo := ""
	if m.phase == phaseWork {
		if m.pomodoroCount > 0 && m.pomodoroCount%total == 0 {
			nextInfo = styleDim.Render(fmt.Sprintf("→  grande pause (%dm)", m.cfg.LongBreak))
		} else {
			left := total - (m.pomodoroCount % total)
			nextInfo = styleDim.Render(fmt.Sprintf("→  grande pause dans %d pomodoro(s)", left))
		}
	}

	// ── Notification ─────────────────────────────────────────────────────────
	notif := ""
	if m.notification != "" {
		notif = styleNotif.Render("✦  " + m.notification)
	}

	// ── Contenu du panneau ────────────────────────────────────────────────────
	divider := styleDim.Render(strings.Repeat("╌", panelWidth-2))

	rows := []string{
		cp(title),
		cp(divider),
		"",
		cp(phaseLine),
		cp(dotsLine),
		"",
	}
	for _, row := range bigRows {
		rows = append(rows, cp(ac.Render(row)))
	}
	rows = append(rows, "")
	rows = append(rows, cp(progressBar))
	rows = append(rows, "")
	rows = append(rows, cp(statusLine))
	if nextInfo != "" {
		rows = append(rows, cp(nextInfo))
	}
	if notif != "" {
		rows = append(rows, "")
		rows = append(rows, cp(notif))
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(col).
		Padding(1, 2).
		Width(panelWidth).
		Render(strings.Join(rows, "\n"))

	// ── Assemblage final ──────────────────────────────────────────────────────
	block := lipgloss.JoinVertical(lipgloss.Center,
		panel,
		"",
		shortcutsBar,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
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
