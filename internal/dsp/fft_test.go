package dsp

import "testing"

func TestFFTImpulse(t *testing.T) {
	x := []complex128{1, 0, 0, 0}
	FFTInPlace(x)
	for i, v := range x {
		if real(v) < 0.99 || real(v) > 1.01 || imag(v) != 0 {
			t.Fatalf("bin %d = %v", i, v)
		}
	}
}
