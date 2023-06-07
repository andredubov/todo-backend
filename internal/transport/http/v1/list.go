package handler

import "net/http"

func (h *Handler) createList(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getAllList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getAllList"))
}

func (h *Handler) getListByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getListByID"))
}

func (h *Handler) updateListByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteListByID(w http.ResponseWriter, r *http.Request) {

}
