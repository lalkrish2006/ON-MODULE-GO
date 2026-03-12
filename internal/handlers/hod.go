package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strconv"
	"strings"
	"time"
)

type HODDashboardData struct {
	User        map[string]interface{}
	PendingODs  []HODDashboardRow
	HistoryODs  []HODDashboardRow
	Search      string
	MonthFilter string
	Name        string
	RegNo       string
	StartDate   string
	EndDate     string
	ODType      string
	Class       string
	YearFilter  string
	FlashSuccess string
	IsAdmin     bool
}

// ... (HODDashboardRow and HODTeamMember structs remain unchanged)

// HODDashboard handler


type HODDashboardRow struct {
	ID          int
	RegisterNo  string
	StudentName string
	Year        string
	Department  string
	Section     string
	ODType      string
	Purpose     string
	CollegeName string
	EventName   string
	DateStr     string
	Status      string
	BadgeClass  string
	TeamMembers []HODTeamMember
	TeamJSON    string // For Modal
}

type HODTeamMember struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	RegNo        string `json:"reg_no"`
	Year         string `json:"year"`
	Dept         string `json:"dept"`
	Section      string `json:"section"`
	Mentor       string `json:"mentor"`
	MentorStatus string `json:"mentor_status"`
}

// HODDashboard handler
func HODDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "hod" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	dept, _ := session.Values["department"].(string)
	if role == "admin" {
		if d := r.URL.Query().Get("department"); d != "" {
			dept = d
		}
	}

	search := r.URL.Query().Get("search")
	month := r.URL.Query().Get("month")
	name := r.URL.Query().Get("name")
	regNo := r.URL.Query().Get("reg_no")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	odType := r.URL.Query().Get("od_type")
	class := r.URL.Query().Get("class")
	yearFilter := r.URL.Query().Get("year")

	// PHP Logic: Single Query for all items
	baseQuery := `
		SELECT DISTINCT ` + ODColumns + ` 
		FROM od_applications o 
		LEFT JOIN od_team_members t ON o.id = t.od_id
		WHERE o.department = ?
		AND NOT (
			(o.status = 'Mentors Rejected') OR 
			EXISTS (SELECT 1 FROM od_team_members t2 WHERE t2.od_id = o.id AND t2.mentor_status = 'Rejected' AND (SELECT COUNT(*) FROM od_team_members t3 WHERE t3.od_id = o.id) = 1)
		)
		AND (
			(o.od_type != 'Internal' AND (o.status = 'Mentors Accepted' OR o.status LIKE 'HOD%'))
			OR
			(o.od_type = 'Internal' AND NOT EXISTS (SELECT 1 FROM od_team_members tm WHERE tm.od_id = o.id AND tm.mentor_status = 'Pending') AND (o.status LIKE 'Mentors%' OR o.status LIKE 'HOD%'))
		)
	`
	args := []interface{}{dept}

	if search != "" {
		like := "%" + search + "%"
		baseQuery += ` AND (
			o.id LIKE ? OR o.register_no LIKE ? OR o.student_name LIKE ? OR
			o.year LIKE ? OR o.department LIKE ? OR o.section LIKE ? OR
			o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
			o.event_name LIKE ? OR t.member_name LIKE ? OR t.member_regno LIKE ? OR t.mentor LIKE ?
		)`
		for i := 0; i < 13; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		baseQuery += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		baseQuery += " AND o.student_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if regNo != "" {
		baseQuery += " AND o.register_no LIKE ?"
		args = append(args, "%"+regNo+"%")
	}

	if startDate != "" {
		baseQuery += " AND (o.from_date >= ? OR o.od_date >= ?)"
		args = append(args, startDate, startDate)
	}

	if endDate != "" {
		baseQuery += " AND (o.to_date <= ? OR o.od_date <= ?)"
		args = append(args, endDate, endDate)
	}

	if odType != "" {
		baseQuery += " AND o.od_type = ?"
		args = append(args, odType)
	}

	if class != "" {
		baseQuery += " AND o.section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		baseQuery += " AND o.year = ?"
		args = append(args, yearFilter)
	}

	baseQuery += " ORDER BY o.id DESC"

	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		log.Println("HOD Query Error:", err)
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	var pendingODs []HODDashboardRow
	var historyODs []HODDashboardRow

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
		tmQuery := "SELECT id, member_name, member_regno, member_year, member_department, member_section, mentor, mentor_status FROM od_team_members WHERE od_id = ?"
		tmRows, err := database.DB.Query(tmQuery, od.ID)
		var team []HODTeamMember
		if err == nil {
			for tmRows.Next() {
				var m HODTeamMember
				var ms sql.NullString
				var mentor sql.NullString
				tmRows.Scan(&m.ID, &m.Name, &m.RegNo, &m.Year, &m.Dept, &m.Section, &mentor, &ms)
				m.MentorStatus = ms.String
				m.Mentor = mentor.String
				team = append(team, m)
			}
			tmRows.Close()
		}

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

		// Normalize ODType
		od.ODType = strings.ToLower(od.ODType)

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
		if strings.Contains(strings.ToLower(od.Status), "accepted") {
			badgeClass = "bg-success"
		} else if strings.Contains(strings.ToLower(od.Status), "rejected") {
			badgeClass = "bg-danger"
		} else if od.Status == "Mentor Accepted" {
			badgeClass = "bg-warning text-dark"
		}

		teamJSON, _ := json.Marshal(team)

		row := HODDashboardRow{
			ID:          od.ID,
			RegisterNo:  od.RegisterNo,
			StudentName: od.StudentName,
			Year:        od.Year,
			Department:  od.Department,
			Section:     od.Section,
			ODType:      od.ODType,
			Purpose:     od.Purpose,
			CollegeName: od.CollegeName.String,
			EventName:   od.EventName.String,
			DateStr:     dateStr,
			Status:      od.Status,
			BadgeClass:  badgeClass,
			TeamMembers: team,
			TeamJSON:    string(teamJSON),
		}

		// Split into Pending vs History
		if od.Status == "Mentors Accepted" || od.Status == "Mentors Reviewed" {
			pendingODs = append(pendingODs, row)
		} else {
			historyODs = append(historyODs, row)
		}
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}
	// Add dept to user map if admin
	if role == "admin" {
		userMap["department"] = dept
	}

	// Extract Flash Message
	flashSuccess, _ := session.Values["flash_success"].(string)
	if flashSuccess != "" {
		delete(session.Values, "flash_success")
		session.Save(r, w)
	}

	data := HODDashboardData{
		User:         userMap,
		PendingODs:   pendingODs,
		HistoryODs:   historyODs,
		Search:       search,
		MonthFilter:  month,
		Name:         name,
		RegNo:        regNo,
		StartDate:    startDate,
		EndDate:      endDate,
		ODType:       odType,
		Class:        class,
		YearFilter:   yearFilter,
		FlashSuccess: flashSuccess,
		IsAdmin:      role == "admin",
	}

	RenderTemplate(w, "templates/hod_dashboard.html", data)
}

// HODAction handler
// HODAction handler
func HODAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", 400)
		return
	}

	action := r.FormValue("action")
	odID := r.FormValue("od_id")
	memberID := r.FormValue("member_id")

	// Case 1: Application-level Action (No member_id)
	if odID != "" && memberID == "" && action != "" {
		status := "HOD Rejected"
		if action == "accept" {
			status = "HOD Accepted"
		}

		_, err := database.DB.Exec("UPDATE od_applications SET status = ? WHERE id = ?", status, odID)
		if err != nil {
			log.Println("HOD Update Error:", err)
			http.Error(w, "Database error", 500)
			return
		}

		// Lab Tech Notification Logic
		if action == "accept" {
			go func(id string) {
				var odType string
				var labRequired int
				var dept string
				
				// Fetch OD Details
				err := database.DB.QueryRow("SELECT od_type, lab_required, department FROM od_applications WHERE id = ?", id).Scan(&odType, &labRequired, &dept)
				if err == nil && strings.ToLower(odType) == "internal" && labRequired == 1 {
					// Fetch Tech Email
					var techEmail, techName string
					err = database.DB.QueryRow("SELECT email, name FROM lab_technicians WHERE department = ? LIMIT 1", dept).Scan(&techEmail, &techName)
					if err == nil && techEmail != "" {
						// Assuming HOD email is in session or we fetch it? 
						// PHP fetches sender HOD email from DB based on session dept.
						// We can just send system notification for now or fetch HOD.
						// Let's use generic system sender for now to keep it simple as we lack full HOD profiles in session easily.
						// Or better, just send the email.
						subject := "Lab Booking Required for Internal OD (ID: " + id + ")"
						body := "Dear " + techName + ",\n\nThe HOD has approved an internal OD application (ID: " + id + ") that requires lab facilities.\n\nPlease check the OD Module."
						services.SendEmail(techEmail, subject, body)
					}
				}
			}(odID)
		}
	}

	// Case 2: Individual Team Member Override
	if memberID != "" && odID != "" && action != "" {
		status := "HOD Rejected"
		if action == "accept" {
			status = "HOD Accepted"
		}

		_, err := database.DB.Exec("UPDATE od_team_members SET mentor_status = ? WHERE id = ?", status, memberID)
		if err != nil {
			log.Println("HOD Member Update Error:", err)
			http.Error(w, "Database error", 500)
			return
		}
	}

	// Redirect back to dashboard
	http.Redirect(w, r, "/hod/dashboard", http.StatusSeeOther)
}

// DownloadHODHistoryPDF handler
func DownloadHODHistoryPDF(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "hod" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	dept, _ := session.Values["department"].(string)
	if role == "admin" {
		if d := r.URL.Query().Get("department"); d != "" {
			dept = d
		}
	}

	search := r.URL.Query().Get("search")
	month := r.URL.Query().Get("month")
	name := r.URL.Query().Get("name")
	regNo := r.URL.Query().Get("reg_no")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	odType := r.URL.Query().Get("od_type")
	class := r.URL.Query().Get("class")
	yearFilter := r.URL.Query().Get("year")

	// Same query logic as HODDashboard but focusing on history items primarily or all?
	// The user asked for OD History, usually meaning all data in the board.
	baseQuery := `
		SELECT DISTINCT ` + ODColumns + ` 
		FROM od_applications o 
		LEFT JOIN od_team_members t ON o.id = t.od_id
		WHERE o.department = ?
	`
	args := []interface{}{dept}

	if search != "" {
		like := "%" + search + "%"
		baseQuery += ` AND (
			o.id LIKE ? OR o.register_no LIKE ? OR o.student_name LIKE ? OR
			o.year LIKE ? OR o.department LIKE ? OR o.section LIKE ? OR
			o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
			o.event_name LIKE ? OR t.member_name LIKE ? OR t.member_regno LIKE ? OR t.mentor LIKE ?
		)`
		for i := 0; i < 13; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		baseQuery += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		baseQuery += " AND o.student_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if regNo != "" {
		baseQuery += " AND o.register_no LIKE ?"
		args = append(args, "%"+regNo+"%")
	}

	if startDate != "" {
		baseQuery += " AND (o.from_date >= ? OR o.od_date >= ?)"
		args = append(args, startDate, startDate)
	}

	if endDate != "" {
		baseQuery += " AND (o.to_date <= ? OR o.od_date <= ?)"
		args = append(args, endDate, endDate)
	}

	if odType != "" {
		baseQuery += " AND o.od_type = ?"
		args = append(args, odType)
	}

	if class != "" {
		baseQuery += " AND o.section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		baseQuery += " AND o.year = ?"
		args = append(args, yearFilter)
	}

	baseQuery += " ORDER BY o.id DESC"

	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		log.Println("PDF Query Error:", err)
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(280, 10, "OD History Report - " + dept)
	pdf.Ln(12)

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)
	headers := []string{"ID", "Name", "Reg No", "Year", "Type", "Dates", "Purpose", "Status"}
	widths := []float64{15, 45, 30, 15, 20, 50, 70, 35}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table Data
	pdf.SetFont("Arial", "", 9)
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

		// Simplified Date logic for PDF
		dateStr := "-"
		formatDate := func(ns sql.NullString) string {
			if !ns.Valid || ns.String == "0000-00-00" { return "-" }
			t, _ := time.Parse("2006-01-02", ns.String[:10])
			return t.Format("02-01-06")
		}

		if od.ODType == "Internal" || od.ODType == "internal" {
			if od.ODDate.Valid && od.ODDate.String != "0000-00-00" {
				dateStr = formatDate(od.ODDate)
			} else {
				dateStr = formatDate(od.FromDate) + " to " + formatDate(od.ToDate)
			}
		} else {
			dateStr = formatDate(od.FromDate) + " to " + formatDate(od.ToDate)
		}

		pdf.CellFormat(widths[0], 10, strconv.Itoa(od.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 10, od.StudentName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 10, od.RegisterNo, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[3], 10, od.Year, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 10, od.ODType, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 10, dateStr, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 10, od.Purpose, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[7], 10, od.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=od_history_"+dept+".pdf")
	pdf.Output(w)
}
