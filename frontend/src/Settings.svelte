<script lang="ts">
  import { getTheme, setTheme } from "./lib/theme";

  let darkMode = $state(getTheme() === "dark");

  function toggleDarkMode() {
    darkMode = !darkMode;
    setTheme(darkMode ? "dark" : "light");
  }
</script>

<h1>Settings</h1>

<div class="setting-row">
  <div class="setting-text">
    <span class="setting-label">Mode nuit</span>
    <p class="setting-desc">Bascule l'interface sur un thème sombre.</p>
  </div>
  <button
    class="switch"
    class:on={darkMode}
    role="switch"
    aria-checked={darkMode}
    aria-label="Mode nuit"
    onclick={toggleDarkMode}
  >
    <span class="switch-thumb"></span>
  </button>
</div>

<style>
  .setting-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.5rem;
    max-width: 420px;
    margin-top: 1.75rem;
    padding: 1.25rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
  }

  .setting-label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--color-text);
  }

  .setting-desc {
    margin: 0.25rem 0 0;
    font-size: 0.75rem;
    color: var(--color-text-secondary);
  }

  .switch {
    all: unset;
    box-sizing: border-box;
    flex-shrink: 0;
    width: 40px;
    height: 24px;
    padding: 2px;
    border-radius: 999px;
    background: var(--color-border-strong);
    cursor: pointer;
    transition: background-color 0.15s ease;
  }

  .switch.on {
    background: var(--color-accent);
  }

  .switch-thumb {
    display: block;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--color-surface);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.25);
    transition: transform 0.15s ease;
  }

  .switch.on .switch-thumb {
    transform: translateX(16px);
  }
</style>
