const themeKey = "blogo_theme";
const getDarkMode = () => localStorage[themeKey] === "dark";
const setDarkMode = (dark) => (localStorage[themeKey] = dark ? "dark" : "light");

const updateTheme = (toggle) => {
  let dark = getDarkMode();
  if (toggle) {
    dark = !dark;
    setDarkMode(dark);
  }
  document.documentElement.dataset.theme = dark ? "dark" : "light";
};

window.onload = () => {
  const themeToggler = document.getElementById("theme-toggler");
  themeToggler.updateStatus = function () {
    const sunIcon = "bi bi-sun-fill";
    const moonIcon = "bi bi-moon-fill";
    this.className = getDarkMode() ? moonIcon : sunIcon;
  };
  themeToggler.addEventListener("click", () => {
    updateTheme(true);
    themeToggler.updateStatus();
  });
  themeToggler.updateStatus();
};

updateTheme(false);
