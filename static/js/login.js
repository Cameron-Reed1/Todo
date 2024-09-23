function on_load() {
    const login_box = document.getElementById("login-box");
    login_box.addEventListener("htmx:afterRequest", function(evt) {
        if (evt.detail.successful) {
            window.location.pathname = "/";
        }
    });
}



if (document.readyState === "completed") {
    on_load();
} else {
    document.addEventListener("DOMContentLoaded", on_load);
}
