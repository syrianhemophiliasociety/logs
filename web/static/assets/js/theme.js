"use strict";

const themes = {
  black: {
    primary: "#000000",
    primary20: "#00000033",
    primary30: "#0000004c",
    primary69: "#000000b0",
    secondary: "#ffffff",
    secondary20: "#ffffff33",
    secondary30: "#ffffff4c",
    secondary69: "#ffffffb0",
    accent: "#d12129",
    accent20: "#d1212933",
    accent30: "#d121294c",
    accent69: "#d12129b0",
  },
  white: {
    primary: "#ffffff",
    primary20: "#ffffff33",
    primary30: "#ffffff4c",
    primary69: "#ffffffb0",
    secondary: "#000000",
    secondary20: "#00000033",
    secondary30: "#0000004c",
    secondary69: "#000000b0",
    accent: "#d12129",
    accent20: "#d1212933",
    accent30: "#d121294c",
    accent69: "#d12129b0",
  },
};

/**
 * @param {string} themeName
 */
function changeTheme(themeName) {
  const theme = themes[themeName];
  if (!theme) {
    return;
  }
  window.Utils.setCookie("theme-name", themeName);
  const style = document.documentElement.style;
  switch (themeName) {
    case "white":
      document.body.style.backgroundImage = `url("/assets/images/shs-bg-logo.webp")`;
      break;
    case "black":
    default:
      document.body.style.backgroundImage = `url("/assets/images/shs-bg-logo-dark.webp")`;
      break;
  }

  style.setProperty("--primary-color", theme.primary);
  style.setProperty("--primary-color-20", theme.primary20);
  style.setProperty("--primary-color-30", theme.primary30);
  style.setProperty("--primary-color-69", theme.primary69);
  style.setProperty("--secondary-color", theme.secondary);
  style.setProperty("--secondary-color-20", theme.secondary20);
  style.setProperty("--secondary-color-30", theme.secondary30);
  style.setProperty("--secondary-color-69", theme.secondary69);
  style.setProperty("--accent-color", theme.accent);
  style.setProperty("--accent-color-20", theme.accent20);
  style.setProperty("--accent-color-30", theme.accent30);
  style.setProperty("--accent-color-69", theme.accent69);
  //document.getElementById("popover-theme-switcher").style.display = "none";
}

(() => {
  const userTheme = window.Utils.getCookie("theme-name");
  if (userTheme) {
    changeTheme(userTheme);
    return;
  }

  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: light)").matches
  ) {
    changeTheme("white");
  }

  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    changeTheme("black");
  }
})();

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", (e) => {
    changeTheme(e.matches ? "black" : "white");
  });

window.Theme = { changeTheme };
