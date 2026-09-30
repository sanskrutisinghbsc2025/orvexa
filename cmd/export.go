package cmd

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"orvexa/internal/resume"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/spf13/cobra"
)

//go:embed templates/classic.html
var classicTemplateContent string

//go:embed templates/modern.html
var modernTemplateContent string

//go:embed templates/minimal.html
var minimalTemplateContent string

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("template", "t", "classic", "Template to use (classic, modern, minimal)")
	exportCmd.Flags().StringP("output", "o", "", "Output filename (e.g. MyResume.pdf)")
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your resume to a high-quality PDF",
	Run: func(cmd *cobra.Command, args []string) {
		templateName, _ := cmd.Flags().GetString("template")
		outputFlag, _ := cmd.Flags().GetString("output")

		// 1. Read Data
		res, err := resume.Read()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		// Choose template content
		var tmplContent string
		switch templateName {
		case "classic":
			tmplContent = classicTemplateContent
		case "modern":
			tmplContent = modernTemplateContent
		case "minimal":
			tmplContent = minimalTemplateContent
		default:
			fmt.Printf("[ERROR] Unknown template %q. Choose classic, modern or minimal.\n", templateName)
			return
		}

		tmpl, err := template.New("resume").Funcs(template.FuncMap{
			"split": func(s, sep string) []string {
				res := strings.Split(s, sep)
				var final []string
				for _, v := range res {
					final = append(final, strings.TrimSpace(v))
				}
				return final
			},
		}).Parse(tmplContent)
		if err != nil {
			fmt.Println("[ERROR] Failed to load template:", err)
			return
		}

		browserPath, err := findBrowser()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		fmt.Printf("[INFO] Generating %s PDF for %s...\n", templateName, res.Basics.Name)

		// 2. Setup Temporary HTTP Server for PDF Generation. Port 0 lets the OS
		// pick a free port, and the listener is ready before the browser starts.
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			fmt.Println("[ERROR] Could not start local render server:", err)
			return
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if err := tmpl.Execute(w, res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		server := &http.Server{Handler: mux}
		go server.Serve(listener)
		defer server.Close()

		// 3. Use Rod to capture PDF
		u, err := launcher.New().
			Bin(browserPath).
			Headless(true).
			Leakless(false).
			Launch()
		if err != nil {
			fmt.Println("[ERROR] Could not launch browser:", err)
			return
		}
		browser := rod.New().ControlURL(u)
		if err := browser.Connect(); err != nil {
			fmt.Println("[ERROR] Could not connect to browser:", err)
			return
		}
		defer browser.Close()

		page, err := browser.Page(proto.TargetCreateTarget{URL: "http://" + listener.Addr().String()})
		if err != nil {
			fmt.Println("[ERROR] Could not open render page:", err)
			return
		}
		if err := page.WaitLoad(); err != nil {
			fmt.Println("[ERROR] Render page failed to load:", err)
			return
		}

		pdfStream, err := page.PDF(&proto.PagePrintToPDF{
			PrintBackground: true,
			PaperWidth:      toPtr(8.27),
			PaperHeight:     toPtr(11.69),
		})
		if err != nil {
			fmt.Println("[ERROR] PDF Rendering failed:", err)
			return
		}

		pdfData, err := io.ReadAll(pdfStream)
		if err != nil {
			fmt.Println("[ERROR] Failed to read rendered PDF:", err)
			return
		}

		// 4. Determine Filename & Save
		filename := outputFlag
		if filename == "" {
			title := strings.ToUpper(templateName[:1]) + templateName[1:]
			filename = fmt.Sprintf("%s_%s_Resume.pdf", resume.SafeFilename(res.Basics.Name), title)
		}

		if err := os.WriteFile(filename, pdfData, 0644); err != nil {
			fmt.Println("[ERROR] Could not save PDF:", err)
			return
		}

		fmt.Printf("[SUCCESS] PDF Exported: %s\n", filename)
		fmt.Printf("[INFO] Using template: %s\n", templateName)
	},
}

// findBrowser locates a Chromium-based browser for PDF rendering. The
// ORVEXA_BROWSER environment variable takes priority, then Edge (both
// Program Files locations), then whatever rod can find (e.g. Chrome).
func findBrowser() (string, error) {
	var candidates []string
	if p := os.Getenv("ORVEXA_BROWSER"); p != "" {
		candidates = append(candidates, p)
	}
	for _, env := range []string{"ProgramFiles(x86)", "ProgramFiles"} {
		if base := os.Getenv(env); base != "" {
			candidates = append(candidates, filepath.Join(base, "Microsoft", "Edge", "Application", "msedge.exe"))
		}
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	if p, ok := launcher.LookPath(); ok {
		return p, nil
	}
	return "", fmt.Errorf("no Chrome or Edge found; install one or set ORVEXA_BROWSER to its path")
}

// Helper for Rod library
func toPtr(f float64) *float64 { return &f }
