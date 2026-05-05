window.nexusTheme = {
    getTheme: function () {
        return localStorage.getItem("theme") || "light";
    },

    setTheme: function (theme) {
        localStorage.setItem("theme", theme);

        if (theme === "dark") {
            document.documentElement.classList.add("dark");
            document.body.classList.add("dark");
        } else {
            document.documentElement.classList.remove("dark");
            document.body.classList.remove("dark");
        }
    },

    toggleTheme: function () {
        const current = localStorage.getItem("theme") || "light";
        const next = current === "dark" ? "light" : "dark";
        this.setTheme(next);
        return next;
    }
};