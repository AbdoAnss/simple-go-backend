package handler

import (
	"backend/db"
	"backend/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func AddNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convertir l'ID en entier
	studentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Vérifier que l'étudiant existe
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE id = ?)", studentID).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(
		"INSERT INTO notes (student_id, course_name, note_value) VALUES (?, ?, ?)",
		studentID, note.Course, note.Value,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupérer la note créée pour confirmer
	err = db.DB.QueryRow(
		"SELECT id, course_name, note_value FROM notes WHERE id = ?",
		id,
	).Scan(&note.ID, &note.Course, &note.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}
