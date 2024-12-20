package router

import (
	"backend/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Peut être remplacé par un tableau de domaines autorisés
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes pour les étudiants
	r.Get("/students", handler.GetStudents)
	r.Post("/students", handler.CreateStudent)
	r.Get("/students/{id}", handler.GetStudentDetails)
	r.Delete("/students/{id}", handler.DeleteStudent)

	// Routes pour les notes
	r.Post("/students/{id}/notes", handler.AddNote)

	return r
}
