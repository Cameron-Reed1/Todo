package pages

import (
	"net/http"

	"github.com/Cameron-Reed1/todo-web/pages/templates"
)

func RootPage(w http.ResponseWriter, r *http.Request) {
    username, err := validateSessionAndGetUsername(r)
    if err != nil {
        w.Header().Add("Location", "/login")
        w.WriteHeader(http.StatusFound)
        return
    }

    templates.RootPage(username, false).Render(r.Context(), w)
}
