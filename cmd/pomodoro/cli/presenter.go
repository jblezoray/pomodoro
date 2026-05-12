package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pomodoro/cmd/pomodoro/model"
)

type Presenter struct {
	m              *model.AppModel
	phaseStartedAt time.Time
}

func New(m *model.AppModel) *Presenter {
	return &Presenter{m: m}
}

func (p *Presenter) Run() error {
	fmt.Println("🍅 Dead simple Pomodoro")
	fmt.Println("Commands: [p/<enter>] pause/resume  [s] skip  [r] reset  [t] test sound  [q/ctrl-c] quit")
	fmt.Println()

	tickInterval := time.Second
	for _, arg := range os.Args[1:] {
		if arg == "--fast" {
			tickInterval = time.Second / 60
			break
		}
	}
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	input := readStdin()
	p.printLine()

	for {
		select {
		case <-ticker.C:
			p.onTick()
		case line, ok := <-input:
			if !ok {
				return nil
			}
			if quit := p.onInput(line); quit {
				return nil
			}
		}
	}
}

func (p *Presenter) onTick() {
	if !p.m.Running() {
		return
	}
	if p.phaseStartedAt.IsZero() {
		p.phaseStartedAt = time.Now()
	}
	if p.m.Tick() {
		p.clearLine()
		fmt.Printf("✓  %s\n", p.m.Notification())
		if p.m.Cfg().SoundEnabled {
			beep(p.m.Cfg())
		}
		p.phaseStartedAt = time.Time{}
		if p.m.AutoStart() {
			p.m.Start()
			p.phaseStartedAt = time.Now()
			p.printPhaseStart()
		}
	}
	p.printLine()
}

func (p *Presenter) onInput(line string) (quit bool) {
	p.clearLine()
	switch strings.TrimSpace(strings.ToLower(line)) {
	case "q", "quit":
		fmt.Println("Bye.")
		return true
	case "p", "":
		freshStart := !p.m.Running() && p.m.Remaining() == p.m.Total()
		if p.m.TogglePause() {
			if freshStart {
				p.phaseStartedAt = time.Now()
				if p.m.Cfg().SoundEnabled {
					beepStart(p.m.Cfg())
				}
				p.printPhaseStart()
			}
		}
	case "s":
		p.m.Skip()
		p.phaseStartedAt = time.Now()
		if p.m.Cfg().SoundEnabled {
			beepStart(p.m.Cfg())
		}
		p.printPhaseStart()
	case "r":
		p.m.Reset()
		p.phaseStartedAt = time.Time{}
	case "t":
		if p.m.Cfg().SoundEnabled {
			beep(p.m.Cfg())
		}
	}
	p.printLine()
	return false
}

func readStdin() <-chan string {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return ch
}

func (p *Presenter) printPhaseStart() {
	endAt := time.Now().Add(p.m.Remaining())
	fmt.Printf("▶ [%s] %s→%s\n", phaseLabel(p.m.Phase()), p.phaseStartedAt.Format("15:04"), endAt.Format("15:04"))
}

func (p *Presenter) printLine() {
	mins := int(p.m.Remaining().Minutes())
	secs := int(p.m.Remaining().Seconds()) % 60
	icon := "⏸"
	if p.m.Running() {
		icon = "▶"
	}
	timeRange := ""
	if !p.phaseStartedAt.IsZero() {
		endAt := time.Now().Add(p.m.Remaining())
		timeRange = fmt.Sprintf(" %s→%s", p.phaseStartedAt.Format("15:04"), endAt.Format("15:04"))
	}
	fmt.Printf("\r[%s] %02d:%02d %s%s  > ", phaseLabel(p.m.Phase()), mins, secs, icon, timeRange)
}

func (p *Presenter) clearLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", 80))
}

func phaseLabel(p model.Phase) string {
	switch p {
	case model.PhaseShortBreak:
		return "Short break"
	case model.PhaseLongBreak:
		return "Long break"
	}
	return "Work"
}
