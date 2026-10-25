package main

import "net/http"

func loginJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/login.js")
}

func decksJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/decks.js")
}

func editDeckJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/edit-deck.js")
}

func editCardJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/edit-card.js")
}

func studyJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/study.js")
}

func createUserJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/create-user.js")
}

func createDeckJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/create-deck.js")
}

func createCardJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/create-card.js")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/index.html")
}

func decksHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/decks.html")
}

func editDeckHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/edit-deck.html")
}

func editCardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/edit-card.html")
}

func studyHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/study.html")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/create-user.html")
}

func createDeckHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/create-deck.html")
}

func createCardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/create-card.html")
}
