// Mobile navigation toggle (ported from the Astro site's header.js; the
// theme-toggle half moved to theme-sync.js / the library component).
(function () {
  var toggle = document.getElementById("nav-toggle");
  var navLinks = document.getElementById("nav-links");
  if (!toggle || !navLinks) return;

  var menuIcon = toggle.querySelector(".menu-icon");
  var closeIcon = toggle.querySelector(".close-icon");

  function setOpen(open) {
    navLinks.classList.toggle("open", open);
    if (menuIcon) menuIcon.classList.toggle("hidden", open);
    if (closeIcon) closeIcon.classList.toggle("hidden", !open);
    toggle.setAttribute("aria-expanded", String(open));
  }

  toggle.addEventListener("click", function () {
    setOpen(!navLinks.classList.contains("open"));
  });

  navLinks.querySelectorAll(".nav-link").forEach(function (link) {
    link.addEventListener("click", function () {
      setOpen(false);
    });
  });
})();
