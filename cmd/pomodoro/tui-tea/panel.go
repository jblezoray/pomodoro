package tuitea

import (
	"fmt"
	"strings"

	"pomodoro/cmd/pomodoro/model"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// fixed inner width of the panel (excluding border and padding)
const panelWidth = 44

// static styles (created once)
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
		styleKey.Render("SPACE"), styleDesc.Render(" Start/Pause  "),
		styleKey.Render("S"), styleDesc.Render(" Next session  "),
		styleKey.Render("R"), styleDesc.Render(" Reset  "),
		styleKey.Render("T"), styleDesc.Render(" Test sound  "),
		styleKey.Render("Q"), styleDesc.Render(" Quit"),
	)
)

// PanelRenderer renders a centered bordered panel.
type PanelRenderer struct {
	progress progress.Model
}

func NewPanelRenderer() *PanelRenderer {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
		progress.WithWidth(panelWidth-4),
	)
	return &PanelRenderer{progress: p}
}

func (r *PanelRenderer) accent(p model.Phase) lipgloss.Style {
	switch p {
	case model.PhaseShortBreak:
		return styleAccentShort
	case model.PhaseLongBreak:
		return styleAccentLong
	}
	return styleAccentWork
}

func (r *PanelRenderer) color(p model.Phase) lipgloss.Color {
	switch p {
	case model.PhaseShortBreak:
		return colorShort
	case model.PhaseLongBreak:
		return colorLong
	}
	return colorWork
}

func (r *PanelRenderer) phaseLabel(s model.State) string {
	switch s.Phase {
	case model.PhaseWork:
		return s.WorkLabel
	case model.PhaseShortBreak:
		return s.ShortBreakLabel
	case model.PhaseLongBreak:
		return s.LongBreakLabel
	}
	return ""
}

func (r *PanelRenderer) Render(s model.State) string {
	if s.Width == 0 {
		return "Loading..."
	}

	ac := r.accent(s.Phase)
	col := r.color(s.Phase)

	cp := func(str string) string {
		return lipgloss.PlaceHorizontal(panelWidth, lipgloss.Center, str)
	}

	// Title
	title := ac.Render("🍅  POMODORO TIMER")

	// Phase + dots
	phaseIcon := "⚡"
	if s.Phase != model.PhaseWork {
		phaseIcon = "☕"
	}
	phaseLine := ac.Render(phaseIcon + "  " + strings.ToUpper(r.phaseLabel(s)))

	total := s.PomodorosBeforeLongBreak
	done := s.PomodoroCount % total
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
	cycleStr := styleDim.Render(fmt.Sprintf("   cycle %d", (s.PomodoroCount/total)+1))
	dotsLine := dotsBuf.String() + cycleStr

	// Clock
	mins := int(s.Remaining.Minutes())
	secs := int(s.Remaining.Seconds()) % 60
	bigRows := renderBigTime(fmt.Sprintf("%02d:%02d", mins, secs))

	// Progress bar
	var pct float64
	if s.Total > 0 {
		pct = float64(s.Total-s.Remaining) / float64(s.Total)
	}
	progressBar := r.progress.ViewAs(pct)

	// Status
	var statusLine string
	if s.Running {
		statusLine = ac.Render("▶  Running")
	} else {
		statusLine = styleDim.Render("⏸  Paused")
	}

	// Next break info
	nextInfo := ""
	if s.Phase == model.PhaseWork {
		if s.PomodoroCount > 0 && s.PomodoroCount%total == 0 {
			nextInfo = styleDim.Render(fmt.Sprintf("→  long break (%dm)", s.LongBreak))
		} else {
			left := total - (s.PomodoroCount % total)
			nextInfo = styleDim.Render(fmt.Sprintf("→  long break in %d pomodoro(s)", left))
		}
	}

	// Notification
	notif := ""
	if s.Notification != "" {
		notif = styleNotif.Render("✦  " + s.Notification)
	}

	// Panel content
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

	// Final assembly
	block := lipgloss.JoinVertical(lipgloss.Center,
		panel,
		"",
		shortcutsBar,
	)

	return lipgloss.Place(s.Width, s.Height, lipgloss.Center, lipgloss.Center, block)
}
