(function () {
    var html = document.documentElement;
    var saved = localStorage.getItem("theme");

    function applyTheme(theme) {
        if (theme === "dark") {
            html.setAttribute("data-theme", "dark");
        } else if (theme === "light") {
            html.setAttribute("data-theme", "light");
        } else {
            html.removeAttribute("data-theme");
        }
    }

    if (saved === "dark" || saved === "light") {
        applyTheme(saved);
    } else {
        var prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
        applyTheme(prefersDark ? "dark" : "light");
    }

    window.__setTheme = function (theme) {
        applyTheme(theme);
        if (theme === "dark" || theme === "light") {
            localStorage.setItem("theme", theme);
        } else {
            localStorage.removeItem("theme");
        }
    };

    window.__toggleTheme = function () {
        var current = html.getAttribute("data-theme") || "";
        if (current === "dark") {
            window.__setTheme("light");
        } else {
            window.__setTheme("dark");
        }
    };
})();
