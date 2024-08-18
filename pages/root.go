package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/pages/templates"
)

func RootPage(w http.ResponseWriter, r *http.Request) {
    templates.RootPage().Render(r.Context(), w)
}
