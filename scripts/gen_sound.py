#!/usr/bin/env python3
"""Generates sounds/done.aiff: ascending C-E-G arpeggio resolving to a chord."""

import array
import math
import os
import subprocess
import tempfile
import wave

RATE = 44100
DURATION = 3.2
STEP = 0.42
CHORD_T0 = STEP * 3


def bell(t, freq, t0, decay=5.0):
    if t < t0:
        return 0.0
    dt = t - t0
    return math.exp(-decay * dt) * (
        math.sin(2 * math.pi * freq * dt) +
        0.3 * math.sin(2 * math.pi * freq * 2 * dt)
    )


NOTES = [
    (523.25, 0),           # C5
    (659.25, STEP),        # E5
    (783.99, STEP * 2),    # G5
    (523.25, CHORD_T0),    # chord
    (659.25, CHORD_T0),
    (783.99, CHORD_T0),
]

raw = [
    sum(bell(i / RATE, freq, t0) for freq, t0 in NOTES)
    for i in range(int(RATE * DURATION))
]

scale = 0.90 / max(abs(v) for v in raw)
samples = [int(max(-32767, min(32767, v * scale * 32767))) for v in raw]

os.makedirs("cmd/pomodoro/cli/sounds", exist_ok=True)

with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp:
    wav_path = tmp.name

with wave.open(wav_path, "wb") as w:
    w.setnchannels(1)
    w.setsampwidth(2)
    w.setframerate(RATE)
    w.writeframes(array.array("h", samples).tobytes())

subprocess.run(["ffmpeg", "-y", "-i", wav_path, "cmd/pomodoro/cli/sounds/done.aiff"], check=True)
os.unlink(wav_path)
print("cmd/pomodoro/cli/sounds/done.aiff generated")
