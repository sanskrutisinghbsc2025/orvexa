package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

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
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

			// The JSON is embedded raw inside a <script> tag. "<" only occurs
			// inside JSON strings, where < is equivalent, so escaping it
			// stops a value like "</script>" from breaking out of the tag.
			safe := bytes.ReplaceAll(data, []byte("<"), []byte("\\u003c"))

			if err := tmpl.Execute(w, template.HTML(safe)); err != nil {
				http.Error(w, "[ERROR] failed to render editor.", http.StatusInternalServerError)
				return
			}
		})

		mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Browsers attach Origin to cross-site POSTs, so this stops other
			// web pages from overwriting resume.json via the local server.
			if origin := r.Header.Get("Origin"); origin != "" && !allowedHost(strings.TrimPrefix(origin, "http://")) {
				http.Error(w, "Forbidden origin", http.StatusForbidden)
				return
			}

			body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
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

		// Reject unexpected Host headers (DNS rebinding) and bind to loopback
		// only so the editor isn't reachable from the network.
		guarded := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowedHost(r.Host) {
				http.Error(w, "Forbidden host", http.StatusForbidden)
				return
			}
			mux.ServeHTTP(w, r)
		})

		if err := http.ListenAndServe("127.0.0.1:9090", guarded); err != nil {
			fmt.Printf("[ERROR] Editor server stopped: %v\n", err)
		}
	},
}

// allowedHost reports whether host is the local editor address.
func allowedHost(host string) bool {
	return host == "localhost:9090" || host == "127.0.0.1:9090"
}
