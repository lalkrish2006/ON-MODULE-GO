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

type CADashboardData struct {
	User         map[string]interface{}
	Applications []CADashboardOD
	Search       string
	MonthFilter  string
	FlashSuccess string
}

type CADashboardOD struct {
	models.ODApplication
	DateStr       string
	BadgeClass    string
	DisplayStatus string // Added for nuanced display logic
	TeamMembers   []models.ODTeamMember
	TeamJSON      string
	IsIndividual  bool
}

// Display struct to handle NullString for JSON
type CATeamMemberDisplay struct {
	MemberName     string `json:"MemberName"`
	MemberRegNo    string `json:"MemberRegNo"`
	Mentor         string `json:"Mentor"`
	MentorStatus   string `json:"MentorStatus"`
	MemberDept     string `json:"MemberDepartment"`
	MemberYear     string `json:"MemberYear"`
	MemberSection  string `json:"MemberSection"`
}

// CADashboard handler
func CADashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "ca" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	search := r.URL.Query().Get("search")
	month := r.URL.Query().Get("month")

	// Get CA's Class Details from Session
	caYear, _ := session.Values["year"].(int)
	caDept, _ := session.Values["department"].(string)
	caSection, _ := session.Values["section"].(string)

	// Logic: Read-only view of finalized ODs
	// PHP ca_dashboard.php: 
	// WHERE (o.od_type = 'internal' AND o.status = 'HOD Accepted') OR (o.od_type = 'external' AND o.status = 'Principal Accepted')
	
	// Fixed: Added alias 'o' to table to match ODColumns (o.id, etc.)
	// Added: Filter by Class (Year, Dept, Section)
	query := `SELECT DISTINCT ` + ODColumns + ` FROM od_applications o 
		LEFT JOIN od_team_members t ON o.id = t.od_id
		WHERE ((o.od_type = 'internal' AND o.status = 'HOD Accepted') 
		   OR (o.od_type = 'external' AND o.status = 'Principal Accepted'))
		AND o.year = ? AND o.department = ? AND o.section = ?`
	
	var args []interface{}
	args = append(args, caYear, caDept, caSection)

	if search != "" {
		like := "%" + search + "%"
		// Search across multiple fields as per PHP
		query += ` AND (
            o.id LIKE ? OR o.register_no LIKE ? OR o.student_name LIKE ? OR
            o.year LIKE ? OR o.department LIKE ? OR o.section LIKE ? OR
            o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
            o.event_name LIKE ? OR t.member_name LIKE ? OR t.member_regno LIKE ?
        )`
		// Append args 12 times
		for i := 0; i < 12; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		query += " AND DATE_FORMAT(o.from_date, '%Y-%m') = ?"
		args = append(args, month)
	}

	query += " ORDER BY o.id DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Println("CA Query Error:", err)
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

	var apps []CADashboardOD
	
	for rows.Next() {
		var od models.ODApplication
		err := rows.Scan(
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Fetch Team
		tmQuery := "SELECT id, od_id, member_name, member_regno, member_department, member_year, member_section, mentor, mentor_status FROM od_team_members WHERE od_id = ?"
		tmRows, err := database.DB.Query(tmQuery, od.ID)
		var team []models.ODTeamMember
		var displayTeam []CATeamMemberDisplay

		if err == nil {
			for tmRows.Next() {
				var m models.ODTeamMember
				tmRows.Scan(&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberDepartment, &m.MemberYear, &m.MemberSection, &m.Mentor, &m.MentorStatus)
				team = append(team, m)
				
				ms := "Pending"
				if m.MentorStatus.Valid {
					ms = m.MentorStatus.String
				}
				displayTeam = append(displayTeam, CATeamMemberDisplay{
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
					// Case 3: Range
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
						dateStr = formatTime(od.FromTime.String) + " to " + formatTime(od.ToTime.String) + ", " + d
					} else {
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
		displayStatus := od.Status
		if od.Status == "Principal Accepted" {
			badgeClass = "bg-success"
		} else if od.Status == "HOD Accepted" {
			badgeClass = "bg-success" // CA sees finalized ones as green typically, matching PHP logic implicitly
		}

		teamJSONBytes, _ := json.Marshal(displayTeam)
		if len(displayTeam) == 0 {
			teamJSONBytes = []byte("[]")
		}

		apps = append(apps, CADashboardOD{
			ODApplication: od,
			DateStr:       dateStr,
			BadgeClass:    badgeClass, 
			DisplayStatus: displayStatus,
			TeamMembers:   team,
			TeamJSON:      string(teamJSONBytes),
			IsIndividual:  len(team) == 0,
		})
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	data := CADashboardData{
		User:         userMap,
		Applications: apps,
		Search:       search,
		MonthFilter:  month,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/ca_dashboard.html", data)
}
