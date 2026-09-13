// Scroll-reveal for [data-animate] sections (ported from animations.js —
// respects prefers-reduced-motion).
(function () {
  // Marks JS as available so the CSS may hide [data-animate] elements for the
  // scroll-reveal effect (without JS they must stay visible).
  document.documentElement.classList.add("js");

  var elements = document.querySelectorAll("[data-animate]");
  if (!window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    var observer = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            entry.target.classList.add("animate-fade-in");
            observer.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.1 },
    );
    elements.forEach(function (el) {
      observer.observe(el);
    });
  } else {
    elements.forEach(function (el) {
      el.classList.add("animate-fade-in");
    });
  }
})();
