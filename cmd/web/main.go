package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"unicode/utf8"

	"ascii-art/modules"
)

// список доступных шаблонов и файлы с ними
var (
	templateFiles = map[string]string{
		"Standard":   "banners/standard.txt",
		"Shadow":     "banners/shadow.txt",
		"Thinkertoy": "banners/thinkertoy.txt",
	}
	// тут будут уже загруженные баннер-карты
	banners = map[string]map[rune][]string{}
	// порядок вывода радио-кнопок
	templateOrder = []string{"Standard", "Shadow", "Thinkertoy"}
	tpl           = template.Must(template.ParseFiles("cmd/web/templates/index.html"))
	errorTpl      = template.Must(template.ParseFiles("cmd/web/templates/error.html"))
)

func main() {
	// Загружаем все баннеры в память
	for name, relPath := range templateFiles {
		abs, err := filepath.Abs(relPath)
		if err != nil {
			log.Fatalf("Abs path error %s: %v", relPath, err)
		}
		m, err := modules.LoadBanner(abs)
		if err != nil {
			log.Fatalf("Couldn't load banner %s: %v", relPath, err)
		}
		banners[name] = m
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("cmd/web/statics"))
	mux.Handle("/statics/", http.StripPrefix("/statics/", fs))
	mux.HandleFunc("/ascii-art", handleGenerate)
	mux.HandleFunc("/", rootHandler)

	log.Println("Server deployed on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// не ровно корень — 404
		renderError(w, http.StatusNotFound, "Not found",
			"Requested page's not found.")
		return
	}

	if r.Method != http.MethodGet {
		// корень, но не GET — Method Not Allowed
		renderError(w, http.StatusMethodNotAllowed, "Method's not supported",
			"Use GET to access this page.")
		return
	}

	// всё ок, рендерим главную
	handleIndex(w, r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "Method's not supported",
			"Use GET to access this page.")
		return
	}

	data := struct {
		Templates     []string
		Selected      string
		InputText     string
		RenderedASCII string
	}{
		Templates:     templateOrder,
		Selected:      "Standard", // шаблон по умолчанию
		InputText:     "",
		RenderedASCII: "",
	}
	if err := tpl.Execute(w, data); err != nil {
		log.Printf("render index error: %v", err)
		renderError(w, http.StatusInternalServerError, "Server error",
			"Couldn't upload main page.")
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed, "Method's not supported",
			"Use POST generate ASCII Art.")
		return
	}
	if err := r.ParseForm(); err != nil {
		renderError(w, http.StatusBadRequest, "Invalid request",
			"Couldn't read sended data.")
		return
	}

	text := r.FormValue("inputText")
	tmpl := r.FormValue("template")

	if utf8.RuneCountInString(text) > 10000 {
		renderError(w, http.StatusBadRequest,
			"Bad Request",
			"Max length of symbols is 10000")
		return
	}

	if !isASCII(text) {
		renderError(w, http.StatusBadRequest,
			"Bad Request",
			"Only ASCII symbols allowed.")
		return
	}

	bannerMap, ok := banners[tmpl]
	if !ok {
		renderError(w, http.StatusNotFound, "Template's not found",
			"Selected banner's not exist.")
		return
	}

	art := modules.RenderBanner(text, bannerMap)

	data := struct {
		Templates     []string
		Selected      string
		InputText     string
		RenderedASCII string
	}{
		Templates:     templateOrder,
		Selected:      tmpl,
		InputText:     text,
		RenderedASCII: art,
	}
	if err := tpl.Execute(w, data); err != nil {
		log.Printf("render ascii-art error: %v", err)
		renderError(w, http.StatusInternalServerError, "Server error",
			"Couldn't show the result.")
	}
}

func renderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	data := struct {
		Code    int
		Title   string
		Message string
	}{
		Code:    code,
		Title:   title,
		Message: message,
	}
	if err := errorTpl.Execute(w, data); err != nil {
		// если даже шаблон ошибки сломался, возвращаем plain-text
		http.Error(w, title+": "+message, code)
	}
}
