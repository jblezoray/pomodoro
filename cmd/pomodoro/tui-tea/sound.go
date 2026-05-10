package tuitea

import (
	"fmt"
	"os/exec"
	"runtime"

	"pomodoro/cmd/pomodoro/model"
)

// beep plays a notification sound.
// If cfg.SoundFile is set, it tries to play it via the native player.
// Otherwise it falls back to the system default sound, then the terminal bell.
func beep(cfg model.Config) {
	if cfg.SoundFile != "" {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("afplay", cfg.SoundFile)
		case "linux":
			// try paplay then aplay
			if err := exec.Command("paplay", cfg.SoundFile).Start(); err == nil {
				return
			}
			cmd = exec.Command("aplay", cfg.SoundFile)
		}
		if cmd != nil {
			if err := cmd.Start(); err == nil {
				return
			}
		}
	}

	// system default sound
	switch runtime.GOOS {
	case "darwin":
		if err := exec.Command("afplay", "/System/Library/Sounds/Blow.aiff").Start(); err == nil {
			return
		}
		if err := exec.Command("osascript", "-e", "beep").Start(); err == nil {
			return
		}
	case "linux":
		for _, args := range [][]string{
			{"paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga"},
			{"aplay", "/usr/share/sounds/alsa/Front_Center.wav"},
		} {
			if err := exec.Command(args[0], args[1:]...).Start(); err == nil {
				return
			}
		}
	}
	fmt.Print("\a")
}
