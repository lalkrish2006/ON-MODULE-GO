package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strings"
	"time"
)

// StudentDashboardData holds data for the student dashboard
type StudentDashboardData struct {
	User        map[string]interface{}
	ODs         []DashboardOD
	FlashSuccess string
}

// DashboardOD extends ODApplication with display fields
type DashboardOD struct {
	models.ODApplication
	DateStr       string
	DisplayStatus string
	DisplayText   string
	BadgeClass    string
	TeamMembers   []models.ODTeamMember
	IsIndividual  bool
	TeamMembersJSON string // Added for JSON serialization in templates
}

// StudentDashboard handler
func StudentDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if role, ok := session.Values["role"].(string); !ok || role != "student" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	regNo, _ := session.Values["register_no"].(string)

	// Fetch OD Applications (Created by student OR where they are a member)
	query := `
		SELECT DISTINCT ` + ODColumns + ` 
		FROM od_applications o
		LEFT JOIN od_team_members t ON o.id = t.od_id
		WHERE o.register_no = ? OR t.member_regno = ?
		ORDER BY o.id DESC`

	rows, err := database.DB.Query(query, regNo, regNo)
	if err != nil {
		log.Println("Error fetching ODs:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dashboardODs []DashboardOD

	for rows.Next() {
		var od models.ODApplication
		err := rows.Scan(
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section, 
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate, 
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide, 
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			log.Println("Scan Error:", err)
			continue
		}
		
		// Normalize ODType to lowercase for consistent logic
		od.ODType = strings.ToLower(od.ODType)

		// Fetch Team Members
		tmQuery := "SELECT id, od_id, member_name, member_regno, member_department, member_year, member_section, mentor, mentor_status FROM od_team_members WHERE od_id = ?"
		tmRows, err := database.DB.Query(tmQuery, od.ID)
		var teamMembers []models.ODTeamMember
		if err == nil {
			for tmRows.Next() {
				var m models.ODTeamMember
				tmRows.Scan(&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberDepartment, &m.MemberYear, &m.MemberSection, &m.Mentor, &m.MentorStatus)
				teamMembers = append(teamMembers, m)
			}
			tmRows.Close()
		}

		// Logic for Display Status
		displayStatus := od.Status // Valid: pending, Mentor Accepted, Hod Accepted, etc.
		// Check if I am rejected in team?
		for _, m := range teamMembers {
			if m.MemberRegNo == regNo {
				status := ""
				if m.MentorStatus.Valid {
					status = m.MentorStatus.String
				}
				if status == "Rejected" && od.Status == "HOD Accepted" {
					displayStatus = "mentor rejected" // Legacy PHP logic quirk
				} else if status == "Rejected" {
                     // If I am rejected specifically?
                }
				break
			}
		}

		// Badge Class
		badgeClass := "bg-secondary"
		switch displayStatus {
		case "pending":
			badgeClass = "bg-warning text-dark"
		case "Mentor Accepted", "HOD Accepted", "Principal Accepted":
			badgeClass = "bg-success"
		case "Mentor Rejected", "HOD Rejected", "Principal Rejected", "mentor rejected":
			badgeClass = "bg-danger"
		}

		// Date Formatting
		dateStr := "-"
		
		// Helper functions for formatting (inline)
		formatDate := func(d string) string {
			d = strings.TrimSpace(d) // Handle potential whitespace
			// Handle "YYYY-MM-DDTHH:MM:SSZ" or "YYYY-MM-DD"
			if len(d) > 10 {
				log.Println("Date truncated:", d, "to", d[:10])
				d = d[:10] // Extract YYYY-MM-DD part if it's a timestamp
			}
			t, err := time.Parse("2006-01-02", d)
			if err != nil { 
				log.Println("Date Parse Error:", err, "for", d)
				return d 
			}
			if t.Year() <= 1 {
				return "-"
			}
			return t.Format("2 Jan 2006")
		}
		formatTime := func(tStr string) string {
			// Try parsing with seconds first, then without
			t, err := time.Parse("15:04:05", tStr)
			if err != nil {
				t, err = time.Parse("15:04", tStr)
				if err != nil { return tStr }
			}
			return t.Format("3:04 pm")
		}

		// Helper to check if date is valid non-zero
		isValidDate := func(ns sql.NullString) bool {
			return ns.Valid && ns.String != "0000-00-00" && len(ns.String) > 0
		}

		if od.ODType == "internal" {
			// Check if we have valid times that are not zero (00:00:00)
			hasTime := od.FromTime.Valid && od.ToTime.Valid && od.FromTime.String != "00:00:00" && od.ToTime.String != "00:00:00"
			
			// Case 2: Period-wise (Time + Date)
			if hasTime && isValidDate(od.ODDate) {
				dateStr = formatTime(od.FromTime.String) + " to " + formatTime(od.ToTime.String) + ", " + formatDate(od.ODDate.String)
			} else if isValidDate(od.FromDate) && isValidDate(od.ToDate) {
				// Case 3: More than a day (or technically single day range)
				f := formatDate(od.FromDate.String)
				t := formatDate(od.ToDate.String)
				if f == t {
					dateStr = f // Display as single date if they are same
				} else {
					dateStr = f + " to " + t
				}
			} else if isValidDate(od.ODDate) {
				// Case 1: Full Day
				dateStr = formatDate(od.ODDate.String)
			}
		} else {
			// External
			if isValidDate(od.FromDate) && isValidDate(od.ToDate) {
				f := formatDate(od.FromDate.String)
				t := formatDate(od.ToDate.String)
				if f == t {
					dateStr = f
				} else {
					dateStr = f + " to " + t
				}
			} else if isValidDate(od.ODDate) {
				// Fallback for single day external if applicable
				dateStr = formatDate(od.ODDate.String)
			}
		}

		// Display Text Mapping
		displayText := displayStatus
		switch displayStatus {
		case "Mentor Accepted":
			displayText = "Mentor Approved"
		case "HOD Accepted":
			displayText = "Approved"
		case "Mentor Rejected":
			displayText = "Rejected (Mentor)"
		case "HOD Rejected":
			displayText = "Rejected (HOD)"
		case "mentor rejected":
			displayText = "Rejected (Mentor)"
        case "Principal Accepted":
            displayText = "Principal Approved"
        case "Principal Rejected":
            displayText = "Rejected (Principal)"
		}

		// JSON for team
		var teamJSON string
		if len(teamMembers) > 0 {
			teamJSONBytes, _ := json.Marshal(teamMembers)
			teamJSON = string(teamJSONBytes)
		} else {
			teamJSON = "[]"
		}

		dashboardODs = append(dashboardODs, DashboardOD{
			ODApplication:   od,
			DateStr:         dateStr,
			DisplayStatus:   displayStatus,
			DisplayText:     displayText,
			BadgeClass:      badgeClass,
			TeamMembers:     teamMembers,
			IsIndividual:    len(teamMembers) == 0,
			TeamMembersJSON: teamJSON,
		})
	}

	// Prepare data
	flashSuccess := "" // Flash messages logic to be implemented
	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	data := StudentDashboardData{
		User:         userMap,
		ODs:          dashboardODs,
		FlashSuccess: flashSuccess,
	}

	RenderTemplate(w, "templates/student_dashboard.html", data)
}
