# pomodoro

A dead-simple Pomodoro timer for the terminal.

- Work / short break / long break cycle
- Auto-start between phases (configurable)
- Embedded notification sound (no external files needed)
- Custom sound file override
- Zero runtime dependencies — pure Go stdlib

## Install

**Requirements:** Go 1.21+

```sh
git clone https://github.com/yourname/pomodoro
cd pomodoro
make build          # produces bin/pomodoro
```

The first `make build` also runs `make sound` automatically to generate and embed the notification chime.

## Usage

```sh
./bin/pomodoro
```

| Key | Action |
|-----|--------|
| `p` / `Enter` | Pause / resume |
| `s` | Skip to next phase |
| `r` | Reset current phase |
| `t` | Test notification sound |
| `q` / `Ctrl-C` | Quit |

The status line shows the current phase, time remaining, and projected start → end times:

```
[Work session] 24:35 ▶ 09:30→09:55  >
```

## Configuration

Drop a `pomodoro.json` in the working directory (or `~/.pomodoro.json`, or `~/.config/pomodoro/config.json`):

```json
{
  "work_duration": 25,
  "short_break": 5,
  "long_break": 15,
  "pomodoros_before_long_break": 4,
  "auto_start_breaks": true,
  "auto_start_pomodoros": true,
  "sound_enabled": true,
  "sound_file": "",
  "work_label": "Work session",
  "short_break_label": "Small break",
  "long_break_label": "Long break"
}
```

Set `sound_file` to an absolute path to override the embedded chime with your own AIFF or WAV file.

## Make targets

| Target | Description |
|--------|-------------|
| `make build` | Build `bin/pomodoro` (generates sound if needed) |
| `make sound` | Regenerate the embedded notification chime |
| `make run` | `go run` without producing a binary |
| `make linux` | Cross-compile for Linux amd64 |
| `make darwin-arm` | Cross-compile for macOS arm64 |
| `make darwin-amd` | Cross-compile for macOS amd64 |
| `make clean` | Remove `bin/` and the generated sound asset |

## License

MIT — see [LICENSE](LICENSE).
