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

// Reset the HTML element style by removing the attribute `style` after
// theme updated to avoid the flash between light and dark colors.
document.documentElement.style.display = "none";

document.addEventListener("DOMContentLoaded", function () {
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
  document.documentElement.removeAttribute("style")
});

updateTheme(false);
