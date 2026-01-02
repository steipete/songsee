package audio

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDecodeWAVPCM8(t *testing.T) {
	data := makeWAVPCM(8, []int32{0, 64, -64, 0}, 44100)
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) != 4 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeWAVPCM24(t *testing.T) {
	data := makeWAVPCM(24, []int32{0, 100000, -100000, 0}, 44100)
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) != 4 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeWAVPCM32(t *testing.T) {
	data := makeWAVPCM(32, []int32{0, 100000, -100000, 0}, 44100)
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) != 4 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeWAVUnsupportedFormat(t *testing.T) {
	data := makeWAVCustom(7, 16, []byte{0, 0}, 44100)
	if _, err := DecodeBytes(data, Options{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWAVUnsupportedBits(t *testing.T) {
	data := makeWAVCustom(1, 12, []byte{0, 0}, 44100)
	if _, err := DecodeBytes(data, Options{}); err == nil {
		t.Fatalf("expected error")
	}
}

func makeWAVPCM(bits int, samples []int32, sampleRate int) []byte {
	data := &bytes.Buffer{}
	for _, s := range samples {
		switch bits {
		case 8:
			b := byte(int(s) + 128)
			data.WriteByte(b)
		case 16:
			_ = binary.Write(data, binary.LittleEndian, int16(s))
		case 24:
			v := uint32(int32(s))
			data.WriteByte(byte(v))
			data.WriteByte(byte(v >> 8))
			data.WriteByte(byte(v >> 16))
		case 32:
			_ = binary.Write(data, binary.LittleEndian, int32(s))
		}
	}
	return makeWAVCustom(1, bits, data.Bytes(), sampleRate)
}

func makeWAVCustom(format uint16, bits int, payload []byte, sampleRate int) []byte {
	buf := &bytes.Buffer{}
	riffSize := 4 + (8 + 16) + (8 + len(payload))

	buf.WriteString("RIFF")
	writeU32(buf, uint32(riffSize))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, format)
	writeU16(buf, 1)
	writeU32(buf, uint32(sampleRate))
	byteRate := sampleRate * (bits / 8)
	writeU32(buf, uint32(byteRate))
	blockAlign := bits / 8
	writeU16(buf, uint16(blockAlign))
	writeU16(buf, uint16(bits))

	buf.WriteString("data")
	writeU32(buf, uint32(len(payload)))
	buf.Write(payload)

	return buf.Bytes()
}
