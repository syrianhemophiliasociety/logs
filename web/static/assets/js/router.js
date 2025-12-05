"use strict";

const mainContentsEl = document.getElementById("main-contents");

const links = [
  {
    check: (l) => l === "/",
    elements: [
      document.getElementById("/"),
      document.getElementById("/?mobile"),
    ],
  },
  {
    check: (l) => l.startsWith("/patient"),
    elements: [
      document.getElementById("/patients"),
      document.getElementById("/patients?mobile"),
    ],
  },
  {
    check: (l) => l === "/blood-tests",
    elements: [
      document.getElementById("/blood-tests"),
      document.getElementById("/blood-tests?mobile"),
    ],
  },
  {
    check: (l) => l === "/viruses",
    elements: [
      document.getElementById("/viruses"),
      document.getElementById("/viruses?mobile"),
    ],
  },
  {
    check: (l) => l === "/medicines",
    elements: [
      document.getElementById("/medicines"),
      document.getElementById("/medicines?mobile"),
    ],
  },
  {
    check: (l) => l.startsWith("/management"),
    elements: [
      document.getElementById("/management"),
      document.getElementById("/management?mobile"),
    ],
  },
];

function updateActiveNavLink() {
  for (const link of links) {
    if (link.check(window.location.pathname)) {
      link.elements.forEach((e) => e?.classList.add("bg-accent-trans-20"));
    } else {
      link.elements.forEach((e) => e?.classList.remove("bg-accent-trans-20"));
    }
  }
}

/**
 * @param {string} path the requested path to update.
 */
async function updateMainContent(path) {
  Utils.showLoading();
  const query = new URLSearchParams(location.search);
  query.set("no_layout", "true");
  htmx
    .ajax("GET", path + "?" + query.toString(), {
      target: "#main-contents",
      swap: "innerHTML",
    })
    .catch(() => {
      window.location.reload();
    })
    .finally(() => {
      Utils.hideLoading();
      updateActiveNavLink();
    });
}

function updateSearchQuery(key, value) {
  const query = new URLSearchParams(location.search);
  query.set(key, value);

  if (history.pushState) {
    const newurl = `${location.protocol}//${location.host}${location.pathname}?${query.toString()}`;
    window.history.pushState({ path: newurl }, "", newurl);
  }
}

function removeSearchQuery(key) {
  const query = new URLSearchParams(location.search);
  query.delete(key);

  if (history.pushState) {
    const newurl = `${location.protocol}//${location.host}${location.pathname}?${query.toString()}`;
    window.history.pushState({ path: newurl }, "", newurl);
  }
}

window.addEventListener("load", () => {
  updateActiveNavLink();
});

window.addEventListener("popstate", async (e) => {
  const mainContentsEl = document.getElementById("main-contents");
  if (!!mainContentsEl && !!e.target.location.pathname) {
    e.stopImmediatePropagation();
    e.preventDefault();
    await updateMainContent(e.target.location.pathname);
    return;
  }
});

document.addEventListener("htmx:afterRequest", function (e) {
  if (!!e.detail && !!e.detail.xhr) {
    const newTitle = e.detail.xhr.getResponseHeader("HX-Title");
    if (newTitle) {
      document.title = newTitle + " - SyrianHemophiliaSocietyLogs";
    }
  }
});

window.Router = {
  updateActiveNavLink,
  updateMainContent,
  updateSearchQuery,
  removeSearchQuery,
};
