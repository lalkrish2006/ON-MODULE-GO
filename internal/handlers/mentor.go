package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/models"
	"od-system/internal/services"
	"strconv"
	"strings"
	"time"
)

// MentorDashboardData holds data for the mentor dashboard
type MentorDashboardData struct {
	User        map[string]interface{}
	Rows        []MentorODRequest
	Search      string
	MonthFilter string
	ODTypeFilter string
	FlashSuccess string
}

// MentorODRequest struct for dashboard display
type MentorODRequest struct {
	models.ODApplication
	MemberID    int // ID from od_team_members table
	MemberName  string
	MemberRegNo string
	MemberYear  string
	MemberDept  string
	MemberSec   string
	MentorStatus string
	DateStr      string // Formatted date string
	TeamJSON     string
}

// MentorDashboard handler
func MentorDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "mentor" && role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	mentorName, _ := session.Values["name"].(string)

	// filters
	search := r.URL.Query().Get("search")
	odTypeFilter := r.URL.Query().Get("od_type")
	month := r.URL.Query().Get("month")

	// PHP Logic: Single Query for all items
	// SELECT t.*, o.* FROM od_team_members t LEFT JOIN od_applications o ON o.id = t.od_id WHERE t.mentor = ?
	query := `
		SELECT 
			t.id, t.od_id, t.member_name, t.member_regno, t.member_year, t.member_department, t.member_section, t.mentor_status,
			o.id, o.od_type, o.purpose, o.college_name, o.event_name, o.from_date, o.to_date, o.od_date, o.from_time, o.to_time, o.status
		FROM od_team_members t
		LEFT JOIN od_applications o ON o.id = t.od_id
		WHERE t.mentor = ?`
	
	args := []interface{}{mentorName}

	if search != "" {
		like := "%" + search + "%"
		query += " AND (t.member_name LIKE ? OR t.member_regno LIKE ?)"
		args = append(args, like, like)
	}
	if odTypeFilter != "" {
		query += " AND o.od_type = ?"
		args = append(args, odTypeFilter)
	}
	if month != "" {
		query += " AND DATE_FORMAT(o.from_date, '%Y-%m') = ?"
		args = append(args, month)
	}

	query += " ORDER BY t.od_id DESC, t.id ASC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Println("Mentor Query Error:", err)
		http.Error(w, "Database Error", 500)
		return
	}
	defer rows.Close()

	var dashboardRows []MentorODRequest
	// We need to fetch ALL team members for the modal JSON?
	// PHP builds `teamByOD` array.
	// We can do a second query for all team members involved in these ODs, OR just fetch all team members for ODs present.
	// Optimization: Fetch all team members for visible OD IDs?
	// For simplicity and parity, let's fetch all relevant team members.
	// Or we can lazy load? PHP loads all upfront: `SELECT * FROM od_team_members ORDER BY od_id ASC` (Heavy?)
	// Actually PHP: `SELECT * FROM od_team_members ORDER BY od_id ASC, id ASC`. It loads ENTIRE table?!
	// That's inefficient but parity. Let's do better: Fetch team members only for the ODs we are showing.
	// But first let's get the main rows.

	// Helper functions for formatting (inline) - Duplicated from dashboard.go for isolation
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

	for rows.Next() {
		var m models.ODTeamMember
		var o models.ODApplication
		var fromTime, toTime interface{} // nullable

		err := rows.Scan(
			&m.ID, &m.ODID, &m.MemberName, &m.MemberRegNo, &m.MemberYear, &m.MemberDepartment, &m.MemberSection, &m.MentorStatus,
			&o.ID, &o.ODType, &o.Purpose, &o.CollegeName, &o.EventName, &o.FromDate, &o.ToDate, &o.ODDate, &fromTime, &toTime, &o.Status,
		)
		if err != nil {
			continue
		}
		
		if fStr, ok := fromTime.([]byte); ok { o.FromTime = sql.NullString{String: string(fStr), Valid: true} }
		if tStr, ok := toTime.([]byte); ok { o.ToTime = sql.NullString{String: string(tStr), Valid: true} }

		// Normalize ODType
		o.ODType = strings.ToLower(o.ODType)

		// Date Formatting Logic
		dateStr := "-"
		if o.ODType == "internal" {
			hasTime := o.FromTime.Valid && o.ToTime.Valid && o.FromTime.String != "00:00:00" && o.ToTime.String != "00:00:00"
			
			if hasTime && isValidDate(o.ODDate) {
				dateStr = formatTime(o.FromTime.String) + " to " + formatTime(o.ToTime.String) + ", " + formatDate(o.ODDate.String)
			} else if isValidDate(o.FromDate) && isValidDate(o.ToDate) {
				f := formatDate(o.FromDate.String)
				t := formatDate(o.ToDate.String)
				if f == t {
					dateStr = f
				} else {
					dateStr = f + " to " + t
				}
			} else if isValidDate(o.ODDate) {
				dateStr = formatDate(o.ODDate.String)
			}
		} else {
			if isValidDate(o.FromDate) && isValidDate(o.ToDate) {
				f := formatDate(o.FromDate.String)
				t := formatDate(o.ToDate.String)
				if f == t {
					dateStr = f
				} else {
					dateStr = f + " to " + t
				}
			} else if isValidDate(o.ODDate) {
				dateStr = formatDate(o.ODDate.String)
			}
		}

		dashboardRows = append(dashboardRows, MentorODRequest{
			ODApplication: o,
			MemberID:      m.ID,
			MemberName:    m.MemberName,
			MemberRegNo:   m.MemberRegNo,
			MemberYear:    m.MemberYear,
			MemberDept:    m.MemberDepartment,
			MemberSec:     m.MemberSection,
			MentorStatus:  m.MentorStatus.String,
			DateStr:       dateStr,
		})
	}

	// Fetch Team Members for Modal (Grouped by ODID)
	// We can select all team members for the OD IDs we found.
	var odIDs []int
	seen := make(map[int]bool)
	for _, r := range dashboardRows {
		if !seen[r.ID] {
			odIDs = append(odIDs, r.ID)
			seen[r.ID] = true
		}
	}

	teamMap := make(map[int][]map[string]interface{})
	if len(odIDs) > 0 {
		// Construct query IN clause
		q := "SELECT id, od_id, member_name, member_regno, mentor, mentor_status FROM od_team_members WHERE od_id IN ("
		qArgs := []interface{}{}
		for i, id := range odIDs {
			if i > 0 { q += "," }
			q += "?"
			qArgs = append(qArgs, id)
		}
		q += ") ORDER BY id ASC"

		tRows, err := database.DB.Query(q, qArgs...)
		if err == nil {
			defer tRows.Close()
			for tRows.Next() {
				var tid, todid int
				var tname, treg, tmentor string
				var tstatus sql.NullString
				tRows.Scan(&tid, &todid, &tname, &treg, &tmentor, &tstatus)
				
				memberData := map[string]interface{}{
					"id": tid,
					"od_id": todid,
					"member_name": tname,
					"member_regno": treg,
					"mentor": tmentor, // PHP field name
					"mentor_status": tstatus.String,
				}
				teamMap[todid] = append(teamMap[todid], memberData)
			}
		}
	}

	// Inject Team JSON into dashboardRows?
	// The template iterates dashboardRows. We can add a field `TeamJSON` to `MentorODRequest`.
	// Let's modify the struct or map. 
	// The `MentorODRequest` struct in `mentor.go` is:
	/*
	type MentorODRequest struct {
		models.ODApplication
		MemberID    int
		MemberName  string
		...
		MentorStatus string // Added locally to valid string
		TeamJSON     string // NEW
	}
	*/
	
	// Re-mapping to include JSON
	var finalRows []MentorODRequest
	for _, r := range dashboardRows {
		tm := teamMap[r.ID]
		if tm == nil { tm = []map[string]interface{}{} } // empty array
		jsonBytes, _ := json.Marshal(tm)
		r.TeamJSON = string(jsonBytes)
		finalRows = append(finalRows, r)
	}
	
	// Helper for user data
	userMap := map[string]interface{}{
		"name": mentorName,
		"role": role,
	}

	data := struct {
		User         map[string]interface{}
		Rows         []MentorODRequest
		Search       string
		ODTypeFilter string
		MonthFilter  string
		FlashSuccess string
	}{
		User:         userMap,
		Rows:         finalRows,
		Search:       search,
		ODTypeFilter: odTypeFilter,
		MonthFilter:  month,
		FlashSuccess: "", // Should take from session if exists
	}

	RenderTemplate(w, "templates/mentor_dashboard.html", data)
}

// MentorAction handler
func MentorAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	session := services.GetSession(r)
	mentorName, _ := session.Values["name"].(string)

	memberID, _ := strconv.Atoi(r.FormValue("member_id"))
	action := r.FormValue("action")

	status := "Rejected"
	if action == "accept" {
		status = "Accepted"
	}

	// 1. Update specific member status
	// We check strict mentor ownership
	_, err := database.DB.Exec("UPDATE od_team_members SET mentor_status=? WHERE id=? AND mentor=?", status, memberID, mentorName)
	if err != nil {
		log.Println("Update Action Error:", err)
		http.Error(w, "Database Error", 500)
		return
	}

	// 2. Fetch OD ID
	var odID int
	err = database.DB.QueryRow("SELECT od_id FROM od_team_members WHERE id=?", memberID).Scan(&odID)
	if err != nil {
		// Should not happen
		return 
	}

	// 3. Check Statuses of ALL Members for this OD
	rows, err := database.DB.Query("SELECT mentor_status FROM od_team_members WHERE od_id=?", odID)
	if err != nil {
		log.Println("Check Status Error:", err)
		return
	}
	defer rows.Close()

	anyRejected := false
	anyPending := false

	for rows.Next() {
		var ms sql.NullString
		rows.Scan(&ms)
		s := ms.String
		if s == "Rejected" {
			anyRejected = true
		}
		if s == "Pending" || s == "" {
			anyPending = true
		}
		if s != "Accepted" {
			// allAccepted = false
		}
	}

	finalStatus := ""
	if anyRejected {
		finalStatus = "Mentors Rejected"
	} else if anyPending {
		finalStatus = "Pending"
	} else {
		finalStatus = "Mentors Accepted"
	}

	// Update OD Application Status
	// Note: Only if status changed? PHP updates unconditionally.
	_, err = database.DB.Exec("UPDATE od_applications SET status=? WHERE id=?", finalStatus, odID)

	// 4. HOD Notification if "Mentors Accepted"
	if finalStatus == "Mentors Accepted" {
		// Fetch HOD Email
		var hodEmail, hodName, mentorEmail string
		// Logic from PHP: JOIN hods on department
		queryHOD := `
			SELECT h.email, h.name 
			FROM hods h
			JOIN od_applications o ON o.department = h.department
			WHERE o.id = ? LIMIT 1`
		err = database.DB.QueryRow(queryHOD, odID).Scan(&hodEmail, &hodName)
		
		// Fetch Mentor Email (My Email)
		_ = database.DB.QueryRow("SELECT mentor_email FROM mentors WHERE name = ?", mentorName).Scan(&mentorEmail)

		if err == nil && hodEmail != "" && mentorEmail != "" {
			subject := fmt.Sprintf("OD Application Approved by All Mentors (ID: %d)", odID)
			body := fmt.Sprintf("Dear HOD,<br><br>The OD application (ID: <strong>%d</strong>) has been <strong>ACCEPTED</strong> by all required mentors.<br>It is now pending your final approval.<br><br>Regards,<br>OD Management System", odID)
			
			// Send Email
			go services.SendEmail(hodEmail, subject, body)
		}
	}

	http.Redirect(w, r, "/mentor/dashboard", http.StatusSeeOther)
}

// Function parity note: check_all_mentors_accepted in PHP checks:
// SELECT * FROM od_team_members WHERE od_id = '$od_id'
// loop: if any member's mentor_status != 'Accepted', return false.
// My query `SELECT COUNT(*) ... WHERE mentor_status != 'Accepted'` is equivalent.
