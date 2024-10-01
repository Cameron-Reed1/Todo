package pages

import (
	"math/rand/v2"
	"net/http"
	"os"
	"path"
)


func RandomImage(w http.ResponseWriter, r *http.Request, basePath string) {
    file, err := os.Open(basePath)
    if err != nil {
        return
    }
    defer file.Close()

    files, err := file.Readdir(-1)
    if err != nil {
        return
    }

    idx := rand.UintN(uint(len(files)))

    w.Header().Add("Cache-Control", "max-age=300")
    http.ServeFile(w, r, path.Join(basePath, files[idx].Name()))
}
