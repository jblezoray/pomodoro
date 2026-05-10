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
	m *model.AppModel
}

func New(m *model.AppModel) *Presenter {
	return &Presenter{m: m}
}

func (p *Presenter) Run() error {
	fmt.Println("🍅 Pomodoro — simple mode")
	fmt.Println("Commands: [p] pause/resume  [s] skip  [r] reset  [t] sound  [q] quit")
	fmt.Println()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// input goroutine only reads stdin — all model access stays in the select loop
	input := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input <- scanner.Text()
		}
		close(input)
	}()

	p.printLine()

	for {
		select {
		case <-ticker.C:
			if !p.m.Running() {
				continue
			}
			if phaseEnded := p.m.Tick(); phaseEnded {
				p.clearLine()
				fmt.Printf("✓  %s\n", p.m.Notification())
				if p.m.Cfg().SoundEnabled {
					fmt.Print("\a")
				}
				if p.m.AutoStart() {
					p.m.Start()
				}
			}
			p.printLine()

		case line, ok := <-input:
			if !ok {
				return nil
			}
			p.clearLine()
			switch strings.TrimSpace(strings.ToLower(line)) {
			case "q", "quit":
				fmt.Println("Bye.")
				return nil
			case "p", "":
				p.m.TogglePause()
			case "s":
				p.m.Skip()
			case "r":
				p.m.Reset()
			case "t":
				if p.m.Cfg().SoundEnabled {
					fmt.Print("\a")
				}
			}
			p.printLine()
		}
	}
}

// printLine rewrites the current terminal line with status + prompt.
func (p *Presenter) printLine() {
	m := p.m
	mins := int(m.Remaining().Minutes())
	secs := int(m.Remaining().Seconds()) % 60
	icon := "⏸"
	if m.Running() {
		icon = "▶"
	}
	fmt.Printf("\r[%s] %02d:%02d %s  > ", phaseLabel(m.Phase()), mins, secs, icon)
}

// clearLine erases the current terminal line.
func (p *Presenter) clearLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", 60))
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
