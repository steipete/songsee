package main

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestRunMP3E2E(t *testing.T) {
	input := testdataPath(t, "sine.mp3")
	outDir := t.TempDir()
	outPath := filepath.Join(outDir, "spectro.jpg")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{
		"--width", "320",
		"--height", "180",
		"--start", "0.2",
		"--duration", "0.5",
		"--style", "magma",
		"--output", outPath,
		input,
	}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if stdout.String() == "" {
		t.Fatalf("expected stdout output")
	}
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("missing output: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("empty output")
	}

	file, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer func() { _ = file.Close() }()
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("decode image: %v", err)
	}
	if img.Bounds().Dx() != 320 || img.Bounds().Dy() != 180 {
		t.Fatalf("size mismatch")
	}
	if flatImage(img) {
		t.Fatalf("image appears flat")
	}
}

func TestRunFromStdinPNG(t *testing.T) {
	outDir := t.TempDir()
	outPath := filepath.Join(outDir, "spectro.png")

	wav := makeWAV([]int16{0, 2000, -2000, 0, 1000, -1000}, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{
		"--format", "png",
		"--output", outPath,
		"-",
	}, bytes.NewReader(wav), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	file, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer func() { _ = file.Close() }()
	if _, err := png.Decode(file); err != nil {
		t.Fatalf("decode png: %v", err)
	}
}

func TestRunVersion(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--version"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d", exit)
	}
	if stdout.String() == "" {
		t.Fatalf("expected version output")
	}
}

func TestRunHelp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--help"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d", exit)
	}
	if stdout.String() == "" && stderr.String() == "" {
		t.Fatalf("expected help output")
	}
}

func TestRunInvalidWindow(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--window", "1000", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr usage")
	}
}

func TestRunUnknownFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--nope"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr output")
	}
}

func TestRunBadFormat(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "gif", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunBadFreqRange(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--min-freq", "100", "--max-freq", "50", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunUnknownStyle(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--style", "nope", "-"}, bytes.NewReader(makeWAV([]int16{0, 1}, 44100, 1)), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunBadSize(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--width", "0", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunBadWindowZero(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--window", "0", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunBadHopZero(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--hop", "0", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunNegativeStart(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--start=-1", "-"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
}

func TestRunMissingFile(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"nope.wav"}, bytes.NewReader(nil), stdout, stderr)
	if exit != 1 {
		t.Fatalf("expected error exit, got %d", exit)
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr output")
	}
}

func TestRunMissingInput(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{}, bytes.NewReader(nil), stdout, stderr)
	if exit != 2 {
		t.Fatalf("expected usage exit, got %d", exit)
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr output")
	}
}

func TestRunNoSamplesDecoded(t *testing.T) {
	wav := makeWAV([]int16{}, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"-"}, bytes.NewReader(wav), stdout, stderr)
	if exit != 1 {
		t.Fatalf("expected error exit, got %d", exit)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("no samples")) {
		t.Fatalf("expected no samples error")
	}
}

func TestRunSliceError(t *testing.T) {
	wav := makeWAV([]int16{0, 1, -1, 0}, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--start", "2", "--duration", "1", "-"}, bytes.NewReader(wav), stdout, stderr)
	if exit != 1 {
		t.Fatalf("expected error exit, got %d", exit)
	}
}

func TestRunSliceVerbose(t *testing.T) {
	samples := make([]int16, 44100)
	wav := makeWAV(samples, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--verbose", "--start", "0", "--duration", "0.2", "--output", "-", "-"}, bytes.NewReader(wav), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if !bytes.Contains(stderr.Bytes(), []byte("slice:")) {
		t.Fatalf("expected slice output")
	}
}

func TestRunStyleAffectsOutput(t *testing.T) {
	wav := makeWAV(genSineMixSamples(44100), 44100, 1)
	outClassic := runToBytes(t, wav, "classic")
	outMagma := runToBytes(t, wav, "magma")
	if bytes.Equal(outClassic, outMagma) {
		t.Fatalf("expected different output for different styles")
	}
}

func TestRunOutputStdout(t *testing.T) {
	wav := makeWAV([]int16{0, 1000, -1000, 0, 500, -500}, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{
		"--format", "png",
		"--output", "-",
		"-",
	}, bytes.NewReader(wav), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected image bytes on stdout")
	}
	if _, err := png.Decode(bytes.NewReader(stdout.Bytes())); err != nil {
		t.Fatalf("decode stdout png: %v", err)
	}
}

func TestRunOutputAuto(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "png", input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	outPath := filepath.Join(dir, "input.png")
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunOutputExtOverride(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out.png")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunOutputAppendDefault(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(output + ".jpg"); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunOutputJpgExt(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out.jpg")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunWriteImageError(t *testing.T) {
	wav := makeWAV([]int16{0, 1000, -1000, 0}, 44100, 1)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--output", "/nope/dir/out.jpg", "-"}, bytes.NewReader(wav), stdout, stderr)
	if exit != 1 {
		t.Fatalf("expected error exit, got %d", exit)
	}
}

func TestRunFormatFlagKeepsOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "customout")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "png", "--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunQuiet(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out.jpg")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--quiet", "--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("expected quiet stdout")
	}
}

func TestDie(t *testing.T) {
	stderr := &bytes.Buffer{}
	if code := die(stderr, errSentinel{}); code != 1 {
		t.Fatalf("expected code 1")
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr output")
	}
}

type errSentinel struct{}

func (errSentinel) Error() string { return "boom" }

func TestRunInputDashDefaultOutput(t *testing.T) {
	tmp := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "png", "-"}, bytes.NewReader(makeWAV([]int16{0, 1, -1}, 44100, 1)), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(tmp, "songsee.png")); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestRunFormatJPEG(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "jpeg", "--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("missing output: %v", err)
	}
}

func TestWriteImageUnknownFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	err := writeImage("-", "gif", image.NewRGBA(image.Rect(0, 0, 1, 1)), buf)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunVerbose(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.wav")
	if err := os.WriteFile(input, makeWAV([]int16{0, 1, -1, 0}, 44100, 1), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	output := filepath.Join(dir, "out.jpg")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--verbose", "--output", output, input}, bytes.NewReader(nil), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	if !bytes.Contains(stderr.Bytes(), []byte("decoded:")) {
		t.Fatalf("expected verbose output")
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

func flatImage(img image.Image) bool {
	bounds := img.Bounds()
	minLum := uint32(0xFFFFFFFF)
	maxLum := uint32(0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := (r + g + b) / 3
			if lum < minLum {
				minLum = lum
			}
			if lum > maxLum {
				maxLum = lum
			}
		}
	}
	return maxLum-minLum < 1000
}

func makeWAV(samples []int16, sampleRate int, channels int) []byte {
	if channels < 1 {
		channels = 1
	}
	dataLen := len(samples) * 2
	riffSize := 4 + (8 + 16) + (8 + dataLen)

	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeU32(buf, uint32(riffSize))
	buf.WriteString("WAVE")

	buf.WriteString("fmt ")
	writeU32(buf, 16)
	writeU16(buf, 1)
	writeU16(buf, uint16(channels))
	writeU32(buf, uint32(sampleRate))
	byteRate := sampleRate * channels * 2
	writeU32(buf, uint32(byteRate))
	blockAlign := channels * 2
	writeU16(buf, uint16(blockAlign))
	writeU16(buf, 16)

	buf.WriteString("data")
	writeU32(buf, uint32(dataLen))
	for _, s := range samples {
		writeU16(buf, uint16(s))
	}

	return buf.Bytes()
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

func runToBytes(t *testing.T, wav []byte, style string) []byte {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exit := run([]string{"--format", "png", "--style", style, "--output", "-", "-"}, bytes.NewReader(wav), stdout, stderr)
	if exit != 0 {
		t.Fatalf("exit %d stderr=%s", exit, stderr.String())
	}
	return stdout.Bytes()
}

func genSineMixSamples(n int) []int16 {
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n)
		v := 0.5*math.Sin(2*math.Pi*440*t) + 0.4*math.Sin(2*math.Pi*880*t)
		out[i] = int16(v * 15000)
	}
	return out
}
