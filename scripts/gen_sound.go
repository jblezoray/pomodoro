//go:build ignore

package main

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
)

const sampleRate = 44100

// Note frequencies in Hz — add more as needed.
const (
	C4 = 261.63
	D4 = 293.66
	E4 = 329.63
	F4 = 349.23
	G4 = 392.00
	A4 = 440.00
	B4 = 493.88
	C5 = 523.25
	D5 = 587.33
	E5 = 659.25
	F5 = 698.46
	G5 = 783.99
	A5 = 880.00
	B5 = 987.77
	C6 = 1046.50
)

// note is a frequency triggered at time t0 (seconds).
type note struct{ freq, t0 float64 }

// melody defines a sound to generate.
type melody struct {
	outPath  string
	duration float64 // seconds
	notes    []note
}

// melodies lists all sounds to generate.
// To add a new melody: append a new entry and run `go run scripts/gen_sound.go`.
var melodies = []melody{
	{
		outPath:  "cmd/pomodoro/cli/sounds/done.wav",
		duration: 2.0,
		notes: []note{
			{C5, 0.00},
			{E5, 0.15},
			{G5, 0.30},
			{E5, 0.45},
			{G5, 0.45},
			{C6, 0.45},
		},
	},
	{
		outPath:  "cmd/pomodoro/cli/sounds/start.wav",
		duration: 2.0,
		notes: []note{
			{G4, 0.0},
			{C4, 0.0},
			{E5, 0.0},
			{C6, 0.0},
			{C5, 0.10},
			{E5, 0.10},
			{G5, 0.10},
		},
	},
}

// bell is the synthesizer: an exponentially-decaying harmonic bell tone.
func bell(t, freq, t0 float64) float64 {
	if t < t0 {
		return 0
	}
	dt := t - t0
	return math.Exp(-5*dt) * (math.Sin(2*math.Pi*freq*dt) + 0.3*math.Sin(4*math.Pi*freq*dt))
}

// render synthesizes a melody into normalized 16-bit PCM samples.
func render(m melody) []int16 {
	n := int(sampleRate * m.duration)
	raw := make([]float64, n)
	peak := 0.0
	for i := range raw {
		t := float64(i) / sampleRate
		v := 0.0
		for _, nt := range m.notes {
			v += bell(t, nt.freq, nt.t0)
		}
		raw[i] = v
		if a := math.Abs(v); a > peak {
			peak = a
		}
	}
	scale := 0.90 / peak
	samples := make([]int16, n)
	for i, v := range raw {
		s := v * scale * 32767
		samples[i] = int16(math.Max(-32767, math.Min(32767, s)))
	}
	return samples
}

// writeWAV writes 16-bit mono PCM samples to a WAV file.
func writeWAV(path string, samples []int16) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dataSize := uint32(len(samples) * 2)
	write := func(v any) {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			panic(err)
		}
	}

	f.Write([]byte("RIFF"))
	write(uint32(36 + dataSize))
	f.Write([]byte("WAVE"))
	f.Write([]byte("fmt "))
	write(uint32(16))
	write(uint16(1))              // PCM
	write(uint16(1))              // mono
	write(uint32(sampleRate))
	write(uint32(sampleRate * 2)) // byte rate
	write(uint16(2))              // block align
	write(uint16(16))             // bits per sample
	f.Write([]byte("data"))
	write(dataSize)
	write(samples)
}

func main() {
	for _, m := range melodies {
		samples := render(m)
		writeWAV(m.outPath, samples)
		println("generated:", m.outPath)
	}
}
