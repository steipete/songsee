package audio

import (
	"bytes"
	"os"
	"testing"
)

func TestResolveFFmpeg(t *testing.T) {
	path, err := resolveFFmpeg("")
	if err != nil {
		t.Fatalf("resolveFFmpeg: %v", err)
	}
	if path == "" {
		t.Fatalf("empty ffmpeg path")
	}
}

func TestResolveFFmpegExplicit(t *testing.T) {
	found, err := resolveFFmpeg("")
	if err != nil {
		t.Fatalf("resolveFFmpeg: %v", err)
	}
	path, err := resolveFFmpeg(found)
	if err != nil {
		t.Fatalf("resolveFFmpeg explicit: %v", err)
	}
	if path == "" {
		t.Fatalf("empty ffmpeg path")
	}
}

func TestResolveFFmpegMissing(t *testing.T) {
	t.Setenv("PATH", "")
	if _, err := resolveFFmpeg(""); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWithFFmpegFile(t *testing.T) {
	path := testdataPath(t, "sine.mp3")
	pcm, err := DecodeWithFFmpeg(path, nil, 22050, "")
	if err != nil {
		t.Fatalf("DecodeWithFFmpeg: %v", err)
	}
	if pcm.SampleRate != 22050 {
		t.Fatalf("sample rate = %d", pcm.SampleRate)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeWithFFmpegStdin(t *testing.T) {
	path := testdataPath(t, "sine.mp3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	pcm, err := DecodeWithFFmpeg("", bytes.NewReader(data), 44100, "")
	if err != nil {
		t.Fatalf("DecodeWithFFmpeg stdin: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeWithFFmpegBadPath(t *testing.T) {
	_, err := DecodeWithFFmpeg("missing.mp3", nil, 0, "/no/such/ffmpeg")
	if err == nil {
		t.Fatalf("expected error")
	}
}
