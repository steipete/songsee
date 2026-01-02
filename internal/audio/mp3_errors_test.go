package audio

import (
	"bytes"
	"testing"
)

func TestDecodeMP3Error(t *testing.T) {
	if _, err := decodeMP3(bytes.NewReader([]byte("not mp3"))); err == nil {
		t.Fatalf("expected error")
	}
}

func TestBinaryReadError(t *testing.T) {
	var s int16
	if err := binaryRead(bytes.NewReader([]byte{0x01}), &s); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeMP3IfCorrupt(t *testing.T) {
	data := []byte{'I', 'D', '3', 0x03, 0x00}
	_, ok, err := DecodeMP3If(bytes.NewReader(data))
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}
