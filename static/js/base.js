window.onload = function () {
  const className = "selected";

  let node = document.querySelector(`.nav ul a.${className}`);
  if (node) {
    node.classList.remove(className);
  }

  node = document.querySelector(`.nav ul a[href='${window.location.pathname}']`);
  if (node) {
    node.classList.add(className);
  }
};
