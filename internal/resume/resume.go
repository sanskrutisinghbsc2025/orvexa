// Package resume provides the canonical data model, I/O, and validation
// for Orvexa's resume.json schema.
package resume

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const ResumeFile = "resume.json"

// Resume is the canonical, single source-of-truth data model for Orvexa.
type Resume struct {
	Basics       Basics       `json:"basics"`
	SectionOrder []string     `json:"sectionOrder"`
	Education    []Education  `json:"education"`
	Skills       Skills       `json:"skills"`
	Experience   []Experience `json:"experience"`
	Projects     []Project    `json:"projects"`
}

// Basics holds top-level personal/contact information.
type Basics struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

// Education represents one academic institution.
type Education struct {
	School   string `json:"school"`
	Degree   string `json:"degree"`
	Date     string `json:"date"`
	CGPA     string `json:"cgpa"`
	Location string `json:"location,omitempty"`
}

// Experience represents one job or internship.
type Experience struct {
	Company  string   `json:"company"`
	Role     string   `json:"role"`
	Date     string   `json:"date"`
	Location string   `json:"location,omitempty"`
	Points   []string `json:"points"`
}

// Project represents one project entry.
type Project struct {
	Name   string   `json:"name"`
	Tech   string   `json:"tech"`
	Points []string `json:"points"`
	Link   string   `json:"link,omitempty"`
}

// Skills is a flexible key-value map (e.g. "Languages" → "Go, Python, TS").
// Using map[string]string matches the actual resume.json schema.
type Skills map[string]string

// Read loads and parses resume.json from the current directory.
func Read() (*Resume, error) {
	data, err := os.ReadFile(ResumeFile)
	if err != nil {
		return nil, fmt.Errorf("resume.json not found — run 'orvexa init' first")
	}
	var r Resume
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("resume.json is malformed JSON: %w", err)
	}
	// Apply defaults for missing optional fields
	if len(r.SectionOrder) == 0 {
		r.SectionOrder = []string{"education", "skills", "experience", "projects"}
	}
	return &r, nil
}

// Write serialises and persists the Resume to resume.json.
func Write(r *Resume) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize resume: %w", err)
	}
	return os.WriteFile(ResumeFile, data, 0644)
}

// ReadRaw returns the raw JSON bytes of resume.json.
func ReadRaw() ([]byte, error) {
	data, err := os.ReadFile(ResumeFile)
	if err != nil {
		return nil, fmt.Errorf("resume.json not found — run 'orvexa init' first")
	}
	return data, nil
}

// Validate performs semantic checks on a Resume and returns warnings.
func Validate(r *Resume) []string {
	var warnings []string

	if strings.TrimSpace(r.Basics.Name) == "" {
		warnings = append(warnings, "basics.name is empty")
	}
	if strings.TrimSpace(r.Basics.Email) == "" {
		warnings = append(warnings, "basics.email is empty")
	}
	for i, exp := range r.Experience {
		if len(exp.Points) == 0 {
			warnings = append(warnings, fmt.Sprintf("experience[%d] (%s) has no bullet points", i, exp.Company))
		}
	}
	for i, p := range r.Projects {
		if len(p.Points) == 0 {
			warnings = append(warnings, fmt.Sprintf("projects[%d] (%s) has no description points", i, p.Name))
		}
	}
	return warnings
}

// SafeFilename converts a person's name to a safe filename prefix.
// "John Doe" → "John_Doe"
func SafeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Resume"
	}
	return strings.ReplaceAll(name, " ", "_")
}

// WriteGitIgnore creates or overwrites the .gitignore file.
func WriteGitIgnore(content string) error {
	return os.WriteFile(".gitignore", []byte(content), 0644)
}
