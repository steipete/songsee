---
layout: default
title: Spec
description: songsee spectral pipeline, defaults, and rendering details.
---

<section class="section">
  <div class="kicker">Spec</div>
  <h1 class="section-title">songsee spectral pipeline</h1>
  <p class="section-sub">
    This page captures the core algorithm and defaults used by songsee for repeatable, high quality
    spectrogram images.
  </p>
</section>

<section class="section">
  <h2 class="section-title">Decode</h2>
  <div class="card">
    <p>
      WAV and MP3 decode natively. Any other format falls back to ffmpeg. Input can be a file path or
      stdin ("-"). Default sample rate for ffmpeg output is 44100 Hz.
    </p>
  </div>
</section>

<section class="section">
  <h2 class="section-title">Spectrogram</h2>
  <div class="card">
    <p>
      Windowed frames use a Hann window. FFT runs on each frame and the magnitude is converted to
      decibels using 20 * log10(mag + 1e-9). The default window size is 2048 samples with a hop size
      of 512 samples.
    </p>
    <p>
      Frames are computed as 1 + (len(samples) - window + hop - 1) / hop, and bins are window/2 + 1.
      Bin spacing is sampleRate / windowSize.
    </p>
  </div>
</section>

<section class="section">
  <h2 class="section-title">Rendering</h2>
  <div class="card">
    <p>
      Each output pixel maps to a time frame and frequency bin. Values are normalized by the global
      min/max in the computed spectrogram unless clamp values are provided. Frequency range can be
      restricted via min/max frequency in Hz.
    </p>
    <p>
      Output size defaults to 1920x1080. JPEG quality is 95. PNG output is available via --format.
    </p>
  </div>
</section>

<section class="section">
  <h2 class="section-title">Palettes</h2>
  <div class="card">
    <p>
      Palettes map normalized values to RGBA colors. Available names: classic, magma, inferno,
      viridis, gray, clawd. The clawd palette shifts from deep ocean to coral highlights.
    </p>
  </div>
</section>

<section class="section">
  <h2 class="section-title">CLI defaults</h2>
  <div class="code-block">
    --format jpg
    --width 1920
    --height 1080
    --window 2048
    --hop 512
    --sample-rate 44100
    --style classic
  </div>
</section>
