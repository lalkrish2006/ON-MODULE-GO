package handlers

import (
	"encoding/json"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/middleware"
)

// GetStudentDetails returns student name for a given register number
func GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	middleware.RequireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		regNo := r.URL.Query().Get("register_no")
		if regNo == "" {
			http.Error(w, "Missing register_no", http.StatusBadRequest)
			return
		}

		var name, year, dept, section string
		query := "SELECT name, year, department, section FROM students WHERE register_no = ?"
		err := database.DB.QueryRow(query, regNo).Scan(&name, &year, &dept, &section)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    true,
			"name":       name,
			"year":       year,
			"department": dept,
			"section":    section,
		})
	})).ServeHTTP(w, r)
}

// GetMentors returns a list of mentors based on filters
func GetMentors(w http.ResponseWriter, r *http.Request) {
	middleware.RequireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dept := r.URL.Query().Get("department")
		year := r.URL.Query().Get("year")
		section := r.URL.Query().Get("section")

		// Debug print
		// log.Printf("Fetching mentors for: %s %s %s", dept, year, section)

		query := "SELECT name FROM mentors WHERE department = ? AND year = ? AND section = ?"
		rows, err := database.DB.Query(query, dept, year, section)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		mentors := []string{}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				mentors = append(mentors, name)
			}
		}

		json.NewEncoder(w).Encode(mentors)
	})).ServeHTTP(w, r)
}
