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

type PrincipalDashboardData struct {
	User         map[string]interface{}
	Applications []PrincipalDashboardOD
	Search       string
	FlashSuccess string
}

type PrincipalDashboardOD struct {
	models.ODApplication
	DateStr     string
	BadgeClass  string
	TeamMembers []models.ODTeamMember
	TeamJSON    string
}

// PrincipalDashboard handler
func PrincipalDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "principal" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	// Logic: Fetch ODs approved by HOD and (od_type='external')
	// PHP: WHERE (status = 'HOD Accepted' OR status = 'Principal Accepted' OR status = 'Principal Rejected') AND od_type = 'external'

	// Fix: Added alias 'o' to od_applications
	query := "SELECT " + ODColumns + " FROM od_applications o WHERE (o.status = 'HOD Accepted' OR o.status = 'Principal Accepted' OR o.status = 'Principal Rejected') AND o.od_type = 'external' ORDER BY o.id DESC"

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Principal Query Error:", err)
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	// Helper functions (inline)
	formatDate := func(d string) string {
		d = strings.TrimSpace(d)
		if len(d) > 10 {
			d = d[:10]
		}
		t, err := time.Parse("2006-01-02", d)
		if err != nil { return d }
		if t.Year() <= 1 { return "-" }
		return t.Format("2 Jan 2006")
	}
	formatTime := func(tStr string) string {
		t, err := time.Parse("15:04:05", tStr)
		if err != nil {
			t, err = time.Parse("15:04", tStr)
			if err != nil { return tStr }
		}
		return t.Format("3:04 pm")
	}
	isValidDate := func(ns sql.NullString) bool {
		return ns.Valid && ns.String != "0000-00-00" && len(ns.String) > 0
	}

	var apps []PrincipalDashboardOD
	for rows.Next() {
		var od models.ODApplication
		err := rows.Scan(
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			log.Println("Principal Scan Error:", err)
			continue
		}

		// Fetch Team
		tmQuery := "SELECT id, od_id, member_name, member_regno, member_department, member_year, member_section, mentor, mentor_status FROM od_team_members WHERE od_id = ?"
		tmRows, err := database.DB.Query(tmQuery, od.ID)
		var team []models.ODTeamMember
		if err == nil {
			for tmRows.Next() {
				var m models.ODTeamMember
				tmRows.Scan(&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberDepartment, &m.MemberYear, &m.MemberSection, &m.Mentor, &m.MentorStatus)
				team = append(team, m)
			}
			tmRows.Close()
		}

		// Normalize ODType
		od.ODType = strings.ToLower(od.ODType)

		// Date Formatting Logic
		// Date Formatting Logic
		dateStr := "-"
		
		// Internal OD
		if od.ODType == "internal" {
			hasTime := od.FromTime.Valid && od.ToTime.Valid && od.FromTime.String != "00:00:00" && od.ToTime.String != "00:00:00"
			
			// Case 3: More than a day (FromDate & ToDate valid)
			if isValidDate(od.FromDate) && isValidDate(od.ToDate) {
				f := formatDate(od.FromDate.String)
				t := formatDate(od.ToDate.String)
				if f != "-" && t != "-" {
					if f == t {
						// Same day
						if hasTime {
							// Case 2: Period-wise on same day (fallback if data is stored this way)
							dateStr = formatTime(od.FromTime.String) + " to " + formatTime(od.ToTime.String) + ", " + f
						} else {
							// Case 1: Full Day
							dateStr = f
						}
					} else {
						// Case 3: Range
						dateStr = f + " to " + t
					}
				}
			} else if isValidDate(od.ODDate) {
				// We have a single OD Date
				d := formatDate(od.ODDate.String)
				if d != "-" {
					if hasTime {
						// Case 2: Period-wise
						dateStr = formatTime(od.FromTime.String) + " to " + formatTime(od.ToTime.String) + ", " + d
					} else {
						// Case 1: Full Day
						dateStr = d
					}
				}
			}
		} else {
			// External OD
			if isValidDate(od.FromDate) && isValidDate(od.ToDate) {
				f := formatDate(od.FromDate.String)
				t := formatDate(od.ToDate.String)
				if f != "-" && t != "-" {
					if f == t {
						dateStr = f
					} else {
						dateStr = f + " to " + t
					}
				}
			} else if isValidDate(od.ODDate) {
				dateStr = formatDate(od.ODDate.String)
			}
		}

		badgeClass := "bg-secondary"
		if od.Status == "Principal Accepted" {
			badgeClass = "bg-success"
		} else if od.Status == "Principal Rejected" {
			badgeClass = "bg-danger"
		} else if od.Status == "HOD Accepted" {
			badgeClass = "bg-warning text-dark"
		}

		teamJSON, _ := json.Marshal(team)
		if len(team) == 0 {
			teamJSON = []byte("[]")
		}

		apps = append(apps, PrincipalDashboardOD{
			ODApplication: od,
			DateStr:       dateStr,
			BadgeClass:    badgeClass,
			TeamMembers:   team,
			TeamJSON:      string(teamJSON),
		})
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	data := PrincipalDashboardData{
		User:         userMap,
		Applications: apps,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/principal_dashboard.html", data)
}

// PrincipalAction handler
func PrincipalAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", 400)
		return
	}

	action := r.FormValue("action")
	odID := r.FormValue("od_id")

	if action == "" || odID == "" {
		http.Error(w, "Missing fields", 400)
		return
	}

	status := "Principal Rejected"
	if action == "accept" {
		status = "Principal Accepted"
	}

	_, err := database.DB.Exec("UPDATE od_applications SET status = ? WHERE id = ?", status, odID)
	if err != nil {
		log.Println("Principal Update Error:", err)
		http.Error(w, "Database error", 500)
		return
	}

	// Redirect back to dashboard
	http.Redirect(w, r, "/principal/dashboard", http.StatusSeeOther)
}
