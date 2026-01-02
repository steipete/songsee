// Package dsp provides spectral analysis utilities.
package dsp

import "math"

// FFTInPlace computes the in-place FFT for length power-of-two slices.
func FFTInPlace(x []complex128) {
	n := len(x)
	if n <= 1 {
		return
	}

	// Bit-reversal permutation.
	j := 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j &= ^bit
		}
		j |= bit
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
	}

	for size := 2; size <= n; size <<= 1 {
		angle := -2 * math.Pi / float64(size)
		wlen := complex(math.Cos(angle), math.Sin(angle))
		for i := 0; i < n; i += size {
			w := complex(1, 0)
			for j := 0; j < size/2; j++ {
				u := x[i+j]
				v := w * x[i+j+size/2]
				x[i+j] = u + v
				x[i+j+size/2] = u - v
				w *= wlen
			}
		}
	}
}
