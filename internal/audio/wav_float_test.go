package audio

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDecodeWAVFloat32(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4+(8+16)+(8+8))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 3)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*4)
	writeU16(buf, 4)
	writeU16(buf, 32)

	buf.WriteString("data")
	writeU32(buf, 8)
	_ = binary.Write(buf, binary.LittleEndian, float32(0.5))
	_ = binary.Write(buf, binary.LittleEndian, float32(-0.25))

	pcm, err := DecodeBytes(buf.Bytes(), Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) != 2 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeWAVFloat64(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4+(8+16)+(8+16))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 3)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*8)
	writeU16(buf, 8)
	writeU16(buf, 64)

	buf.WriteString("data")
	writeU32(buf, 16)
	_ = binary.Write(buf, binary.LittleEndian, float64(0.5))
	_ = binary.Write(buf, binary.LittleEndian, float64(-0.25))

	pcm, err := DecodeBytes(buf.Bytes(), Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) != 2 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func writeU16(buf *bytes.Buffer, v uint16) {
	buf.WriteByte(byte(v))
	buf.WriteByte(byte(v >> 8))
}

func writeU32(buf *bytes.Buffer, v uint32) {
	buf.WriteByte(byte(v))
	buf.WriteByte(byte(v >> 8))
	buf.WriteByte(byte(v >> 16))
	buf.WriteByte(byte(v >> 24))
}
