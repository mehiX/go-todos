package todos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
}

func listTodos(svc Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		all, err := svc.ListAll(r.Context())
		if err != nil {
			log.Printf("Serving all: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(all); err != nil {
			log.Printf("Encoding all: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
		}
	}
}

func createTodo(svc Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var td Todo
		if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
			log.Printf("Decoding body to create: %v\n", err)
			handleError(w, err, http.StatusBadRequest)
			return
		}

		newTd, err := svc.Add(r.Context(), td)
		if err != nil {
			log.Printf("Creating todo: %v\n", err)
			handleError(w, fmt.Errorf("todo not saved"), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(newTd); err != nil {
			log.Printf("Encoding new todo: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
		}
	}
}

func getTodo(svc Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		t, ok := r.Context().Value(TodoCtxKey).(*Todo)
		if !ok || t == nil {
			log.Println("no todo from request context")
			handleError(w, fmt.Errorf("not found"), http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(t); err != nil {
			log.Printf("Encoding todo: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
		}
	}
}

func deleteTodo(svc Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		t, ok := r.Context().Value(TodoCtxKey).(*Todo)
		if !ok || t == nil {
			log.Println("no todo from request context")
			handleError(w, fmt.Errorf("not found"), http.StatusNotFound)
			return
		}

		if err := svc.Delete(r.Context(), t.ID); err != nil {
			log.Printf("Deleting todo: %v\n", err)
			handleError(w, fmt.Errorf("todo not deleted"), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "{}")
	}
}

func updateTodo(svc Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		old, ok := r.Context().Value(TodoCtxKey).(*Todo)
		if !ok || old == nil {
			log.Println("no todo from request context")
			handleError(w, fmt.Errorf("not found"), http.StatusNotFound)
			return
		}

		var newTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			log.Printf("decoding body for update: %v\n", err)
			handleError(w, err, http.StatusBadRequest)
			return
		}

		updated, err := svc.Update(r.Context(), old.ID, newTodo)
		if err != nil {
			log.Printf("updating todo: %v\n", err)
			handleError(w, fmt.Errorf("todo not updated"), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(updated); err != nil {
			log.Printf("Encoding new todo: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
		}
	}
}

func searchByTag(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		tags := strings.Split(r.URL.Query().Get("q"), ",")
		if len(tags) == 0 {
			handleError(w, fmt.Errorf("provide a list of tags in the q query parameter"), http.StatusBadRequest)
			return
		}

		withTag, err := svc.FindByTags(r.Context(), tags)
		if err != nil {
			log.Printf("searching by tags: %v\n", err)
			handleError(w, fmt.Errorf("error searching by tags"), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(withTag); err != nil {
			log.Printf("Encoding all: %v\n", err)
			handleError(w, err, http.StatusInternalServerError)
		}
	}
}
