package handlers

import (
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/services"
	"strconv"
	"time"
)

// SubmitOD handles the OD application submission
func SubmitOD(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form error", http.StatusBadRequest)
		return
	}

	// Extract Form Data
	registerNo := r.FormValue("registerNo")
	studentName := r.FormValue("studentName")
	yearVal, _ := strconv.Atoi(r.FormValue("year"))
	department := r.FormValue("department")
	section := r.FormValue("section")
	mentor := r.FormValue("mentor")
	purpose := r.FormValue("purpose")
	odType := r.FormValue("odType")
	
	// Optional fields
	collegeName := r.FormValue("college_name")
	eventName := r.FormValue("event_name")
	
	// Date/Time Logic
	var fromDate, toDate, odDate, fromTime, toTime string
	// noOfDays removed

	if r.FormValue("fullDay") == "on" {
		odDate = r.FormValue("od_date")
	} else if r.FormValue("periodwise") == "on" {
		odDate = time.Now().Format("2006-01-02") // Or derived?
		// Periodwise logic
		fromTime = r.FormValue("from_time")
		toTime = r.FormValue("to_time")
	} else if r.FormValue("daywise") == "on" {
		fromDate = r.FormValue("from_date")
		toDate = r.FormValue("to_date")
	} else if odType == "external" {
		fromDate = r.FormValue("from_date_ext")
		toDate = r.FormValue("to_date_ext")
	}

	// Boolean conversions
	reqBonafide := 0
	if r.FormValue("request_bonafide") == "on" {
		reqBonafide = 1
	}
	labRequired := 0
	if r.FormValue("labRequired") == "on" {
		labRequired = 1
	}
	sysRequired := 0
	if r.FormValue("systemRequired") == "on" {
		sysRequired = 1
	}
	labName := r.FormValue("labName")

	// Insert OD Application
	query := `INSERT INTO od_applications 
		(register_no, student_name, year, department, section, od_type, purpose, 
		college_name, event_name, from_date, to_date, od_date, from_time, to_time, 
		status, request_bonafide, lab_required, lab_name, system_required, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	// Handle Nulls/Empty strings for SQL?
	// database/sql handles empty strings as empty strings. passing null explicitly requires sql.NullString.
	// We can pass empty string if allowed.
	
	// Helper to handle empty dates as valid NULL or string?
	// MySQL DATE column accepts NULL.
	// We should pass nil if empty.
	var fd, td, odD, ft, tt interface{}
	fd = fromDate
	if fromDate == "" { fd = nil }
	td = toDate
	if toDate == "" { td = nil }
	odD = odDate
	if odDate == "" { odD = nil }
	ft = fromTime
	if fromTime == "" { ft = nil }
	tt = toTime
	if toTime == "" { tt = nil }

	// Default status
	status := "Pending"

	res, err := database.DB.Exec(query, 
		registerNo, studentName, yearVal, department, section, odType, purpose,
		collegeName, eventName, fd, td, odD, ft, tt,
		status, reqBonafide, labRequired, labName, sysRequired,
	)

	if err != nil {
		log.Println("Insert OD Error:", err)
		http.Error(w, "Database Insert Error", 500)
		return
	}

	odID, _ := res.LastInsertId()

	// Insert Primary Student as Team Member?
	// The PHP logic: `od_team_members` contains the applicant AND other members?
	// Let's check `submit_od.php`.
	// Yes, "Insert Main Student into Team Members"
	
	teamQuery := `INSERT INTO od_team_members 
		(od_id, member_name, member_regno, member_department, member_year, member_section, mentor, mentor_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	// Main student status is Pending initially? Or auto-approved by self? 
	// The mentor status column refers to MENTOR approval. So "Pending".
	_, err = database.DB.Exec(teamQuery, odID, studentName, registerNo, department, yearVal, section, mentor, "Pending")
	if err != nil {
		log.Println("Insert main member error:", err)
	}

	// Notify Mentor
	// Fetch mentor email?
	// We don't have mentor email readily available. We should fetch it from `mentors` table.
	// Or pass it if we had it.
	// For now, call service logic which can fetch it.
	// Process Additional Team Members & Collect Info for Notification
	type TeamMemberInfo struct {
		Name   string
		Mentor string
	}
	var teamMembers []TeamMemberInfo

	teamCount, _ := strconv.Atoi(r.FormValue("teamCount"))
	for i := 0; i < teamCount; i++ {
		idx := strconv.Itoa(i)
		
		// Extract all fields for DB
		mReg := r.FormValue("member_regno_" + idx)
		mName := r.FormValue("member_name_" + idx)
		mYearVal, _ := strconv.Atoi(r.FormValue("member_year_" + idx))
		mDept := r.FormValue("member_department_" + idx)
		mSec := r.FormValue("member_section_" + idx)
		mMentor := r.FormValue("member_mentor_" + idx)

		// Insert into DB if valid
		if mReg != "" {
			_, err = database.DB.Exec(teamQuery, odID, mName, mReg, mDept, mYearVal, mSec, mMentor, "Pending")
			if err != nil {
				log.Println("Insert team member error:", err)
			}
			// Add to slice for notification
			if mName != "" && mMentor != "" {
				teamMembers = append(teamMembers, TeamMemberInfo{Name: mName, Mentor: mMentor})
			}
		}
	}

	// Async Notification
	go func(mainMentor, mainStudent string, members []TeamMemberInfo) {
		notifiedMentors := make(map[string]bool)

		notify := func(mName, sName string) {
			if mName == "" || notifiedMentors[mName] { return }
			var email string
			// Assuming mentors table has mentor_email column
			err := database.DB.QueryRow("SELECT mentor_email FROM mentors WHERE name = ?", mName).Scan(&email)
			if err == nil && email != "" {
				services.SendODNotification(email, sName, odType)
				notifiedMentors[mName] = true
			}
		}

		notify(mainMentor, mainStudent)
		for _, m := range members {
			notify(m.Mentor, m.Name)
		}
	}(mentor, studentName, teamMembers)

	// Process Additional Team Members (DB Insert already done above, just loop for logic consistency if needed, but we did DB insert in main thread)
	// Note: The loop regarding DB insert of team members was lines 155-171. 
	// We need to keep that DB insert logic separate from the async notification logic we just added above.
	// Wait, I am replacing the DB insert loop if I target lines 155-171! 
	// I must NOT remove the DB insert logic.
	// The user instruction "Replace logic with proper calls..." implies I should ADD the notification calls, not remove DB inserts.
	// Let's scroll up and see where I am inserting.
	// The target content I selected is the END of the function roughly. 
	// Ah, I need to be careful.
	// Let's replace the block *after* the DB inserts, or modify the DB insert loop to collect data for async?
	// The DB insert loop parses FormValue again.
	// Safe bet: Keep DB insert loop as is (it's synchronous).
	// Add the Notification block separately at the end.
	// But `submit_od.php` does it inline.
	// In Go, it's better to do it async to not block response.
	// I will just append the notification logic or wrap it.
	
	// Actually, let's leave the DB insert loop alone (previous tool call logic).
	// I will target the *existing* async notification block (lines 143-150) and expand it to include team members.
	// But the team member *data* is in the request form. I should parse it inside the goroutine or pass it?
	// Accessing `r.FormValue` inside goroutine is race-y if request is closed?
	// `r` might be closed when handler returns.
	// Correct way: Extract all needed data *before* goroutine.



	// Set Flash success
	session.Values["flash_success"] = "OD Application Submitted Successfully!"
	session.Save(r, w)

	http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
}
