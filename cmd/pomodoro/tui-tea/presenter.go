package tuitea

import (
	"time"

	"pomodoro/cmd/pomodoro/model"

	tea "github.com/charmbracelet/bubbletea"
)

// Presenter runs the BubbleTea UI loop for an AppModel.
type Presenter struct {
	model *model.AppModel
}

func New(m *model.AppModel) *Presenter {
	return &Presenter{model: m}
}

func (p *Presenter) Run() error {
	prog := tea.NewProgram(
		teaAdapter{model: p.model, renderer: NewPanelRenderer()},
		tea.WithAltScreen(),
	)
	_, err := prog.Run()
	return err
}

// teaAdapter bridges AppModel to the tea.Model interface.
type teaAdapter struct {
	model    *model.AppModel
	renderer model.Renderer
	width    int
	height   int
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (a teaAdapter) Init() tea.Cmd { return nil }

func (a teaAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tickMsg:
		if !a.model.Running() {
			return a, nil
		}
		phaseEnded := a.model.Tick()
		if phaseEnded {
			if a.model.Cfg().SoundEnabled {
				beep(a.model.Cfg())
			}
			if a.model.AutoStart() {
				a.model.Start()
				return a, doTick()
			}
			return a, nil
		}
		return a, doTick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case " ":
			if a.model.TogglePause() {
				return a, doTick()
			}
			return a, nil
		case "s":
			if a.model.Skip() {
				return a, doTick()
			}
			return a, nil
		case "r":
			a.model.Reset()
			return a, nil
		case "t":
			if a.model.Cfg().SoundEnabled {
				beep(a.model.Cfg())
			}
			return a, nil
		}
	}

	return a, nil
}

func (a teaAdapter) View() string {
	cfg := a.model.Cfg()
	return a.renderer.Render(model.State{
		Phase:                    a.model.Phase(),
		Remaining:                a.model.Remaining(),
		Total:                    a.model.Total(),
		Running:                  a.model.Running(),
		PomodoroCount:            a.model.PomodoroCount(),
		Notification:             a.model.Notification(),
		Width:                    a.width,
		Height:                   a.height,
		PomodorosBeforeLongBreak: cfg.PomodorosBeforeLongBreak,
		LongBreak:                cfg.LongBreak,
		WorkLabel:                cfg.WorkLabel,
		ShortBreakLabel:          cfg.ShortBreakLabel,
		LongBreakLabel:           cfg.LongBreakLabel,
	})
}
