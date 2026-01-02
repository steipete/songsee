package audio

import "fmt"

// Slice returns a time-based slice of audio in seconds.
func Slice(a Audio, startSec, durationSec float64) (Audio, error) {
	if startSec < 0 || durationSec < 0 {
		return Audio{}, fmt.Errorf("slice: start and duration must be >= 0")
	}
	if a.SampleRate <= 0 {
		return Audio{}, fmt.Errorf("slice: invalid sample rate")
	}
	if len(a.Samples) == 0 {
		return Audio{}, fmt.Errorf("slice: empty samples")
	}

	start := int(startSec * float64(a.SampleRate))
	if start >= len(a.Samples) {
		return Audio{}, fmt.Errorf("slice: start beyond end")
	}
	end := len(a.Samples)
	if durationSec > 0 {
		end = start + int(durationSec*float64(a.SampleRate))
		if end > len(a.Samples) {
			end = len(a.Samples)
		}
		if end <= start {
			return Audio{}, fmt.Errorf("slice: duration too short")
		}
	}

	return Audio{SampleRate: a.SampleRate, Samples: a.Samples[start:end]}, nil
}
