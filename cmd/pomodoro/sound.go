package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// beep joue un son de notification.
// Si cfg.SoundFile est renseigné, on tente de le jouer via le lecteur natif.
// Sinon on utilise le son système par défaut, puis le bell terminal.
func beep(cfg Config) {
	if cfg.SoundFile != "" {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("afplay", cfg.SoundFile)
		case "linux":
			// tente paplay puis aplay
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

	// son système par défaut
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
