function on_load() {
    let create_start_input = document.getElementById("create-item-start");
    let create_due_input = document.getElementById("create-item-due");

    let create_form = document.getElementById("create-item");

    create_form.addEventListener("htmx:configRequest", function(evt) {
        evt.detail.parameters["start"] = create_start_input.value ? create_start_input.valueAsNumber / 1000 : 0;
        evt.detail.parameters["due"] = create_due_input.value ? create_due_input.valueAsNumber / 1000 : 0;
    });

    create_form.addEventListener("htmx:afterRequest", function(evt) {
        if (evt.detail.successful) {
            window.location.hash = '';
            create_form.reset();
        }
    });

    let edit_start_input = document.getElementById("edit-item-start");
    let edit_due_input = document.getElementById("edit-item-due");

    let edit_form = document.getElementById("edit-item");

    edit_form.addEventListener("htmx:configRequest", function(evt) {
        let target = document.getElementById(edit_form.getAttribute("data-id"));

        evt.detail.parameters["id"] = parseInt(target.getAttribute("data-idnum"));
        evt.detail.parameters["start"] = edit_start_input.value ? edit_start_input.valueAsNumber / 1000 : 0;
        evt.detail.parameters["due"] = edit_due_input.value ? edit_due_input.valueAsNumber / 1000 : 0;
    });

    edit_form.addEventListener("htmx:afterRequest", function(evt) {
        if (evt.detail.successful) {
            window.location.hash = '';
            edit_form.reset();
        }
    });

    edit_form.addEventListener("htmx:beforeSwap", function(evt) {
        let target = document.getElementById(edit_form.getAttribute("data-id"));
        if (target === null || target === undefined) {
            evt.detail.shouldSwap = false;
        } else {
            evt.detail.target = target;
        }
    });

    document.querySelector("#overdue-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let overdue_items = document.querySelector("#overdue-list > .todo-list-items");
        let due = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-due"));
        let id = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-idnum"));

        let target = overdue_items.children[overdue_items.children.length - 1];
        for (let i = 1; i < overdue_items.children.length; i++) {
            let other_due = parseInt(today_items.children[i].getAttribute("data-due"));
            let other_id = parseInt(today_items.children[i].getAttribute("data-idnum"));

            if (other_due > due || (other_due === due && other_id > id)) {
                target = overdue_items.children[i - 1];
                break;
            }
        }

        evt.detail.target = target;

        overdue_items.setAttribute("data-item-count", parseInt(overdue_items.getAttribute("data-item-count")) + 1);
    });

    let checkbox = document.getElementById("show-completed");
    checkbox.addEventListener("change", function(evt) {
        document.documentElement.setAttribute("data-show-completed", checkbox.checked);
    });

    document.querySelector("#today-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let today_items = document.querySelector("#today-list > .todo-list-items");
        let due = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-due"));
        let id = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-idnum"));

        let target = today_items.children[today_items.children.length - 1];
        for (let i = 1; i < today_items.children.length; i++) {
            let other_due = parseInt(today_items.children[i].getAttribute("data-due"));
            let other_id = parseInt(today_items.children[i].getAttribute("data-idnum"));

            if ((other_due > due && due !== 0) || (other_due === due && other_id > id)) {
                target = today_items.children[i - 1];
                break;
            }
        }

        evt.detail.target = target;

        today_items.setAttribute("data-item-count", parseInt(today_items.getAttribute("data-item-count")) + 1);
    });

    document.querySelector("#upcoming-list > .new-item").addEventListener("htmx:oobBeforeSwap", function(evt) {
        let upcoming_items = document.querySelector("#upcoming-list > .todo-list-items");
        let start = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-start"));
        let id = parseInt(evt.detail.fragment.firstElementChild.getAttribute("data-idnum"));

        let target = upcoming_items.children[0];
        for (let i = 1; i < upcoming_items.children.length; i++) {
            let other_start = parseInt(today_items.children[i].getAttribute("data-start"));
            let other_id = parseInt(today_items.children[i].getAttribute("data-idnum"));

            if (other_start > start || (other_start === start && other_id > id)) {
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



document.addEventListener("keyup", (event) => {
    if (event.key === "Escape") {
        window.location.hash = '';
    }
});



function edit(item_id) {
    let item = document.getElementById(item_id);
    let start = parseInt(item.getAttribute("data-start"));
    let due = parseInt(item.getAttribute("data-due"));
    let text = item.querySelector(".todo-text").textContent;

    let edit_form = document.getElementById("edit-item");
    document.getElementById("edit-item-name").value = text;
    if (start != 0) {
        document.getElementById("edit-item-start").valueAsNumber = start * 1000;
    } else {
        document.getElementById("edit-item-start").value = null;
    }
    if (due != 0) {
        document.getElementById("edit-item-due").valueAsNumber = due * 1000;
    } else {
        document.getElementById("edit-item-due").value = null;
    }
    //document.getElementByid("edit-item-start").value = new Date(start * 100).toISOString().slice(0, 16);
    //document.getElementByid("edit-item-due").value = new Date(due * 100).toISOString().slice(0, 16);
    edit_form.setAttribute("data-id", item_id);
    window.location.hash = "edit-item";
}

