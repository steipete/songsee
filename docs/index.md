---
layout: default
title: Home
description: Generate modern spectrogram images from audio files with a fast, scriptable CLI.
body_class: home
---

<section class="hero">
  <div class="hero-copy">
    <div class="kicker reveal delay-1">Spectral imaging CLI</div>
    <h1 class="hero-title reveal delay-2">See sound as living color.</h1>
    <p class="hero-sub reveal delay-3">
      songsee turns audio into precise, high-resolution spectrograms and feature panels. Fast decode
      paths for WAV and MP3, ffmpeg fallback for everything else, and palette styles that make science
      look cinematic.
    </p>
    <div class="hero-actions reveal delay-4">
      <a class="btn primary" href="#install">Install</a>
      <a class="btn" href="https://github.com/steipete/songsee">GitHub</a>
    </div>
    <div class="hero-meta reveal delay-4">Hann window. Log magnitude. 2048 / 512 defaults.</div>
  </div>
  <div class="hero-visual">
    <div class="spectral-panel" role="img" aria-label="Animated spectrogram preview">
      <div class="spectral-caption">Spectrogram preview</div>
    </div>
  </div>
</section>

<section class="section">
  <div class="kicker">Why songsee</div>
  <h2 class="section-title">A focused pipeline for modern spectrograms.</h2>
  <p class="section-sub">
    Decode audio into mono samples, window it with Hann, run FFT, and render log-magnitude frames into
    a crisp image. The CLI stays small, reliable, and scriptable.
  </p>

  <div class="feature-grid">
    <div class="card">
      <h3>Precise controls</h3>
      <p>Window, hop, min/max frequency, output dimensions, and time slicing for exact framing.</p>
    </div>
    <div class="card">
      <h3>Fast decode paths</h3>
      <p>Native WAV/MP3 decoding with ffmpeg fallback for everything else.</p>
    </div>
    <div class="card">
      <h3>Palette styles</h3>
      <p>classic, magma, inferno, viridis, and gray for a bold spectral aesthetic.</p>
    </div>
    <div class="card">
      <h3>Feature panels</h3>
      <p>mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux — rendered as single or grid views.</p>
    </div>
    <div class="card">
      <h3>Clean output</h3>
      <p>JPEG or PNG output, default quality 95, and stable results for batch workflows.</p>
    </div>
  </div>
</section>

<section class="section" id="install">
  <div class="kicker">Install</div>
  <h2 class="section-title">One command. Instant spectrograms.</h2>
  <div class="code-block">
    go install github.com/steipete/songsee/cmd/songsee@latest
  </div>
  <div class="domain-note">
    songsee.ai, songsee.app, and songsee.dev all redirect to songsee.sh.
  </div>
</section>

<section class="section" id="usage">
  <div class="kicker">Usage</div>
  <h2 class="section-title">CLI ready for pipes, batches, and automation.</h2>
  <div class="code-block">
    songsee track.mp3
    songsee track.wav --style magma --width 2048 --height 1024 -o spectro.png
    cat track.mp3 | songsee - --style gray --format png
    songsee track.mp3 --start 12.5 --duration 8 --output slice.jpg
    songsee track.mp3 --viz spectrogram,mel,chroma --width 2048 --height 1024
  </div>
</section>

<section class="section">
  <div class="kicker">Palettes</div>
  <h2 class="section-title">Color maps with character.</h2>
  <p class="section-sub">Pick a palette by name for instant visual tone shifts.</p>
  <div class="palette-row">
    <div class="palette classic" title="classic"></div>
    <div class="palette magma" title="magma"></div>
    <div class="palette inferno" title="inferno"></div>
    <div class="palette viridis" title="viridis"></div>
    <div class="palette gray" title="gray"></div>
  </div>
</section>

<section class="section">
  <div class="kicker">Specs</div>
  <h2 class="section-title">Detailed pipeline notes.</h2>
  <p class="section-sub">
    Windowing, bin mapping, normalization, and rendering details live in the spec.
  </p>
  <div class="hero-actions">
    <a class="btn" href="{{ '/spec/' | relative_url }}">Read spec</a>
  </div>
</section>
