// Scroll-reveal for [data-animate] sections (ported from animations.js —
// respects prefers-reduced-motion).
(function () {
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
