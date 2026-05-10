//go:build ignore

package main

import (
	"encoding/binary"
	"math"
	"os"
)

const (
	rate     = 44100
	duration = 3.2
	step     = 0.42
)

type note struct{ freq, t0 float64 }

func bell(t, freq, t0 float64) float64 {
	if t < t0 {
		return 0
	}
	dt := t - t0
	return math.Exp(-5*dt) * (math.Sin(2*math.Pi*freq*dt) + 0.3*math.Sin(4*math.Pi*freq*dt))
}

func main() {
	chord := step * 3
	notes := []note{
		{523.25, 0},      // C5
		{659.25, step},   // E5
		{783.99, step*2}, // G5
		{523.25, chord},  // chord
		{659.25, chord},
		{783.99, chord},
	}

	n := int(rate * duration)
	raw := make([]float64, n)
	peak := 0.0
	for i := range raw {
		t := float64(i) / rate
		v := 0.0
		for _, nt := range notes {
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

	if err := os.MkdirAll("cmd/pomodoro/cli/sounds", 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create("cmd/pomodoro/cli/sounds/done.wav")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dataSize := uint32(n * 2)
	write := func(v any) {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			panic(err)
		}
	}

	// RIFF header
	f.Write([]byte("RIFF"))
	write(uint32(36 + dataSize))
	f.Write([]byte("WAVE"))
	// fmt chunk
	f.Write([]byte("fmt "))
	write(uint32(16))
	write(uint16(1))        // PCM
	write(uint16(1))        // mono
	write(uint32(rate))     // sample rate
	write(uint32(rate * 2)) // byte rate
	write(uint16(2))        // block align
	write(uint16(16))       // bits per sample
	// data chunk
	f.Write([]byte("data"))
	write(dataSize)
	write(samples)
}
