package audio

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeWAVBytes(t *testing.T) {
	samples := make([]int16, 1000)
	for i := range samples {
		samples[i] = int16(2000 * math.Sin(2*math.Pi*float64(i)/50))
	}
	data := makeWAV(samples, 44100, 1)
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if pcm.SampleRate != 44100 {
		t.Fatalf("sample rate = %d", pcm.SampleRate)
	}
	if len(pcm.Samples) != len(samples) {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeWAVIfNotWAV(t *testing.T) {
	_, ok, err := DecodeWAVIf(bytesReader([]byte("NOTWAVE12345")))
	if err != nil {
		t.Fatalf("DecodeWAVIf error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}

func TestDecodeWAVIfValid(t *testing.T) {
	data := makeWAV([]int16{0, 1000, -1000}, 44100, 1)
	pcm, ok, err := DecodeWAVIf(bytesReader(data))
	if err != nil {
		t.Fatalf("DecodeWAVIf error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if len(pcm.Samples) != 3 {
		t.Fatalf("samples = %d", len(pcm.Samples))
	}
}

func TestDecodeMP3File(t *testing.T) {
	path := testdataPath(t, "sine.mp3")
	pcm, err := DecodeFile(path, Options{})
	if err != nil {
		t.Fatalf("DecodeFile: %v", err)
	}
	if pcm.SampleRate == 0 || len(pcm.Samples) == 0 {
		t.Fatalf("invalid decode result")
	}
}

func TestDecodeMP3IfValid(t *testing.T) {
	path := testdataPath(t, "sine.mp3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	pcm, ok, err := DecodeMP3If(bytesReader(data))
	if err != nil {
		t.Fatalf("DecodeMP3If: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeFileUnknownExt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audio.bin")
	if err := os.WriteFile(path, makeWAV([]int16{0, 1, -1}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	pcm, err := DecodeFile(path, Options{})
	if err != nil {
		t.Fatalf("DecodeFile: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeFileWAV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audio.wav")
	if err := os.WriteFile(path, makeWAV([]int16{0, 1, -1}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	pcm, err := DecodeFile(path, Options{})
	if err != nil {
		t.Fatalf("DecodeFile: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeFileFFmpegFallbackError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audio.bin")
	if err := os.WriteFile(path, []byte("garbagegarbagegarbage"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := DecodeFile(path, Options{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeBytesFFmpegFallbackError(t *testing.T) {
	_, err := DecodeBytes([]byte("not audio"), Options{})
	if err == nil {
		t.Fatalf("expected error for garbage data")
	}
}

func TestDecodeMP3Bytes(t *testing.T) {
	path := testdataPath(t, "sine.mp3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatalf("DecodeBytes: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeReader(t *testing.T) {
	samples := []int16{0, 1000, -1000, 0}
	data := makeWAV(samples, 48000, 1)
	pcm, err := DecodeReader(bytesReader(data), Options{})
	if err != nil {
		t.Fatalf("DecodeReader: %v", err)
	}
	if pcm.SampleRate != 48000 {
		t.Fatalf("sample rate = %d", pcm.SampleRate)
	}
}

func TestDecodeReaderError(t *testing.T) {
	_, err := DecodeReader(errReader{}, Options{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSlice(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}
	out, err := Slice(a, 0.2, 0.5)
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(out.Samples) != 5 {
		t.Fatalf("slice samples = %d", len(out.Samples))
	}
}

func TestSliceErrors(t *testing.T) {
	_, err := Slice(Audio{SampleRate: 10, Samples: []float64{1}}, -1, 1)
	if err == nil {
		t.Fatalf("expected error for negative start")
	}
	_, err = Slice(Audio{SampleRate: 0, Samples: []float64{1}}, 0, 1)
	if err == nil {
		t.Fatalf("expected error for sample rate")
	}
	_, err = Slice(Audio{SampleRate: 10, Samples: []float64{}}, 0, 1)
	if err == nil {
		t.Fatalf("expected error for empty samples")
	}
	_, err = Slice(Audio{SampleRate: 10, Samples: []float64{1}}, 2, 0)
	if err == nil {
		t.Fatalf("expected error for start")
	}
}

func TestSliceFullDuration(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2, 3}}
	out, err := Slice(a, 0, 0)
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}
	if len(out.Samples) != len(a.Samples) {
		t.Fatalf("expected full slice")
	}
}

func TestSliceDurationTooShort(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2}}
	if _, err := Slice(a, 0, 0.01); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeMP3IfNotMP3(t *testing.T) {
	_, ok, err := DecodeMP3If(bytesReader([]byte("NOTMP3DATA")))
	if err != nil {
		t.Fatalf("DecodeMP3If error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}

func TestDecodeWAVUnsupported(t *testing.T) {
	_, err := decodeWAV(bytesReader([]byte("NOTWAVE12345")))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func testdataPath(t *testing.T, name string) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	root := filepath.Dir(filepath.Dir(wd))
	path := filepath.Join(root, "testdata", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("missing testdata: %v", err)
	}
	return path
}

func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, os.ErrInvalid }
