package handlers

import (
	"fmt"
	"log"
	"net/http"
	"od-system/internal/database"
	"od-system/internal/services"
	"od-system/internal/utils"
	"strings"
)

// TableConfig defines the structure for user role tables
type TableConfig struct {
	Table  string
	PK     string
	Fields []string
}

// UserRolesConfig holds the configuration for each role
var UserRolesConfig = map[string]TableConfig{
	"student":        {Table: "students", PK: "register_no", Fields: []string{"name", "department", "year", "section"}},
	"hod":            {Table: "hods", PK: "register_no", Fields: []string{"name", "department"}},
	"mentor":         {Table: "mentors", PK: "register_no", Fields: []string{"name", "department", "year", "section"}},
	"admin":          {Table: "admin_users", PK: "register_no", Fields: []string{"name", "department"}},
	"principal":      {Table: "principals", PK: "register_no", Fields: []string{"name"}},
	"lab_technician": {Table: "lab_technicians", PK: "register_no", Fields: []string{"name", "department"}},
	"cas":            {Table: "cas", PK: "register_no", Fields: []string{"name", "department", "year", "section"}},
	"jas":            {Table: "jas", PK: "register_no", Fields: []string{"name", "department"}},
}

type AdminDashboardData struct {
	User         map[string]interface{}
	Roles        map[string]TableConfig
	SelectedRole string
	TableConfig  *TableConfig
	UserData     []map[string]interface{}
	Message      string
	MessageType  string // "success" or "error"
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	if role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	selectedRole := r.URL.Query().Get("role")
	message := r.URL.Query().Get("message")
	messageType := r.URL.Query().Get("type")

	var tableConfig *TableConfig
	var userData []map[string]interface{}

	if config, ok := UserRolesConfig[selectedRole]; ok {
		tableConfig = &config

		// Build Query
		columns := append([]string{config.PK}, config.Fields...)
		query := fmt.Sprintf("SELECT %s FROM %s ORDER BY name ASC", strings.Join(columns, ", "), config.Table)

		rows, err := database.DB.Query(query)
		if err != nil {
			log.Println("Admin Query Error:", err)
			http.Error(w, "Database error", 500)
			return
		}
		defer rows.Close()

		// Dynamic Fetch
		cols, _ := rows.Columns()
		for rows.Next() {
			// Create a slice of interface{} to hold values
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			// Map values to column names
			rowMap := make(map[string]interface{})
			for i, col := range cols {
				val := values[i]
				b, ok := val.([]byte)
				var v interface{}
				if ok {
					v = string(b)
				} else {
					v = val
				}
				rowMap[col] = v
			}
			userData = append(userData, rowMap)
		}
	}

	userMap := make(map[string]interface{})
	for k, v := range session.Values {
		if strKey, ok := k.(string); ok {
			userMap[strKey] = v
		}
	}

	data := AdminDashboardData{
		User:         userMap,
		Roles:        UserRolesConfig,
		SelectedRole: selectedRole,
		TableConfig:  tableConfig,
		UserData:     userData,
		Message:      message,
		MessageType:  messageType,
	}

	RenderTemplate(w, "templates/admin_dashboard.html", data)
}

func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", 400)
		return
	}

	updateRole := r.FormValue("role")
	identifier := r.FormValue("identifier")

	config, ok := UserRolesConfig[updateRole]
	if !ok || identifier == "" {
		http.Redirect(w, r, "/admin/dashboard?error=invalid_request", http.StatusSeeOther)
		return
	}

	// Build Update Query
	var setClauses []string
	var args []interface{}

	for _, field := range config.Fields {
		val := r.FormValue(field)
		// We update even if empty, as user might want to clear a field.
		// PHP code checks `isset($_POST[$field])` which is true for empty string in form post
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, val)
	}

	args = append(args, identifier) // Add PK for WHERE clause

	if len(setClauses) > 0 {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", config.Table, strings.Join(setClauses, ", "), config.PK)
		_, err := database.DB.Exec(query, args...)

		msg := ""
		msgType := ""
		if err != nil {
			log.Println("Admin Update Error:", err)
			msg = "Error updating data: " + err.Error()
			msgType = "error"
		} else {
			msg = fmt.Sprintf("Data for %s in %s updated successfully!", identifier, updateRole)
			msgType = "success"
		}

		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=%s&type=%s", updateRole, msg, msgType), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s", updateRole), http.StatusSeeOther)
}

func AdminAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", 400)
		return
	}

	addRole := r.FormValue("role")
	identifier := strings.TrimSpace(r.FormValue("identifier"))
	rawPassword := strings.TrimSpace(r.FormValue("password"))

	config, ok := UserRolesConfig[addRole]
	if !ok || identifier == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=Invalid+request&type=error", addRole), http.StatusSeeOther)
		return
	}

	if rawPassword == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=Password+is+required&type=error", addRole), http.StatusSeeOther)
		return
	}

	// Date picker sends YYYY-MM-DD; user requested to keep this format for login.
	dobForHashing := rawPassword

	// Hash password using MD5 (matches login logic)
	hashedPassword := utils.HashPasswordMD5(dobForHashing)

	// Build INSERT query: PK + role fields + password column
	allCols := append([]string{config.PK}, config.Fields...)
	allCols = append(allCols, "password")

	placeholders := make([]string, len(allCols))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	args := make([]interface{}, len(allCols))
	args[0] = identifier
	for i, field := range config.Fields {
		args[i+1] = r.FormValue(field)
	}
	args[len(allCols)-1] = hashedPassword

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		config.Table,
		strings.Join(allCols, ", "),
		strings.Join(placeholders, ", "))

	_, err := database.DB.Exec(query, args...)
	msg := ""
	msgType := ""
	if err != nil {
		log.Println("Admin Add Error:", err)
		msg = "Error adding entry: " + err.Error()
		msgType = "error"
	} else {
		msg = fmt.Sprintf("New %s entry '%s' added successfully!", addRole, identifier)
		msgType = "success"
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=%s&type=%s", addRole, msg, msgType), http.StatusSeeOther)
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", 400)
		return
	}

	delRole := r.FormValue("role")
	identifier := strings.TrimSpace(r.FormValue("identifier"))

	config, ok := UserRolesConfig[delRole]
	if !ok || identifier == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=Invalid+request&type=error", delRole), http.StatusSeeOther)
		return
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", config.Table, config.PK)
	_, err := database.DB.Exec(query, identifier)

	msg := ""
	msgType := ""
	if err != nil {
		log.Println("Admin Delete Error:", err)
		msg = "Error removing entry: " + err.Error()
		msgType = "error"
	} else {
		msg = fmt.Sprintf("Entry '%s' removed from %s successfully!", identifier, delRole)
		msgType = "success"
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/dashboard?role=%s&message=%s&type=%s", delRole, msg, msgType), http.StatusSeeOther)
}

func AdminViewODs(w http.ResponseWriter, r *http.Request) {
	session := services.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	role, _ := session.Values["role"].(string)
	registerNo, _ := session.Values["register_no"].(string)

	if role != "admin" {
		http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
		return
	}

	// Fetch Department
	var department string
	err := database.DB.QueryRow("SELECT department FROM admin_users WHERE register_no = ?", registerNo).Scan(&department)
	if err != nil {
		log.Println("Admin View ODs Error:", err)
		http.Error(w, "Could not fetch admin department", 500)
		return
	}

	// Redirect to HOD Dashboard
	http.Redirect(w, r, "/hod/dashboard?department="+department+"&access=admin", http.StatusSeeOther)
}
