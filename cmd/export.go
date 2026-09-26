package cmd

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"orvexa/internal/resume"
	"os"
	"strings"
	"time"

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
		case "modern":
			tmplContent = modernTemplateContent
		case "minimal":
			tmplContent = minimalTemplateContent
		default:
			tmplContent = classicTemplateContent
		}

		fmt.Printf("[INFO] Generating %s PDF for %s...\n", templateName, res.Basics.Name)

		// 2. Setup Temporary HTTP Server for PDF Generation
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.New("resume").Funcs(template.FuncMap{
				"split": func(s, sep string) []string {
					res := strings.Split(s, sep)
					var final []string
					for _, v := range res {
						final = append(final, strings.TrimSpace(v))
					}
					return final
				},
			})

			tmpl, err := tmpl.Parse(tmplContent)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			tmpl.Execute(w, res)
		})

		server := &http.Server{Addr: ":7331", Handler: mux}
		go server.ListenAndServe()

		// 3. Use Rod to capture PDF
		time.Sleep(500 * time.Millisecond)
		edgePath := `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`
		u := launcher.New().
			Bin(edgePath).
			Headless(true).
			Leakless(false).
			MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.MustClose()

		page := browser.MustPage("http://localhost:7331")
		page.MustWaitLoad()

		pdfStream, err := page.PDF(&proto.PagePrintToPDF{
			PrintBackground: true,
			PaperWidth:      toPtr(8.27),
			PaperHeight:     toPtr(11.69),
		})

		if err != nil {
			fmt.Println("[ERROR] PDF Rendering failed:", err)
			server.Close()
			return
		}

		pdfData, _ := io.ReadAll(pdfStream)
		server.Close()

		// 4. Determine Filename & Save
		filename := outputFlag
		if filename == "" {
			filename = fmt.Sprintf("%s_%s_Resume.pdf", resume.SafeFilename(res.Basics.Name), strings.Title(templateName))
		}

		os.WriteFile(filename, pdfData, 0644)

		fmt.Printf("[SUCCESS] PDF Exported: %s\n", filename)
		fmt.Printf("[INFO] Using template: %s\n", templateName)
	},
}

// Helper for Rod library
func toPtr(f float64) *float64 { return &f }
