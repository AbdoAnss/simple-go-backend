package main

import (
	"backend/db"
	"backend/router"
	"fmt"
	"net/http"
)

func main() {
	// Initialisation de la base de données
	db.InitDB()

	// Configuration du routeur
	r := router.SetupRouter()

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur: 8080...")
	http.ListenAndServe(":8080", r)
}
