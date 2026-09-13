// Hero copy button (ported from copy-code.js).
(function () {
  var btn = document.getElementById("copy-btn");
  if (!btn) return;

  btn.addEventListener("click", function () {
    var code = btn.getAttribute("data-code");
    navigator.clipboard
      .writeText(code || "")
      .then(function () {
        btn.textContent = "Copied!";
        btn.classList.add("text-success");
        setTimeout(function () {
          btn.textContent = "Copy";
          btn.classList.remove("text-success");
        }, 2000);
      })
      .catch(function () {
        btn.textContent = "Failed";
        setTimeout(function () {
          btn.textContent = "Copy";
        }, 2000);
      });
  });
})();
