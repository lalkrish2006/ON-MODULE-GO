package handlers

import (
	"database/sql"
	"encoding/json"
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

type CADashboardData struct {
	User         map[string]interface{}
	Applications []CADashboardOD
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
	name := r.URL.Query().Get("name")
	regNo := r.URL.Query().Get("reg_no")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	odType := r.URL.Query().Get("od_type")
	class := r.URL.Query().Get("class")
	yearFilter := r.URL.Query().Get("year")

	// Get CA's Class Details from Session
	caYear, _ := session.Values["year"].(int)
	caDept, _ := session.Values["department"].(string)
	caSection, _ := session.Values["section"].(string)

	// Logic: Read-only view of finalized ODs
	// PHP ca_dashboard.php: 
	// WHERE (o.od_type = 'internal' AND o.status = 'HOD Accepted') OR (o.od_type = 'external' AND o.status = 'Principal Accepted')
	
	// Fixed: Added alias 'o' to table to match ODColumns (o.id, etc.)
	// Added: Filter by Class (Year, Dept, Section)
	query := `SELECT t.member_name, t.member_regno, t.member_year, t.member_section, ` + ODColumns + ` FROM od_team_members t 
		JOIN od_applications o ON o.id = t.od_id
		WHERE ((o.od_type = 'internal' AND o.status = 'HOD Accepted') 
		   OR (o.od_type = 'external' AND o.status = 'Principal Accepted'))
		AND t.member_year = ? AND t.member_department = ? AND t.member_section = ?`
	
	var args []interface{}
	args = append(args, caYear, caDept, caSection)

	if search != "" {
		like := "%" + search + "%"
		query += ` AND (
            o.id LIKE ? OR t.member_regno LIKE ? OR t.member_name LIKE ? OR
            t.member_year LIKE ? OR t.member_department LIKE ? OR t.member_section LIKE ? OR
            o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
            o.event_name LIKE ?
        )`
		for i := 0; i < 10; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		query += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		query += " AND t.member_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if regNo != "" {
		query += " AND t.member_regno LIKE ?"
		args = append(args, "%"+regNo+"%")
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
		query += " AND t.member_section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		query += " AND t.member_year = ?"
		args = append(args, yearFilter)
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

		// Use member specific data for display
		od.StudentName = mName
		od.RegisterNo = mReg
		od.Year = strconv.Itoa(mYear)
		od.Section = mSection

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
		Name:         name,
		RegNo:        regNo,
		StartDate:    startDate,
		EndDate:      endDate,
		ODType:       odType,
		Class:        class,
		YearFilter:   yearFilter,
		FlashSuccess: "",
	}

	RenderTemplate(w, "templates/ca_dashboard.html", data)
}

// DownloadCAHistoryPDF handler
func DownloadCAHistoryPDF(w http.ResponseWriter, r *http.Request) {
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

	caYear, _ := session.Values["year"].(int)
	caDept, _ := session.Values["department"].(string)
	caSection, _ := session.Values["section"].(string)

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
		   OR (o.od_type = 'external' AND o.status = 'Principal Accepted'))
		AND t.member_year = ? AND t.member_department = ? AND t.member_section = ?`
	
	args := []interface{}{caYear, caDept, caSection}

	if search != "" {
		like := "%" + search + "%"
		query += ` AND (
            o.id LIKE ? OR t.member_regno LIKE ? OR t.member_name LIKE ? OR
            t.member_year LIKE ? OR t.member_department LIKE ? OR t.member_section LIKE ? OR
            o.od_type LIKE ? OR o.purpose LIKE ? OR o.college_name LIKE ? OR
            o.event_name LIKE ?
        )`
		for i := 0; i < 10; i++ {
			args = append(args, like)
		}
	}

	if month != "" {
		query += " AND (DATE_FORMAT(o.from_date, '%Y-%m') = ? OR DATE_FORMAT(o.od_date, '%Y-%m') = ?)"
		args = append(args, month, month)
	}

	if name != "" {
		query += " AND t.member_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if regNo != "" {
		query += " AND t.member_regno LIKE ?"
		args = append(args, "%"+regNo+"%")
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
		query += " AND t.member_section = ?"
		args = append(args, class)
	}

	if yearFilter != "" {
		query += " AND t.member_year = ?"
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
	pdf.Cell(280, 10, "OD History Report - " + caDept + " (" + strconv.Itoa(caYear) + "-" + caSection + ")")
	pdf.Ln(12)

	headers := []string{"ID", "Name", "Reg No", "Year", "Type", "Dates", "Purpose", "Status"}
	widths := []float64{15, 45, 30, 15, 20, 50, 70, 35}

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

		if pdf.GetY()+rowHeight > 275 {
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
	w.Header().Set("Content-Disposition", "attachment; filename=ca_od_history.pdf")
	pdf.Output(w)
}
