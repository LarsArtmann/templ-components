// Theme sync: keeps the library ThemeToggle's aria-checked in sync and
// follows OS preference changes when the visitor has no stored choice.
// The toggle click behavior itself ships inside the library component.
(function () {
  var media = window.matchMedia("(prefers-color-scheme: dark)");

  function syncToggles() {
    var isDark = document.documentElement.classList.contains("dark");
    document.querySelectorAll("[data-theme-toggle]").forEach(function (btn) {
      btn.setAttribute("aria-checked", String(isDark));
    });
  }

  syncToggles();

  media.addEventListener("change", function () {
    var stored = null;
    try {
      stored = localStorage.getItem("theme");
    } catch (e) {}
    if (!stored) {
      // No explicit choice: mirror the OS like the library ThemeScript does.
      document.documentElement.classList.toggle("dark", media.matches);
      document.documentElement.style.colorScheme = media.matches
        ? "dark"
        : "light";
    }
    syncToggles();
  });
})();
