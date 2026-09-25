package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"orvexa/internal/resume"

	"github.com/spf13/cobra"
)

//go:embed templates/editor.html
var editorTemplateContent string

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the Orvexa Live Form Editor",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := resume.ReadRaw()
			if err != nil {
				http.Error(w, "[ERROR] resume.json not found. Run 'orvexa init' first.", http.StatusNotFound)
				return
			}

			tmpl, err := template.New("editor").Parse(editorTemplateContent)
			if err != nil {
				http.Error(w, "[ERROR] failed to load editor template.", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, template.HTML(data)); err != nil {
				http.Error(w, "[ERROR] failed to render editor.", http.StatusInternalServerError)
				return
			}
		})

		http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "[ERROR] failed to read request.", http.StatusBadRequest)
				return
			}

			var temp resume.Resume
			if err := json.Unmarshal(body, &temp); err != nil {
				http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
				return
			}

			if err := resume.Write(&temp); err != nil {
				http.Error(w, "[ERROR] failed to save resume.", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		fmt.Println("[INFO] Orvexa Editor starting...")
		fmt.Println("[INFO] Local Network Link: http://localhost:9090")
		fmt.Println("[INFO] Press Ctrl+C to disconnect from the network.")

		if err := http.ListenAndServe(":9090", nil); err != nil {
			fmt.Printf("[ERROR] Editor server stopped: %v\n", err)
		}
	},
}
