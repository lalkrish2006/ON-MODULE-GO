package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strconv"
	"strings"
	"time"
)

type JADashboardData struct {
	User         map[string]interface{}
	Applications []JADashboardOD
	Search       string
	MonthFilter  string
	Name         string
	RegNo        string
	StartDate    string
	EndDate      string
	ODType       string
	Class        string
	YearFilter   string
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
	name := r.URL.Query().Get("name")
	regNo := r.URL.Query().Get("reg_no")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	odType := r.URL.Query().Get("od_type")
	class := r.URL.Query().Get("class")
	yearFilter := r.URL.Query().Get("year")
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
	query := `SELECT t.member_name, t.member_regno, t.member_year, t.member_section, ` + ODColumns + ` FROM od_team_members t 
		JOIN od_applications o ON o.id = t.od_id
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
		query += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		query += " AND (o.student_name LIKE ? OR t.member_name LIKE ?)"
		args = append(args, "%"+name+"%", "%"+name+"%")
	}

	if regNo != "" {
		query += " AND (o.register_no LIKE ? OR t.member_regno LIKE ?)"
		args = append(args, "%"+regNo+"%", "%"+regNo+"%")
	}

	if startDate != "" {
		query += " AND (o.from_date >= ? OR o.od_date >= ?)"
		args = append(args, startDate, startDate)
	}

	if endDate != "" {
		query += " AND (o.to_date <= ? OR o.od_date <= ?)"
		args = append(args, endDate, endDate)
	}

	if odType != "" {
		query += " AND LOWER(o.od_type) = LOWER(?)"
		args = append(args, odType)
	}

	if class != "" {
		query += " AND o.section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		query += " AND o.year = ?"
		args = append(args, yearFilter)
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
		var mName, mReg, mSection string
		var mYear int
		err := rows.Scan(
			&mName, &mReg, &mYear, &mSection,
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			log.Println("Scan Error:", err)
			continue
		}

		// Use member specific data
		od.StudentName = mName
		od.RegisterNo = mReg
		od.Year = strconv.Itoa(mYear)
		od.Section = mSection

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
		Name:         name,
		RegNo:        regNo,
		StartDate:    startDate,
		EndDate:      endDate,
		ODType:       odType,
		Class:        class,
		YearFilter:   yearFilter,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/ja_dashboard.html", data)
}

// DownloadJAHistoryPDF handler
func DownloadJAHistoryPDF(w http.ResponseWriter, r *http.Request) {
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
	name := r.URL.Query().Get("name")
	regNo := r.URL.Query().Get("reg_no")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	odType := r.URL.Query().Get("od_type")
	class := r.URL.Query().Get("class")
	yearFilter := r.URL.Query().Get("year")

	query := `SELECT t.member_name, t.member_regno, t.member_year, t.member_section, ` + ODColumns + ` FROM od_team_members t 
		JOIN od_applications o ON o.id = t.od_id
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
		query += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		query += " AND (o.student_name LIKE ? OR t.member_name LIKE ?)"
		args = append(args, "%"+name+"%", "%"+name+"%")
	}

	if regNo != "" {
		query += " AND (o.register_no LIKE ? OR t.member_regno LIKE ?)"
		args = append(args, "%"+regNo+"%", "%"+regNo+"%")
	}

	if startDate != "" {
		query += " AND (o.from_date >= ? OR o.od_date >= ?)"
		args = append(args, startDate, startDate)
	}

	if endDate != "" {
		query += " AND (o.to_date <= ? OR o.od_date <= ?)"
		args = append(args, endDate, endDate)
	}

	if odType != "" {
		query += " AND LOWER(o.od_type) = LOWER(?)"
		args = append(args, odType)
	}

	if class != "" {
		query += " AND o.section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		query += " AND o.year = ?"
		args = append(args, yearFilter)
	}

	query += " ORDER BY o.id DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(275, 10, "JA OD History Report")
	pdf.Ln(12)

	headers := []string{"ID", "Name", "Reg No", "Year", "Type", "Dates", "Purpose", "Status"}
	widths := []float64{15, 40, 25, 10, 20, 45, 90, 30}

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	for rows.Next() {
		var od models.ODApplication
		var mName, mReg, mSection string
		var mYear int
		rows.Scan(
			&mName, &mReg, &mYear, &mSection,
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)

		// Use member specific data
		od.StudentName = mName
		od.RegisterNo = mReg
		od.Year = strconv.Itoa(mYear)
		od.Section = mSection

		dateStr := "-"
		formatDate := func(ns sql.NullString) string {
			if !ns.Valid || ns.String == "0000-00-00" { return "-" }
			t, _ := time.Parse("2006-01-02", ns.String[:10])
			return t.Format("02-01-06")
		}

		if strings.ToLower(od.ODType) == "internal" {
			if od.ODDate.Valid && od.ODDate.String != "0000-00-00" {
				dateStr = formatDate(od.ODDate)
			} else {
				dateStr = formatDate(od.FromDate) + " to " + formatDate(od.ToDate)
			}
		} else {
			dateStr = formatDate(od.FromDate) + " to " + formatDate(od.ToDate)
		}

		purposeWidth := widths[6]
		lines := pdf.SplitLines([]byte(od.Purpose), purposeWidth)
		lineCount := len(lines)
		if lineCount == 0 { lineCount = 1 }
		cellHeight := 5.0
		rowHeight := float64(lineCount) * cellHeight
		if rowHeight < 10 { rowHeight = 10 }

		if pdf.GetY()+rowHeight > 180 {
			pdf.AddPage()
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(200, 200, 200)
			for i, h := range headers {
				pdf.CellFormat(widths[i], 10, h, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 9)
		}

		curX, curY := pdf.GetXY()
		pdf.CellFormat(widths[0], rowHeight, strconv.Itoa(od.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], rowHeight, od.StudentName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], rowHeight, od.RegisterNo, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[3], rowHeight, od.Year, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], rowHeight, od.ODType, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], rowHeight, dateStr, "1", 0, "C", false, 0, "")
		pdf.MultiCell(widths[6], cellHeight, od.Purpose, "1", "L", false)
		pdf.SetXY(curX+widths[0]+widths[1]+widths[2]+widths[3]+widths[4]+widths[5]+widths[6], curY)
		pdf.CellFormat(widths[7], rowHeight, od.Status, "1", 0, "C", false, 0, "")
		pdf.SetXY(curX, curY+rowHeight)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=ja_od_history.pdf")
	pdf.Output(w)
}
