// Package main provides the songsee CLI entrypoint.
package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/steipete/songsee/internal/audio"
	"github.com/steipete/songsee/internal/dsp"
	"github.com/steipete/songsee/internal/render"
)

var version = "dev"

type cli struct {
	Input      string           `arg:"" help:"file path or '-' for stdin"`
	Output     string           `short:"o" help:"output image path"`
	Format     string           `help:"output format: jpg or png" default:"jpg"`
	Width      int              `help:"output width in pixels" default:"1920"`
	Height     int              `help:"output height in pixels" default:"1080"`
	WindowSize int              `name:"window" help:"FFT window size in samples" default:"2048"`
	HopSize    int              `name:"hop" help:"hop size in samples" default:"512"`
	MinFreq    float64          `name:"min-freq" help:"minimum frequency in Hz"`
	MaxFreq    float64          `name:"max-freq" help:"maximum frequency in Hz (0 = Nyquist)"`
	StartSec   float64          `name:"start" help:"start time in seconds"`
	Duration   float64          `name:"duration" help:"duration in seconds (0 = full)"`
	SampleRate int              `name:"sample-rate" help:"ffmpeg output sample rate" default:"44100"`
	Style      string           `help:"palette style: classic, magma, inferno, viridis, gray" default:"classic"`
	FFmpegPath string           `name:"ffmpeg" help:"path to ffmpeg binary"`
	Quiet      bool             `short:"q" help:"suppress stdout output"`
	Verbose    bool             `short:"v" help:"verbose stderr output"`
	Version    kong.VersionFlag `name:"version" help:"print version"`
}

type exitPanic struct {
	code int
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	formatSet := hasFlag(args, "--format")
	cfg := cli{}
	exitCode := -1

	parser, err := kong.New(&cfg,
		kong.Name("songsee"),
		kong.Description("generate a classic spectrogram image"),
		kong.Vars{"version": version},
		kong.Writers(stdout, stderr),
		kong.Exit(func(code int) { panic(exitPanic{code: code}) }),
	)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "songsee:", err)
		return 1
	}

	var ctx *kong.Context
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				if exit, ok := recovered.(exitPanic); ok {
					exitCode = exit.code
					return
				}
				panic(recovered)
			}
		}()
		ctx, err = parser.Parse(args)
	}()
	if exitCode >= 0 {
		return exitCode
	}
	if err != nil {
		if parseErr, ok := err.(*kong.ParseError); ok {
			_, _ = fmt.Fprintln(stderr, "songsee:", parseErr)
			if parseErr.Context != nil {
				parseErr.Context.Stdout = stderr
				_ = parseErr.Context.PrintUsage(false)
			}
			return 2
		}
		_, _ = fmt.Fprintln(stderr, "songsee:", err)
		return 1
	}

	input := cfg.Input
	if input == "" {
		if ctx != nil {
			ctx.Stdout = stderr
			_ = ctx.PrintUsage(false)
		}
		return 2
	}

	if cfg.MaxFreq > 0 && cfg.MaxFreq <= cfg.MinFreq {
		return dieUsage(stderr, ctx, "--max-freq must be > --min-freq")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return dieUsage(stderr, ctx, "--width and --height must be > 0")
	}
	if cfg.WindowSize <= 0 || cfg.HopSize <= 0 {
		return dieUsage(stderr, ctx, "--window and --hop must be > 0")
	}
	if !isPowerOfTwo(cfg.WindowSize) {
		return dieUsage(stderr, ctx, "--window must be a power of two")
	}
	if cfg.StartSec < 0 || cfg.Duration < 0 {
		return dieUsage(stderr, ctx, "--start and --duration must be >= 0")
	}

	format := strings.ToLower(cfg.Format)
	if format != "jpg" && format != "jpeg" && format != "png" {
		return dieUsage(stderr, ctx, "--format must be jpg or png")
	}
	if format == "jpeg" {
		format = "jpg"
	}

	output := cfg.Output
	if output == "" {
		if input == "-" {
			output = "songsee." + format
		} else {
			ext := strings.ToLower(filepath.Ext(input))
			base := strings.TrimSuffix(filepath.Base(input), ext)
			output = filepath.Join(filepath.Dir(input), base+"."+format)
		}
	} else {
		ext := strings.ToLower(filepath.Ext(output))
		switch ext {
		case ".png":
			format = "png"
		case ".jpg", ".jpeg":
			format = "jpg"
		default:
			if !formatSet {
				output = output + "." + format
			}
		}
	}

	if cfg.Verbose {
		_, _ = fmt.Fprintf(stderr, "input: %s\n", input)
		_, _ = fmt.Fprintf(stderr, "output: %s (%s)\n", output, format)
	}

	opts := audio.Options{SampleRate: cfg.SampleRate, FFmpegPath: cfg.FFmpegPath}
	var pcm audio.Audio
	if input == "-" {
		pcm, err = audio.DecodeReader(stdin, opts)
	} else {
		pcm, err = audio.DecodeFile(input, opts)
	}
	if err != nil {
		return die(stderr, err)
	}
	if len(pcm.Samples) == 0 {
		return die(stderr, errors.New("no samples decoded"))
	}
	if cfg.Verbose {
		_, _ = fmt.Fprintf(stderr, "decoded: %d samples @ %d Hz\n", len(pcm.Samples), pcm.SampleRate)
	}
	if cfg.StartSec > 0 || cfg.Duration > 0 {
		pcm, err = audio.Slice(pcm, cfg.StartSec, cfg.Duration)
		if err != nil {
			return die(stderr, err)
		}
		if cfg.Verbose {
			_, _ = fmt.Fprintf(stderr, "slice: %0.2fs + %0.2fs => %d samples\n", cfg.StartSec, cfg.Duration, len(pcm.Samples))
		}
	}

	spec := dsp.ComputeSpectrogram(pcm.Samples, pcm.SampleRate, cfg.WindowSize, cfg.HopSize)
	style := strings.ToLower(strings.TrimSpace(cfg.Style))
	palette, err := render.PaletteByName(style)
	if err != nil {
		return dieUsage(stderr, ctx, "unknown style")
	}

	img, err := render.Spectrogram(spec, render.Options{
		Width:   cfg.Width,
		Height:  cfg.Height,
		MinFreq: cfg.MinFreq,
		MaxFreq: cfg.MaxFreq,
		Palette: palette,
	})
	if err != nil {
		return die(stderr, err)
	}

	if err := writeImage(output, format, img, stdout); err != nil {
		return die(stderr, err)
	}

	if output != "-" && !cfg.Quiet {
		_, _ = fmt.Fprintln(stdout, output)
	}
	return 0
}

func writeImage(path, format string, img image.Image, stdout io.Writer) error {
	var out io.Writer
	if path == "-" {
		out = stdout
	} else {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()
		out = file
	}

	switch format {
	case "png":
		return png.Encode(out, img)
	case "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	default:
		return fmt.Errorf("unknown format %s", format)
	}
}

func die(stderr io.Writer, err error) int {
	_, _ = fmt.Fprintln(stderr, "songsee:", err)
	return 1
}

func dieUsage(stderr io.Writer, ctx *kong.Context, msg string) int {
	_, _ = fmt.Fprintln(stderr, "songsee:", msg)
	if ctx != nil {
		ctx.Stdout = stderr
		_ = ctx.PrintUsage(false)
	}
	return 2
}

func isPowerOfTwo(v int) bool {
	return v > 0 && (v&(v-1)) == 0
}

func hasFlag(args []string, name string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == name || strings.HasPrefix(arg, name+"=") {
			return true
		}
	}
	return false
}
