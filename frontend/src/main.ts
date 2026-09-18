import { mount } from "svelte";
import App from "./App.svelte";
import { applyTheme, getTheme } from "./lib/theme";
import "./style.css";

applyTheme(getTheme()); // avant le mount, pour éviter un flash en thème clair

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;
