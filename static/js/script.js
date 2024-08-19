document.addEventListener("keyup", (event) => {
    if (event.key === "Escape") {
        window.location.hash = '';
    }
});

function on_load() {
    let start_input = document.getElementById("create-item-form-start");
    let due_input = document.getElementById("create-item-form-due");

    let create_form = document.getElementById("create-item");

    create_form.addEventListener("htmx:configRequest", function(evt) {
        evt.detail.parameters["start"] = start_input.value ? start_input.valueAsNumber / 1000 : 0;
        evt.detail.parameters["due"] = due_input.value ? due_input.valueAsNumber / 1000 : 0;
    });

    create_form.addEventListener("htmx:afterRequest", function(evt) {
        if (evt.detail.successful) {
            window.location.hash = '';
            evt.detail.elt.reset();
        }
    });

    document.querySelector("#overdue-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let overdue_items = document.querySelector("#overdue-list > .todo-list-items");
        let due = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-due"));

        let target = overdue_items.children[overdue_items.children.length - 1];
        for (let i = 1; i < overdue_items.children.length; i++) {
            if (parseInt(overdue_items.children[i].getAttribute("data-due")) > due) {
                target = overdue_items.children[i - 1];
                break;
            }
        }

        evt.detail.target = target;

        overdue_items.setAttribute("data-item-count", parseInt(overdue_items.getAttribute("data-item-count")) + 1);
    });

    document.querySelector("#today-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let today_items = document.querySelector("#today-list > .todo-list-items");
        let due = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-due"));

        let target = today_items.children[today_items.children.length - 1];
        if (due !== 0) {
            for (let i = 1; i < today_items.children.length; i++) {
                if (parseInt(today_items.children[i].getAttribute("data-due")) > due) {
                    target = today_items.children[i - 1];
                    break;
                }
            }
        }

        evt.detail.target = target;

        today_items.setAttribute("data-item-count", parseInt(today_items.getAttribute("data-item-count")) + 1);
    });

    document.querySelector("#upcoming-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let upcoming_items = document.querySelector("#upcoming-list > .todo-list-items");
        let start = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-start"));

        let target = upcoming_items.children[0];
        for (let i = 1; i < upcoming_items.children.length; i++) {
            if (parseInt(upcoming_items.children[i].getAttribute("data-start")) > start) {
                target = upcoming_items.children[i - 1];
                break;
            }
        }

        evt.detail.target = target;

        upcoming_items.setAttribute("data-item-count", parseInt(upcoming_items.getAttribute("data-item-count")) + 1);
    });
}


if (document.readyState === "completed") {
    on_load();
} else {
    document.addEventListener("DOMContentLoaded", on_load);
}

