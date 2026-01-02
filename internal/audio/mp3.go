package audio

import (
	"bytes"
	"errors"
	"io"

	"github.com/hajimehoshi/go-mp3"
)

// DecodeMP3If tries to decode MP3 data, returning ok=false when not MP3.
func DecodeMP3If(r io.ReadSeeker) (Audio, bool, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return Audio{}, false, err
	}
	isMP3 := string(header[0:3]) == "ID3" || (header[0] == 0xFF && header[1]&0xE0 == 0xE0)
	_, _ = r.Seek(0, io.SeekStart)
	if !isMP3 {
		return Audio{}, false, nil
	}
	pcm, err := decodeMP3(r)
	if err != nil {
		return Audio{}, true, err
	}
	return pcm, true, nil
}

func decodeMP3(r io.Reader) (Audio, error) {
	dec, err := mp3.NewDecoder(r)
	if err != nil {
		return Audio{}, err
	}
	pcm, err := io.ReadAll(dec)
	if err != nil {
		return Audio{}, err
	}
	if len(pcm)%2 != 0 {
		return Audio{}, errors.New("mp3: odd pcm length")
	}

	channels := 1
	if len(pcm)%4 == 0 {
		channels = 2
	}
	frames := len(pcm) / (2 * channels)
	out := make([]float64, frames)

	buf := bytes.NewReader(pcm)
	for i := 0; i < frames; i++ {
		var sum float64
		for ch := 0; ch < channels; ch++ {
			var sample int16
			if err := binaryRead(buf, &sample); err != nil {
				return Audio{}, err
			}
			sum += float64(sample) / 32768.0
		}
		out[i] = sum / float64(channels)
	}

	return Audio{SampleRate: dec.SampleRate(), Samples: out}, nil
}

func binaryRead(r io.Reader, v *int16) error {
	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return err
	}
	*v = int16(b[0]) | int16(b[1])<<8
	return nil
}
