/* Header docs search: lazily fetches /search-index.json, filters it
   client-side, and navigates on selection. CSP-safe: external file, no
   inline handlers, no third-party requests. */
(function () {
  "use strict";

  var input = document.getElementById("doc-search-input");
  var panel = document.getElementById("doc-search-results");
  if (!input || !panel) return;

  var MAX_HITS = 8;
  var MIN_QUERY_LENGTH = 2;
  var SNIPPET_BEFORE = 40;
  var SNIPPET_AFTER = 70;

  var index = null;
  var activeHit = -1;

  function loadIndex() {
    if (index) return Promise.resolve(index);

    if (!loadIndex.promise) {
      loadIndex.promise = fetch("/search-index.json")
        .then(function (response) {
          return response.json();
        })
        .then(function (parsed) {
          index = parsed;
          return index;
        })
        .catch(function () {
          loadIndex.promise = null;
          return [];
        });
    }

    return loadIndex.promise;
  }

  function scoreDoc(doc, terms) {
    var title = doc.title.toLowerCase();
    var everything = (doc.title + " " + (doc.description || "") + " " + doc.body).toLowerCase();
    var score = 0;
    var matchedSection = null;

    for (var i = 0; i < terms.length; i++) {
      var term = terms[i];
      var inTitle = title.indexOf(term);
      var inEverything = everything.indexOf(term);

      if (inTitle === -1 && inEverything === -1) return null;

      if (inTitle === 0) score += 40;
      else if (inTitle !== -1) score += 25;

      if (inEverything !== -1) score += 10;

      var sections = doc.sections || [];
      for (var s = 0; s < sections.length; s++) {
        if (sections[s].text.toLowerCase().indexOf(term) !== -1) {
          score += 15;
          if (!matchedSection) matchedSection = sections[s];
          break;
        }
      }
    }

    return { doc: doc, score: score, section: matchedSection, terms: terms };
  }

  function snippetFor(hit) {
    if (hit.section) return hit.section.text;

    var body = hit.doc.body;
    var pos = -1;

    for (var i = 0; i < hit.terms.length && pos === -1; i++) {
      pos = body.toLowerCase().indexOf(hit.terms[i]);
    }

    if (pos === -1) return (hit.doc.description || "").slice(0, SNIPPET_BEFORE + SNIPPET_AFTER);

    var start = Math.max(0, pos - SNIPPET_BEFORE);
    var end = Math.min(body.length, pos + SNIPPET_AFTER);
    var snippet = body.slice(start, end);
    if (start > 0) snippet = "…" + snippet;
    if (end < body.length) snippet += "…";

    return snippet;
  }

  function appendMarked(element, text, terms) {
    var lower = text.toLowerCase();
    var pos = -1;
    var length = 0;

    for (var i = 0; i < terms.length; i++) {
      var at = lower.indexOf(terms[i]);
      if (at !== -1 && (pos === -1 || at < pos)) {
        pos = at;
        length = terms[i].length;
      }
    }

    if (pos === -1) {
      element.textContent = text;
      return;
    }

    element.appendChild(document.createTextNode(text.slice(0, pos)));

    var mark = document.createElement("mark");
    mark.textContent = text.slice(pos, pos + length);
    element.appendChild(mark);

    element.appendChild(document.createTextNode(text.slice(pos + length)));
  }

  function render(query, results) {
    panel.textContent = "";
    activeHit = -1;
    input.removeAttribute("aria-activedescendant");

    if (!results.length) {
      var empty = document.createElement("div");
      empty.className = "doc-search-no";
      empty.textContent = "No results for “" + query + "”";
      panel.appendChild(empty);
      open();
      return;
    }

    results.forEach(function (hit, i) {
      var link = document.createElement("a");
      link.className = "doc-search-hit";
      link.href = hit.url + (hit.section ? "#" + hit.section.id : "");
      link.setAttribute("role", "option");
      link.id = "doc-search-option-" + i;

      var title = document.createElement("span");
      title.className = "doc-search-hit-title";
      appendMarked(title, hit.doc.title, hit.terms);

      var detail = document.createElement("span");
      detail.className = "doc-search-hit-sub";
      appendMarked(detail, snippetFor(hit), hit.terms);

      link.appendChild(title);
      link.appendChild(detail);
      panel.appendChild(link);
    });

    open();
  }

  function open() {
    panel.classList.remove("hidden");
    input.setAttribute("aria-expanded", "true");
  }

  function close() {
    panel.classList.add("hidden");
    input.setAttribute("aria-expanded", "false");
    activeHit = -1;
    input.removeAttribute("aria-activedescendant");
  }

  function setActive(next) {
    var options = panel.querySelectorAll(".doc-search-hit");
    if (!options.length) return;

    if (activeHit >= 0 && options[activeHit]) options[activeHit].classList.remove("active");

    activeHit = (next + options.length) % options.length;

    options[activeHit].classList.add("active");
    options[activeHit].scrollIntoView({ block: "nearest" });
    input.setAttribute("aria-activedescendant", options[activeHit].id);
  }

  function runSearch(query) {
    var terms = query.toLowerCase().split(/\s+/).filter(Boolean);
    var results = [];

    for (var i = 0; i < index.length; i++) {
      var hit = scoreDoc(index[i], terms);
      if (hit) results.push(hit);
    }

    results.sort(function (a, b) {
      return b.score - a.score;
    });

    render(query, results.slice(0, MAX_HITS));
  }

  input.addEventListener("input", function () {
    var query = input.value.trim();

    if (query.length < MIN_QUERY_LENGTH) {
      close();
      return;
    }

    loadIndex().then(function () {
      if (input.value.trim() !== query) return;
      runSearch(query);
    });
  });

  input.addEventListener("keydown", function (event) {
    var options = panel.querySelectorAll(".doc-search-hit");

    switch (event.key) {
      case "ArrowDown":
        if (options.length) {
          event.preventDefault();
          setActive(activeHit + 1);
        }
        break;
      case "ArrowUp":
        if (options.length) {
          event.preventDefault();
          setActive(activeHit === -1 ? options.length - 1 : activeHit - 1);
        }
        break;
      case "Enter": {
        var target = activeHit >= 0 ? options[activeHit] : options[0];
        if (target) {
          event.preventDefault();
          window.location.assign(target.href);
        }
        break;
      }
      case "Escape":
        input.value = "";
        close();
        break;
    }
  });

  panel.addEventListener("pointerdown", function (event) {
    var hit = event.target.closest(".doc-search-hit");
    if (hit) {
      event.preventDefault();
      window.location.assign(hit.href);
    }
  });

  document.addEventListener("click", function (event) {
    if (event.target !== input && !panel.contains(event.target)) close();
  });
})();
