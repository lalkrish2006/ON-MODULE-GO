package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strconv"
	"strings"
	"time"
)

type LabTechDashboardData struct {
	User         map[string]interface{}
	Applications []LabTechOD
	Labs         []string
	LabFilter    string
	FlashSuccess string
}

type LabTechOD struct {
	models.ODApplication
	DateStr     string
	BadgeClass  string
	TeamMembers []models.ODTeamMember
	TeamJSON    string
}

// Display struct to handle NullString for JSON
type TeamMemberDisplay struct {
	MemberName     string `json:"member_name"`
	MemberRegNo    string `json:"member_regno"`
	Mentor         string `json:"mentor"`
	MentorStatus   string `json:"mentor_status"`
	MemberDept     string `json:"member_department"`
	MemberYear     string `json:"member_year"`
	MemberSection  string `json:"member_section"`
}

// LabTechDashboard handler
func LabTechDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "labtech" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	labFilter := r.URL.Query().Get("lab_name")

	// Explicitly list columns to ensure aliases match Scan
	cols := `o.id, o.register_no, o.student_name, o.year, o.department, o.section, 
	o.od_type, o.purpose, o.college_name, o.event_name, o.from_date, o.to_date, 
	o.od_date, o.from_time, o.to_time, o.status, o.request_bonafide, 
	o.lab_required, o.lab_name, o.system_required, o.created_at`

	// Logic: status='HOD Accepted' AND lab_required=1
	query := "SELECT " + cols + " FROM od_applications o WHERE status='HOD Accepted' AND lab_required=1"
	var args []interface{}

	if labFilter != "" {
		query += " AND o.lab_name = ?"
		args = append(args, labFilter)
	}

	query += " ORDER BY o.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Println("LabTech Query Error:", err)
		http.Error(w, "Database error: "+err.Error(), 500)
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
		return t.Format("3:04 PM") // Upper case PM as per request implies standard
	}
	isValidDate := func(ns sql.NullString) bool {
		return ns.Valid && ns.String != "0000-00-00" && len(ns.String) > 0
	}

	var apps []LabTechOD
	for rows.Next() {
		var od models.ODApplication
		err := rows.Scan(
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			log.Println("LabTech Scan Error:", err)
			continue
		}

		// Fetch Team
		tmQuery := "SELECT id, od_id, member_name, member_regno, member_department, member_year, member_section, mentor, mentor_status FROM od_team_members WHERE od_id = ?"
		tmRows, err := database.DB.Query(tmQuery, od.ID)
		var team []models.ODTeamMember
		
		// Create display team for JSON
		var displayTeam []TeamMemberDisplay

		if err == nil {
			for tmRows.Next() {
				var m models.ODTeamMember
				tmRows.Scan(&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberDepartment, &m.MemberYear, &m.MemberSection, &m.Mentor, &m.MentorStatus)
				team = append(team, m)
				
				// Handle NullString for JSON
				ms := "Pending"
				if m.MentorStatus.Valid {
					ms = m.MentorStatus.String
				}
				displayTeam = append(displayTeam, TeamMemberDisplay{
					MemberName:    m.MemberName,
					MemberRegNo:   m.MemberRegNo,
					Mentor:        m.Mentor,
					MentorStatus:  ms,
					MemberDept:    m.MemberDepartment,
					MemberYear:    m.MemberYear,
					MemberSection: m.MemberSection,
				})
			}
			tmRows.Close()
		}

		// Normalize ODType
		od.ODType = strings.ToLower(od.ODType)

		// Date Formatting Logic
		dateStr := "-"
		
		// Internal OD
		if od.ODType == "internal" {
			hasTime := od.FromTime.Valid && od.ToTime.Valid && od.FromTime.String != "00:00:00" && od.ToTime.String != "00:00:00"
			
			if isValidDate(od.FromDate) && isValidDate(od.ToDate) {
				f := formatDate(od.FromDate.String)
				t := formatDate(od.ToDate.String)
				
				if f != t {
					// Case 3: FromDate and ToDate are different -> Range
					dateStr = f + " to " + t
				} else {
					// Same day
					if hasTime {
						// Case 2: Period-wise
						dateStr = formatTime(od.FromTime.String) + " to " + formatTime(od.ToTime.String) + ", " + f
					} else {
						// Case 1: Full Day
						dateStr = f
					}
				}
			} else if isValidDate(od.ODDate) {
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

		badgeClass := "bg-success" 
		
		// Marshall the display struct which has strings, not NullString objects
		teamJSONBytes, _ := json.Marshal(displayTeam)
		if len(displayTeam) == 0 {
			teamJSONBytes = []byte("[]")
		}

		apps = append(apps, LabTechOD{
			ODApplication: od,
			DateStr:       dateStr,
			BadgeClass:    badgeClass,
			TeamMembers:   team,
			TeamJSON:      string(teamJSONBytes),
		})
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	labs := []string{"IOS Lab", "CC4 Lab", "Cloud Lab", "Open Source Lab", "HPC Lab"}

	data := LabTechDashboardData{
		User:         userMap,
		Applications: apps,
		Labs:         labs,
		LabFilter:    labFilter,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/labtech_dashboard.html", data)
}

// LabTechAction handler (Change Lab)
func LabTechAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
    
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Parse error", 400)
        return
    }
    
    idStr := r.FormValue("update_lab_id")
    newLab := r.FormValue("new_lab")
    
    // Check if we have a lab name filter to preserve
    // Referer check or just a hidden input could work, but user didn't request complex preservation
    // I will try to preserve it if it's in the Referer, but safer to just redirect to dashboard
    
    if idStr != "" && newLab != "" {
        id, _ := strconv.Atoi(idStr)
        _, err := database.DB.Exec("UPDATE od_applications SET lab_name = ? WHERE id = ?", newLab, id)
        if err != nil {
            log.Println("Lab Update Error:", err)
        }
    }
    
    // Redirect back
    http.Redirect(w, r, "/labtech/dashboard", http.StatusSeeOther)
}
