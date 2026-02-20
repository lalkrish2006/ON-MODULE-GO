package models

import (
	"database/sql"
	"time"
)

// --- User Roles ---

type Student struct {
	RegisterNo string `db:"register_no"`
	Name       string `db:"name"`
	Password   string `db:"password"`
	Department string `db:"department"`
	Year       string `db:"year"`
	Section    string `db:"section"`
	Email      string `db:"email"`
}

type Mentor struct {
	RegisterNo  string `db:"register_no"`
	Password    string `db:"password"`
	Name        string `db:"name"`
	Department  string `db:"department"`
	Year        string `db:"year"`
	Section     string `db:"section"`
	MentorEmail string `db:"mentor_email"`
}

type HOD struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Email      string `db:"email"`
}

type Principal struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Email      string `db:"email"`
}

type AdminUser struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Department string `db:"department"`
}

type LabTechnician struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Email      string `db:"email"`
}

type CA struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Year       string `db:"year"`
	Section    string `db:"section"`
}

type JA struct {
	RegisterNo string `db:"register_no"`
	Password   string `db:"password"`
	Name       string `db:"name"`
	Department string `db:"department"`
}

// --- OD Application ---

type ODApplication struct {
	ID              int            `db:"id"`
	RegisterNo      string         `db:"register_no"`
	StudentName     string         `db:"student_name"`
	Year            string         `db:"year"`
	Department      string         `db:"department"`
	Section         string         `db:"section"`
	ODType          string         `db:"od_type"` // 'internal' or 'external'
	Purpose         string         `db:"purpose"`
	CollegeName     sql.NullString `db:"college_name"`
	EventName       sql.NullString `db:"event_name"`
	FromDate        sql.NullString `db:"from_date"`
	ToDate          sql.NullString `db:"to_date"`
	ODDate          sql.NullString `db:"od_date"`
	FromTime        sql.NullString `db:"from_time"`
	ToTime          sql.NullString `db:"to_time"`
	Status          string         `db:"status"`
	RequestBonafide int            `db:"request_bonafide"` // 1 or 0
	LabRequired     int            `db:"lab_required"`     // 1 or 0
	LabName         sql.NullString `db:"lab_name"`
	SystemRequired  int            `db:"system_required"` // 1 or 0
	CreatedAt       time.Time      `db:"created_at"`
}

type ODTeamMember struct {
	ID               int            `db:"id" json:"id"`
	ODID             int            `db:"od_id" json:"od_id"`
	MemberName       string         `db:"member_name" json:"member_name"`
	MemberRegNo      string         `db:"member_regno" json:"member_regno"`
	MemberDepartment string         `db:"member_department" json:"member_department"`
	MemberYear       string         `db:"member_year" json:"member_year"`
	MemberSection    string         `db:"member_section" json:"member_section"`
	Mentor           string         `db:"member_mentor" json:"mentor"`
	MentorStatus     sql.NullString `db:"mentor_status" json:"mentor_status"` // 'Pending', 'Accepted', 'Rejected'
}
