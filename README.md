# pomodoro

[![CI](https://github.com/jblezoray/pomodoro/actions/workflows/ci.yml/badge.svg)](https://github.com/jblezoray/pomodoro/actions/workflows/ci.yml)

A dead-simple Pomodoro timer for the terminal.

```
% ./bin/pomodoro
🍅 Dead simple Pomodoro
Commands: [p/<enter>] pause/resume  [s] skip  [r] reset  [t] test sound  [q/ctrl-c] quit

[Work] 25:00 ⏸  >
[Work] 24:51 ▶ 00:01→00:26  > s
[Short break] 04:55 ▶ 00:02→00:06  >
[Short break] 04:55 ⏸ 00:02→00:07  >
```

## Install

**Requirements:** Go 1.21+

```sh
git clone https://github.com/jblezoray/pomodoro
cd pomodoro
make build          # produces bin/pomodoro
```

## Usage

```sh
./bin/pomodoro
```

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

| Target            | Description                                      |
| ----------------- | ------------------------------------------------ |
| `make build`      | Build `bin/pomodoro` (generates sound if needed) |
| `make sound`      | Regenerate the embedded notification chime       |
| `make run`        | `go run` without producing a binary              |
| `make linux`      | Cross-compile for Linux amd64                    |
| `make darwin-arm` | Cross-compile for macOS arm64                    |
| `make darwin-amd` | Cross-compile for macOS amd64                    |
| `make clean`      | Remove `bin/` and the generated sound asset      |

## License

LGPL-3.0 license — see [LICENSE](LICENSE).
