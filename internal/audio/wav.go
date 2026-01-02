package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// DecodeWAVIf tries to decode WAV data, returning ok=false when not WAV.
func DecodeWAVIf(r io.ReadSeeker) (Audio, bool, error) {
	header := make([]byte, 12)
	if _, err := io.ReadFull(r, header); err != nil {
		return Audio{}, false, err
	}
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		_, _ = r.Seek(0, io.SeekStart)
		return Audio{}, false, nil
	}
	_, _ = r.Seek(0, io.SeekStart)
	pcm, err := decodeWAV(r)
	if err != nil {
		return Audio{}, true, err
	}
	return pcm, true, nil
}

func decodeWAV(r io.ReadSeeker) (Audio, error) {
	var (
		fmtFound  bool
		dataFound bool
		fmtChunk  wavFormat
		data      []byte
	)

	header := make([]byte, 12)
	if _, err := io.ReadFull(r, header); err != nil {
		return Audio{}, err
	}
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return Audio{}, ErrUnsupported
	}

	for {
		chunkHeader := make([]byte, 8)
		_, err := io.ReadFull(r, chunkHeader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return Audio{}, err
		}
		chunkID := string(chunkHeader[0:4])
		chunkSize := int(binary.LittleEndian.Uint32(chunkHeader[4:8]))

		switch chunkID {
		case "fmt ":
			fmtFound = true
			buf := make([]byte, chunkSize)
			if _, err := io.ReadFull(r, buf); err != nil {
				return Audio{}, err
			}
			if err := parseWavFormat(buf, &fmtChunk); err != nil {
				return Audio{}, err
			}
		case "data":
			dataFound = true
			data = make([]byte, chunkSize)
			if _, err := io.ReadFull(r, data); err != nil {
				return Audio{}, err
			}
		default:
			// Skip unknown chunk.
			if _, err := r.Seek(int64(chunkSize), io.SeekCurrent); err != nil {
				return Audio{}, err
			}
		}
		if chunkSize%2 == 1 {
			_, _ = r.Seek(1, io.SeekCurrent)
		}
	}

	if !fmtFound || !dataFound {
		return Audio{}, errors.New("wav: missing fmt or data chunk")
	}
	return decodeWavData(fmtChunk, data)
}

type wavFormat struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	BitsPerSample uint16
	Extensible    bool
	SubFormat     [16]byte
}

func parseWavFormat(buf []byte, fmtChunk *wavFormat) error {
	if len(buf) < 16 {
		return errors.New("wav: short fmt chunk")
	}
	fmtChunk.AudioFormat = binary.LittleEndian.Uint16(buf[0:2])
	fmtChunk.NumChannels = binary.LittleEndian.Uint16(buf[2:4])
	fmtChunk.SampleRate = binary.LittleEndian.Uint32(buf[4:8])
	fmtChunk.BitsPerSample = binary.LittleEndian.Uint16(buf[14:16])
	if fmtChunk.AudioFormat == 0xFFFE && len(buf) >= 40 {
		fmtChunk.Extensible = true
		copy(fmtChunk.SubFormat[:], buf[24:40])
	}
	return nil
}

func decodeWavData(fmtChunk wavFormat, data []byte) (Audio, error) {
	format := fmtChunk.AudioFormat
	if fmtChunk.Extensible {
		// PCM subformat GUID 00000001-0000-0010-8000-00aa00389b71
		if isGUID(fmtChunk.SubFormat, 0x00000001) {
			format = 1
		} else if isGUID(fmtChunk.SubFormat, 0x00000003) {
			format = 3
		}
	}

	switch format {
	case 1, 3:
		// PCM or IEEE float.
	default:
		return Audio{}, fmt.Errorf("wav: unsupported format %d", format)
	}

	channels := int(fmtChunk.NumChannels)
	if channels < 1 {
		return Audio{}, errors.New("wav: invalid channel count")
	}

	sampleRate := int(fmtChunk.SampleRate)
	bits := int(fmtChunk.BitsPerSample)
	if bits == 0 {
		return Audio{}, errors.New("wav: invalid bits per sample")
	}

	var samples []float64
	if format == 3 {
		samples = decodeWavFloat(data, bits, channels)
	} else {
		samples = decodeWavPCM(data, bits, channels)
	}
	if samples == nil {
		return Audio{}, fmt.Errorf("wav: unsupported bit depth %d", bits)
	}

	return Audio{SampleRate: sampleRate, Samples: samples}, nil
}

func decodeWavPCM(data []byte, bits, channels int) []float64 {
	bytesPerSample := bits / 8
	frameSize := bytesPerSample * channels
	if frameSize == 0 {
		return nil
	}
	frames := len(data) / frameSize
	out := make([]float64, frames)
	idx := 0
	for i := 0; i < frames; i++ {
		var sum float64
		for ch := 0; ch < channels; ch++ {
			off := idx + ch*bytesPerSample
			var v int32
			switch bits {
			case 8:
				v = int32(int(data[off]) - 128)
			case 16:
				v = int32(int16(binary.LittleEndian.Uint16(data[off : off+2])))
			case 24:
				b := data[off : off+3]
				v = int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16
				if v&0x800000 != 0 {
					v |= ^0xffffff
				}
			case 32:
				v = int32(binary.LittleEndian.Uint32(data[off : off+4]))
			default:
				return nil
			}
			scale := float64(int64(1) << (bits - 1))
			sum += float64(v) / scale
		}
		out[i] = sum / float64(channels)
		idx += frameSize
	}
	return out
}

func decodeWavFloat(data []byte, bits, channels int) []float64 {
	bytesPerSample := bits / 8
	frameSize := bytesPerSample * channels
	if frameSize == 0 {
		return nil
	}
	frames := len(data) / frameSize
	out := make([]float64, frames)
	idx := 0
	for i := 0; i < frames; i++ {
		var sum float64
		for ch := 0; ch < channels; ch++ {
			off := idx + ch*bytesPerSample
			switch bits {
			case 32:
				sum += float64(math.Float32frombits(binary.LittleEndian.Uint32(data[off : off+4])))
			case 64:
				sum += math.Float64frombits(binary.LittleEndian.Uint64(data[off : off+8]))
			default:
				return nil
			}
		}
		out[i] = sum / float64(channels)
		idx += frameSize
	}
	return out
}

func isGUID(b [16]byte, sub uint32) bool {
	return binary.LittleEndian.Uint32(b[0:4]) == sub &&
		binary.LittleEndian.Uint16(b[4:6]) == 0x0000 &&
		binary.LittleEndian.Uint16(b[6:8]) == 0x0010 &&
		b[8] == 0x80 && b[9] == 0x00 &&
		b[10] == 0x00 && b[11] == 0xAA && b[12] == 0x00 && b[13] == 0x38 && b[14] == 0x9B && b[15] == 0x71
}
