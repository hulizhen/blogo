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
  const updateStatus = () => {
    const sunIcon = "bi bi-sun-fill";
    const moonIcon = "bi bi-moon-fill";
    const dark = getDarkMode();
    navbar = document.getElementById("navbar").classList;
    navbar.remove(dark ? "navbar-light" : "navbar-dark");
    navbar.add(dark ? "navbar-dark" : "navbar-light");
    themeToggler.className = dark ? moonIcon : sunIcon;
  };
  themeToggler.addEventListener("click", () => {
    updateTheme(true);
    updateStatus();
  });
  updateStatus();
};

updateTheme(false);
