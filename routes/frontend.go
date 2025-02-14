package routes

import (
	"net/http"
	"os"
)

func RegisterFrontHandler(w http.ResponseWriter, r *http.Request) {
	root, _ := os.Getwd()
	file, err := os.ReadFile(root + "/routes/html/register.html")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return

	}

	w.Write(file)
}
func LoginFrontHandler(w http.ResponseWriter, r *http.Request) {
	root, _ := os.Getwd()
	file, err := os.ReadFile(root + "/routes/html/login.html")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return

	}

	w.Write(file)
}
func RootHandler(w http.ResponseWriter, r *http.Request) {
	root, _ := os.Getwd()
	file, err := os.ReadFile(root + "/routes/html/index.html")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return

	}

	w.Write(file)
}
