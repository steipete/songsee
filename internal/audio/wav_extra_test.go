package audio

import (
	"bytes"
	"testing"
)

func TestDecodeWAVWithExtraChunk(t *testing.T) {
	payload := []byte{0, 0, 0, 0}
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, uint32(4+(8+16)+(8+5)+(8+len(payload))))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 1)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*2)
	writeU16(buf, 2)
	writeU16(buf, 16)

	buf.WriteString("JUNK")
	writeU32(buf, 5)
	buf.Write([]byte{1, 2, 3, 4, 5})
	buf.WriteByte(0)

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

func TestDecodeWAVMissingData(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4+(8+16))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 1)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*2)
	writeU16(buf, 2)
	writeU16(buf, 16)

	if _, err := DecodeBytes(buf.Bytes(), Options{}); err == nil {
		t.Fatalf("expected error for missing data")
	}
}

func TestDecodeWAVHeaderOnly(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4)
	buf.WriteString("WAVE")
	if _, err := decodeWAV(bytes.NewReader(buf.Bytes())); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWAVInvalidChannels(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4+(8+16)+(8+2))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 1)
	writeU16(buf, 0)
	writeU32(buf, 44100)
	writeU32(buf, 44100*2)
	writeU16(buf, 2)
	writeU16(buf, 16)

	buf.WriteString("data")
	writeU32(buf, 2)
	buf.Write([]byte{0, 0})

	if _, err := DecodeBytes(buf.Bytes(), Options{}); err == nil {
		t.Fatalf("expected error for channels")
	}
}

func TestDecodeWAVFloatUnsupportedBits(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, 4+(8+16)+(8+3))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 3)
	writeU16(buf, 1)
	writeU32(buf, 44100)
	writeU32(buf, 44100*3)
	writeU16(buf, 3)
	writeU16(buf, 24)

	buf.WriteString("data")
	writeU32(buf, 3)
	buf.Write([]byte{0, 0, 0})

	if _, err := DecodeBytes(buf.Bytes(), Options{}); err == nil {
		t.Fatalf("expected error for float bit depth")
	}
}
