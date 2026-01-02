(() => {
  const panels = document.querySelectorAll('.spectral-panel');
  if (!panels.length) return;

  const update = (panel, event) => {
    const rect = panel.getBoundingClientRect();
    const x = Math.min(Math.max((event.clientX - rect.left) / rect.width, 0), 1);
    const y = Math.min(Math.max((event.clientY - rect.top) / rect.height, 0), 1);
    panel.style.setProperty('--mx', `${(x * 100).toFixed(1)}%`);
    panel.style.setProperty('--my', `${(y * 100).toFixed(1)}%`);
  };

  panels.forEach((panel) => {
    panel.addEventListener('pointermove', (event) => update(panel, event));
    panel.addEventListener('pointerleave', () => {
      panel.style.removeProperty('--mx');
      panel.style.removeProperty('--my');
    });
  });
})();
