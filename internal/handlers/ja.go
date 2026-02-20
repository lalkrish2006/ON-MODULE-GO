package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strings"
	"time"
)

type JADashboardData struct {
	User         map[string]interface{}
	Applications []JADashboardOD
	Search       string
	MonthFilter  string
	FlashSuccess string
}

type JADashboardOD struct {
	models.ODApplication
	DateStr       string
	BadgeClass    string
	DisplayStatus string
	TeamMembers   []models.ODTeamMember
	TeamJSON      string
	IsIndividual  bool
	BonafideBadget string
}

// Display struct to handle NullString for JSON
type JATeamMemberDisplay struct {
	MemberName     string `json:"member_name"`
	MemberRegNo    string `json:"member_regno"`
	Mentor         string `json:"mentor"`
	MentorStatus   string `json:"mentor_status"`
	MemberDept     string `json:"member_department"`
	MemberYear     string `json:"member_year"`
	MemberSection  string `json:"member_section"`
}

func JADashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "ja" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	search := r.URL.Query().Get("search")
	month := r.URL.Query().Get("month")
	export := r.URL.Query().Get("export")

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

	// Query Construction
	// Using ODColumns from helpers.go which contains "o.id, o.register_no..."
	// Added 'AS o' for explicit aliasing
	query := `SELECT DISTINCT ` + ODColumns + ` FROM od_applications AS o 
		LEFT JOIN od_team_members t ON o.id = t.od_id
		WHERE ((o.od_type = 'internal' AND o.status = 'HOD Accepted') 
		   OR (o.od_type = 'external' AND o.status = 'Principal Accepted'))`
	
	var args []interface{}

	if search != "" {
		like := "%" + search + "%"
		query += ` AND (
            o.id LIKE ? OR o.register_no LIKE ? OR o.student_name LIKE ? OR
            o.year LIKE ? OR o.department LIKE ? OR o.section LIKE ? OR
            o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
            o.event_name LIKE ? OR t.member_name LIKE ? OR t.member_regno LIKE ?
        )`
		for i := 0; i < 12; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		query += " AND DATE_FORMAT(o.from_date, '%Y-%m') = ?"
		args = append(args, month)
	}

	query += " ORDER BY o.id DESC"

	// Log query for debugging
	log.Println("JA Query:", query)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Println("JA Query Error:", err)
		http.Error(w, "Database error: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var apps []JADashboardOD

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
		var displayTeam []JATeamMemberDisplay

		if err == nil {
			for tmRows.Next() {
				var m models.ODTeamMember
				tmRows.Scan(&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberDepartment, &m.MemberYear, &m.MemberSection, &m.Mentor, &m.MentorStatus)
				team = append(team, m)

				ms := "Pending"
				if m.MentorStatus.Valid {
					ms = m.MentorStatus.String
				}
				displayTeam = append(displayTeam, JATeamMemberDisplay{
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

		// Normalize ODType make lowercase for logic
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
		displayStatus := "Principal Accepted"
		if od.ODType == "internal" {
			displayStatus = "HOD Accepted"
			badgeClass = "bg-success"
		} else {
			badgeClass = "bg-primary"
		}

		bonafideBadge := "No"
		if od.RequestBonafide == 1 {
			bonafideBadge = "Required"
		}

		teamJSONBytes, _ := json.Marshal(displayTeam)
		if len(displayTeam) == 0 {
			teamJSONBytes = []byte("[]")
		}

		apps = append(apps, JADashboardOD{
			ODApplication: od,
			DateStr:       dateStr,
			BadgeClass:    badgeClass, 
			DisplayStatus: displayStatus,
			TeamMembers:   team,
			TeamJSON:      string(teamJSONBytes),
			IsIndividual:  len(team) == 0,
			BonafideBadget: bonafideBadge,
		})
	}

	// CSV Export Logic
	if export == "csv" {
		filename := fmt.Sprintf("ja_od_export_%s.csv", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Headers
		writer.Write([]string{"ID", "Register No", "Student Name", "Year", "Department", "Section", "OD Type", "Purpose", "College Name", "Event Name", "Dates", "Status", "Bonafide Required", "Team Members"})

		for _, app := range apps {
			teamString := "Individual"
			if len(app.TeamMembers) > 0 {
				var members []string
				for _, m := range app.TeamMembers {
					members = append(members, fmt.Sprintf("%s [%s]", m.MemberName, m.MemberRegNo))
				}
				teamString = strings.Join(members, "; ")
			}

			// College/Event fallback
			cn := app.CollegeName.String
			if !app.CollegeName.Valid { cn = "N/A" }
			en := app.EventName.String
			if !app.EventName.Valid { en = "N/A" }

			writer.Write([]string{
				fmt.Sprintf("%d", app.ID),
				app.RegisterNo,
				app.StudentName,
				app.Year,
				app.Department,
				app.Section,
				strings.Title(app.ODType),
				app.Purpose,
				cn,
				en,
				app.DateStr,
				app.DisplayStatus,
				app.BonafideBadget,
				teamString,
			})
		}
		return
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	data := JADashboardData{
		User:         userMap,
		Applications: apps,
		Search:       search,
		MonthFilter:  month,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/ja_dashboard.html", data)
}
