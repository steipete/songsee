package audio

import (
	"bytes"
	"testing"
)

func TestDecodeWAVExtensiblePCM(t *testing.T) {
	payload := []byte{0, 0, 0, 0}
	buf := &bytes.Buffer{}
	riffSize := 4 + (8 + 40) + (8 + len(payload))

	buf.WriteString("RIFF")
	writeU32(buf, uint32(riffSize))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 40)
	writeU16(buf, 0xFFFE)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*2)
	writeU16(buf, 2)
	writeU16(buf, 16)
	writeU16(buf, 22)
	writeU16(buf, 16)
	writeU32(buf, 1)
	buf.Write([]byte{
		0x01, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x10, 0x00,
		0x80, 0x00,
		0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71,
	})

	buf.WriteString("data")
	writeU32(buf, uint32(len(payload)))
	buf.Write(payload)

	pcm, err := DecodeBytes(buf.Bytes(), Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}
