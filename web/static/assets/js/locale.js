function changeLocale(localeKey) {
  window.Utils.setCookie("locale", localeKey);
  window.location.reload();
}

window.Locale = { changeLocale };
