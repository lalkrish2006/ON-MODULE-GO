package handlers

import (
	"database/sql"
	"log"
	"od-system/internal/database"
	"od-system/internal/models"
)

// ODColumns explicitly lists fields matching scanODs expectation
const ODColumns = `o.id, o.register_no, o.student_name, o.year, o.department, o.section, 
	o.od_type, o.purpose, o.college_name, o.event_name, o.from_date, o.to_date, 
	o.od_date, o.from_time, o.to_time, o.status, o.request_bonafide, 
	o.lab_required, o.lab_name, o.system_required, o.created_at`

// scanODs helper to scan application rows and fetch team members
func scanODs(rows *sql.Rows, mentorName string) []DashboardOD {
	var ods []DashboardOD
	for rows.Next() {
		var od models.ODApplication
		err := rows.Scan(
			&od.ID, &od.RegisterNo, &od.StudentName, &od.Year, &od.Department, &od.Section,
			&od.ODType, &od.Purpose, &od.CollegeName, &od.EventName, &od.FromDate, &od.ToDate,
			&od.ODDate, &od.FromTime, &od.ToTime, &od.Status, &od.RequestBonafide,
			&od.LabRequired, &od.LabName, &od.SystemRequired, &od.CreatedAt,
		)
		if err != nil {
			log.Println("Scan Error in helper:", err)
			continue
		}

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

		// Determine Status for display
		displayStatus := "Pending"
		if od.Status != "" {
			displayStatus = od.Status
		}
		
		// Map specific status if mentorName is provided
		if mentorName != "" {
			for _, m := range teamMembers {
				if m.Mentor == mentorName {
					if m.MentorStatus.Valid {
						displayStatus = m.MentorStatus.String
					} else {
						displayStatus = "Pending"
					}
				}
			}
		}

        // Format Date
        dateStr := "-"
        if od.FromDate.Valid && od.ToDate.Valid {
             dateStr = od.FromDate.String + " to " + od.ToDate.String
        } else if od.ODDate.Valid {
            dateStr = od.ODDate.String
        }

		ods = append(ods, DashboardOD{
			ODApplication: od,
			DisplayStatus: displayStatus,
			TeamMembers:   teamMembers,
            DateStr:       dateStr,
		})
	}
	return ods
}
