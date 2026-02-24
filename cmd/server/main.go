package main

import (
	"log"
	"net/http"
	"od-system/internal/config"
	"od-system/internal/database"
	"od-system/internal/handlers"
	"od-system/internal/services"
	"od-system/internal/utils"
	"os"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Session Store
	services.Init()

	// Connect to Database
	database.Connect(cfg)
	defer database.DB.Close()

	// Setup Router (Using standard ServeMux for now)
	mux := http.NewServeMux()

	// Static Files
	staticPath := utils.ResolvePath("static")
	fileServer := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Register Handlers
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/student/dashboard", handlers.StudentDashboard)
	mux.HandleFunc("/student/apply", handlers.StudentApply)
	mux.HandleFunc("/api/student", handlers.GetStudentDetails)
	mux.HandleFunc("/api/mentors", handlers.GetMentors)
	mux.HandleFunc("/student/submit", handlers.SubmitOD)
	
	mux.HandleFunc("/mentor/dashboard", handlers.MentorDashboard)
	mux.HandleFunc("/mentor/action", handlers.MentorAction)

	mux.HandleFunc("/hod/dashboard", handlers.HODDashboard)
	mux.HandleFunc("/hod/action", handlers.HODAction)

	mux.HandleFunc("/principal/dashboard", handlers.PrincipalDashboard)
	mux.HandleFunc("/principal/action", handlers.PrincipalAction)
	
	
	// LabTech routes
	mux.HandleFunc("/labtech/dashboard", handlers.LabTechDashboard)
	mux.HandleFunc("/labtech/action", handlers.LabTechAction)

	// New Dashboards
	mux.HandleFunc("/ca/dashboard", handlers.CADashboard)
	mux.HandleFunc("/ja/dashboard", handlers.JADashboard)
	
	// Admin Routes
	mux.HandleFunc("/admin/dashboard", handlers.AdminDashboard)
	mux.HandleFunc("/admin/update", handlers.AdminUpdateUser)
	mux.HandleFunc("/admin/add", handlers.AdminAddUser)
	mux.HandleFunc("/admin/delete", handlers.AdminDeleteUser)
	mux.HandleFunc("/admin/view_ods", handlers.AdminViewODs)

	// Default redirect to login for root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Default port changed to 8082
	}
	serverAddress := ":" + port
	log.Printf("Server starting on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
