// Docs code-copy buttons: binds every .code-copy button (rendered by the
// markdown engine's highlighting wrapper) to its sibling <pre> content.
(function () {
  document.querySelectorAll(".code-block .code-copy").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var pre = btn.parentElement.querySelector("pre");
      if (!pre) return;

      navigator.clipboard
        .writeText(pre.textContent || "")
        .then(function () {
          btn.textContent = "Copied!";
          btn.classList.add("code-copy-done");
          setTimeout(function () {
            btn.textContent = "Copy";
            btn.classList.remove("code-copy-done");
          }, 2000);
        })
        .catch(function () {
          btn.textContent = "Failed";
          setTimeout(function () {
            btn.textContent = "Copy";
          }, 2000);
        });
    });
  });
})();
