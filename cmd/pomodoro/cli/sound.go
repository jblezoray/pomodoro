package cli

import (
	_ "embed"
	"os"
	"os/exec"
	"runtime"

	"pomodoro/cmd/pomodoro/model"
)

//go:embed sounds/done.wav
var embeddedSound []byte

func beep(cfg model.Config) {
	go playSound(cfg)
}

func playSound(cfg model.Config) {
	if cfg.SoundFile != "" {
		if playFile(cfg.SoundFile) {
			return
		}
	}
	if playEmbedded() {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		for _, f := range []string{
			"/System/Library/Sounds/Glass.aiff",
			"/System/Library/Sounds/Ping.aiff",
			"/System/Library/Sounds/Tink.aiff",
		} {
			if exec.Command("afplay", f).Run() == nil {
				return
			}
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

func playEmbedded() bool {
	tmp, err := os.CreateTemp("", "pomodoro-*.aiff")
	if err != nil {
		return false
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(embeddedSound); err != nil {
		tmp.Close()
		return false
	}
	tmp.Close()
	return playFile(tmp.Name())
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
