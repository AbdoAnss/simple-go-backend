package handler

import (
	"backend/db"
	"backend/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Récupérer tous les étudiants
func GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT id, name, created_at FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		var createdAt []uint8 // temporaire pour scanner la date
		if err := rows.Scan(&s.ID, &s.Name, &createdAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Convertir en time.Time
		s.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(createdAt))
		students = append(students, s)
	}

	if err := json.NewEncoder(w).Encode(students); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Créer un nouvel étudiant
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student model.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO students (name) VALUES (?)", student.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	student.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// Récupérer les détails d'un étudiant et ses notes
func GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	// Récupérer l'étudiant
	var student model.Student
	var createdAt []uint8
	err := db.DB.QueryRow("SELECT id, name, created_at FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &createdAt)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	student.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", string(createdAt))

	// Récupérer ses notes
	rows, err := db.DB.Query(`
        SELECT id, course_name, note_value 
        FROM notes 
        WHERE student_id = ?`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.Course, &n.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}

	response := struct {
		Student model.Student `json:"student"`
		Notes   []model.Note  `json:"notes"`
	}{
		Student: student,
		Notes:   notes,
	}

	json.NewEncoder(w).Encode(response)
}

// Supprimer
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	// D'abord supprimer les notes associées (à cause de la clé étrangère)
	_, err := db.DB.Exec("DELETE FROM notes WHERE student_id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensuite supprimer l'étudiant
	result, err := db.DB.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student and associated grades deleted successfully",
	})
}
