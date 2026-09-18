// Mode nuit : préférence stockée localement et appliquée via l'attribut
// data-theme sur <html>, que style.css utilise pour retheme toute l'app.
const STORAGE_KEY = "patio:theme";

export type Theme = "light" | "dark";

export function getTheme(): Theme {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "light" || stored === "dark") return stored;
  } catch {
    // localStorage indisponible (navigation privée, etc.) -> "light" par défaut
  }
  return "light";
}

export function applyTheme(theme: Theme) {
  document.documentElement.setAttribute("data-theme", theme);
}

export function setTheme(theme: Theme) {
  applyTheme(theme);
  try {
    localStorage.setItem(STORAGE_KEY, theme);
  } catch {
    // tant pis, le choix ne survivra pas au redémarrage
  }
}
