package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hieptle/gcp-guardrail/pkg/rego"
)

var (
	port        = flag.String("port", "8080", "HTTP server port")
	policiesDir = flag.String("policies", "../policies", "Path to Rego policies directory")
	templatesDir = flag.String("templates", "templates", "Path to HTML templates directory")
)

type PlaygroundServer struct {
	evaluator    *rego.Evaluator
	templates    *template.Template
	policiesDir  string
}

type EvaluationRequest struct {
	Input       string `json:"input"`
	PackagePath string `json:"package_path"`
	PolicyText  string `json:"policy_text,omitempty"`
}

type TemplateData struct {
	PolicyDirs  []string
	ExampleJSON string
}

func main() {
	flag.Parse()

	// Create the playground server
	server, err := NewPlaygroundServer(*policiesDir, *templatesDir)
	if err != nil {
		log.Fatalf("Failed to create playground server: %v", err)
	}

	// Set up HTTP handlers
	http.HandleFunc("/", server.handleIndex)
	http.HandleFunc("/evaluate", server.handleEvaluate)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the HTTP server
	log.Printf("Starting server on port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

// NewPlaygroundServer creates a new Rego playground server
func NewPlaygroundServer(policiesDir, templatesDir string) (*PlaygroundServer, error) {
	// Load templates
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create policy evaluator
	policyDirs := []string{policiesDir}
	evaluator, err := rego.NewEvaluator(context.Background(), policyDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy evaluator: %w", err)
	}

	return &PlaygroundServer{
		evaluator:    evaluator,
		templates:    templates,
		policiesDir:  policiesDir,
	}, nil
}

// handleIndex serves the main page
func (s *PlaygroundServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get available policy directories
	policyDirs, err := listDirectories(s.policiesDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list policy directories: %v", err), http.StatusInternalServerError)
		return
	}

	// Example JSON for the playground
	exampleJSON := `{
  "resource": {
    "type": "google_storage_bucket",
    "name": "example_bucket",
    "attributes": {
      "name": "example-bucket",
      "location": "US",
      "uniform_bucket_level_access": false,
      "public_access_prevention": "inherited"
    }
  }
}`

	data := TemplateData{
		PolicyDirs:  policyDirs,
		ExampleJSON: exampleJSON,
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
	}
}

// handleEvaluate evaluates policy against input
func (s *PlaygroundServer) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	var req EvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	// Parse the input JSON
	var input interface{}
	if err := json.Unmarshal([]byte(req.Input), &input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Evaluate the policy
	result, err := s.evaluator.Evaluate(req.PackagePath, input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Policy evaluation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// listDirectories returns a list of subdirectories in the specified directory
func listDirectories(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
} 