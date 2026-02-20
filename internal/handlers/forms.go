package handlers

import (
	"html/template"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/services"
	"od-system/internal/utils"
)

type ApplyFormData struct {
	User     map[string]interface{}
	Mentors  []string
}

// StudentApply renders the OD application form
func StudentApply(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if role, ok := session.Values["role"].(string); !ok || role != "student" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	// Prepare user map safely
	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	// Fetch Mentors for the current student
	dept, _ := userMap["department"].(string)
	year, _ := userMap["year"].(int)
	section, _ := userMap["section"].(string)

	mentorsQuery := "SELECT name FROM mentors WHERE department = ? AND year = ? AND section = ?"
	rows, err := database.DB.Query(mentorsQuery, dept, year, section)
	var mentors []string
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var mName string
			rows.Scan(&mName)
			mentors = append(mentors, mName)
		}
	} else {
		log.Println("Error fetching mentors for form:", err)
	}

	tmplPath := utils.ResolvePath("templates/student_od_apply.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Template Error:", err)
		http.Error(w, "Template Error", 500)
		return
	}

	data := ApplyFormData{
		User:    userMap,
		Mentors: mentors,
	}

	tmpl.Execute(w, data)
}
