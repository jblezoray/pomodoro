package cli

import (
	"os/exec"
	"runtime"

	"pomodoro/cmd/pomodoro/model"
)

func beep(cfg model.Config) {
	go playSound(cfg)
}

func playSound(cfg model.Config) {
	if cfg.SoundFile != "" {
		if playFile(cfg.SoundFile) {
			return
		}
	}

	switch runtime.GOOS {
	case "darwin":
		for _, f := range []string{
			"/System/Library/Sounds/Glass.aiff",
			"/System/Library/Sounds/Blow.aiff",
			"/System/Library/Sounds/Tink.aiff",
			"/System/Library/Sounds/Ping.aiff",
		} {
			if exec.Command("afplay", f).Run() == nil {
				return
			}
		}
		if exec.Command("osascript", "-e", "beep").Run() == nil {
			return
		}
	case "linux":
		for _, args := range [][]string{
			{"paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga"},
			{"aplay", "/usr/share/sounds/alsa/Front_Center.wav"},
		} {
			if exec.Command(args[0], args[1:]...).Run() == nil {
				return
			}
		}
	}
}

func playFile(path string) bool {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("afplay", path).Run() == nil
	case "linux":
		if exec.Command("paplay", path).Run() == nil {
			return true
		}
		return exec.Command("aplay", path).Run() == nil
	}
	return false
}
