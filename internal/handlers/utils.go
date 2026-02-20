package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"od-system/internal/utils"
	"path/filepath"
)

// RenderTemplate parses and executes a template with shared functions
func RenderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) {
	// Resolve the absolute path
	resolvedPath := utils.ResolvePath(tmplPath)

	funcMap := template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"get": func(m map[string]interface{}, key string) interface{} {
			if val, ok := m[key]; ok {
				return val
			}
			return ""
		},
		// Add other helpers here if needed, e.g. date formatting
	}

	name := filepath.Base(resolvedPath)
	tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(resolvedPath)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		// Log the error so we can see it in the terminal
		log.Printf("Template Execution Error (%s): %v", name, err)
		// Only try to send an error response/header if one hasn't been sent?
		// Effective Go http doesn't let us check fácilmente, but logging is the important part here.
		// We'll keep the http.Error logic but the log is key.
		// http.Error(w, "Template Execution Error: "+err.Error(), http.StatusInternalServerError)
	}
}
