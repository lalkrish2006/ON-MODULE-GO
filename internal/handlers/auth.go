package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/services"
	"od-system/internal/utils"
)

// LoginPageData holds data for rendering the login page
type LoginPageData struct {
	Error     string
	LogoutMsg bool
}

// Login serves the login page (GET) and handles authentication (POST)
func Login(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		role, _ := session.Values["role"].(string)
		redirectMap := map[string]string{
			"student":   "/student/dashboard",
			"mentor":    "/mentor/dashboard",
			"hod":       "/hod/dashboard",
			"ca":        "/ca/dashboard",
			"ja":        "/ja/dashboard",
			"principal": "/principal/dashboard",
			"labtech":   "/labtech/dashboard",
			"admin":     "/admin/dashboard",
		}
		if target, exists := redirectMap[role]; exists {
			http.Redirect(w, r, target, http.StatusSeeOther)
			return
		}
	}

	// Prevent Caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method == http.MethodGet {
		tmplPath := utils.ResolvePath("templates/login.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Printf("Template error: open %s: %v", tmplPath, err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		data := LoginPageData{
			Error:     r.URL.Query().Get("error"),
			LogoutMsg: r.URL.Query().Get("msg") == "loggedout",
		}
		if data.Error == "unauthorized" {
			data.Error = "❌ Unauthorized access! Please login."
		} else if data.Error == "invalid_role" {
			data.Error = "❌ Invalid Role!"
		}

		tmpl.Execute(w, data)
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Form parse error", http.StatusBadRequest)
			return
		}

		role := r.FormValue("role")
		registerNo := r.FormValue("register_no")
		password := r.FormValue("password")

		// Map role to table
		tableMap := map[string]string{
			"student":   "students",
			"mentor":    "mentors",
			"ca":        "cas",
			"ja":        "jas",
			"hod":       "hods",
			"principal": "principals",
			"labtech":   "lab_technicians",
			"admin":     "admin_users",
		}

		table, ok := tableMap[role]
		if !ok {
			http.Redirect(w, r, "/login?error=invalid_role", http.StatusSeeOther)
			return
		}

		// Query DB
		// Note: We scan into generic map/interace or struct fields needed for session
		// For simplicity, we just fetch name, password, department, year, section if available
		var storedPassword, name, department, section string
		var year int


		// Prepare query based on columns available in all tables?
		// No, columns differ. We need to be dynamic or specific.
		// Common: register_no, password, name.
		// Others: department (most), year (student, mentor, ca), section (student, mentor, ca).

		// Let's use a dynamic approach or struct specific queries is safer but verbose.
		// Given we just need to validate and set session, we can `SELECT *` and scan needed fields or specific cols.
		// But Scan needs exact count.
		// Safer approach: Query generic fields for password check, then extra fields for session.

		// Simplified: Select password and common fields.
		// Actually, let's use the models? But we don't know which one until runtime if we want strongly typed.
		// We can switch on role again for querying.
		
		query := "SELECT password, name FROM " + table + " WHERE register_no = ?"
		err = database.DB.QueryRow(query, registerNo).Scan(&storedPassword, &name)
		
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?error=not_found", http.StatusSeeOther)
			return
		} else if err != nil {
			log.Println("DB Error:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Check Password
		if utils.HashPasswordMD5(password) != storedPassword {
			// Try plain text for temporary debugging or known issue? No, strict matching user req.
			http.Redirect(w, r, "/login?error=invalid_password", http.StatusSeeOther)
			return
		}

		// Fetch additional details for session if needed (dept, year, etc.)
		// This is used in dashboards.
		if role == "student" || role == "mentor" || role == "ca" {
			// Fetch extra fields
			_ = database.DB.QueryRow("SELECT department, year, section FROM "+table+" WHERE register_no = ?", registerNo).Scan(&department, &year, &section)
		} else if role != "principal" { // HOD, Admin, Labtech, JA have department
             // Principal table might not have department? Schema analysis needed.
             // PHP code: `if (isset($user['department'])) $_SESSION['department'] = ...`
			_ = database.DB.QueryRow("SELECT department FROM "+table+" WHERE register_no = ?", registerNo).Scan(&department)
		}

		// Set Session
		session := services.GetSession(r)
		session.Values["authenticated"] = true
		session.Values["role"] = role
		session.Values["register_no"] = registerNo
		session.Values["name"] = name
		session.Values["department"] = department
		session.Values["year"] = year
		session.Values["section"] = section
		session.Save(r, w)

		// Redirect
		redirectMap := map[string]string{
			"student":   "/student/dashboard",
			"mentor":    "/mentor/dashboard",
			"hod":       "/hod/dashboard",
			"ca":        "/ca/dashboard",
			"ja":        "/ja/dashboard",
			"principal": "/principal/dashboard",
			"labtech":   "/labtech/dashboard",
			"admin":     "/admin/dashboard",
		}
		
		if target, ok := redirectMap[role]; ok {
			http.Redirect(w, r, target, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login?error=invalid_role", http.StatusSeeOther)
		}
	}
}

// Logout clears the session
func Logout(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // delete immediately
	session.Save(r, w)
	http.Redirect(w, r, "/login?msg=loggedout", http.StatusSeeOther)
}
