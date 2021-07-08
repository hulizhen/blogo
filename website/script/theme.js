const themeKey = "blogo_theme";

const updateTheme = (toggle) => {
  let dark = localStorage[themeKey] === "dark";
  if (toggle) {
    dark = !dark;
  }
  localStorage[themeKey] = dark ? "dark" : "light";

  document.documentElement.dataset.theme = dark ? "dark" : "light";
};

window.onload = () => {
  const themeToggler = document.getElementById("theme-toggler");

  themeToggler.updateStatus = function () {
    this.innerHTML = localStorage[themeKey] === "dark" ? "🌙" : "🌞";
  };
  themeToggler.addEventListener("click", () => {
    updateTheme(true);
    themeToggler.updateStatus();
  });

  themeToggler.updateStatus();
};

updateTheme(false);
