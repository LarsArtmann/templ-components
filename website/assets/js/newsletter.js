// Newsletter popup: opens the Buttondown confirmation window on submit.
// Replaces the Astro site's inline onsubmit="" handler so pages stay
// CSP-safe (no inline event handlers).
(function () {
  var form = document.getElementById("newsletter-form");
  if (!form) return;

  form.addEventListener("submit", function () {
    window.open(
      "about:blank",
      "popupwindow",
      "scrollbars=yes,width=560,height=540",
    );
  });
})();
