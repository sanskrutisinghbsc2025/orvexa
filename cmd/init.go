package cmd

import (
	"fmt"
	"orvexa/internal/resume"
	"orvexa/internal/vcs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Orvexa resume network",
	Run: func(cmd *cobra.Command, args []string) {
		printBrand()

		// 1. Initialize Git via VCS package
		_, err := vcs.Init()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		// 2. Create the Mock Resume (John Doe)
		mockResume := resume.Resume{
			Basics: resume.Basics{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Phone:    "+1 555-0199",
				LinkedIn: "linkedin.com/in/johndoe",
				GitHub:   "github.com/johndoe",
			},
			SectionOrder: []string{"education", "skills", "experience", "projects"},
			Education: []resume.Education{
				{
					School:   "University of Technology",
					Degree:   "B.S. in Computer Science",
					Date:     "2018 - 2022",
					CGPA:     "3.9/4.0",
					Location: "San Francisco, CA",
				},
			},
			Skills: resume.Skills{
				"Languages": "Golang, Python, TypeScript, SQL",
				"Cloud":     "AWS, Docker, Kubernetes",
				"AI/ML":     "PyTorch, Scikit-Learn, OpenAI API",
			},
			Experience: []resume.Experience{
				{
					Company: "Tech Solutions Inc.",
					Role:    "Software Engineer",
					Date:    "2022 - Present",
					Points: []string{
						"Led development of a high-throughput data pipeline in Go.",
						"Reduced cloud infrastructure costs by 25% through container optimization.",
						"Mentored junior developers on best practices for version control.",
					},
				},
			},
			Projects: []resume.Project{
				{
					Name: "Distributed Crawler",
					Tech: "Golang, Redis, Docker",
					Points: []string{
						"Built a concurrent web crawler capable of processing 10k pages/minute.",
						"Implemented Redis-based deduplication logic to prevent redundant crawls.",
					},
				},
			},
		}

		err = resume.Write(&mockResume)
		if err != nil {
			fmt.Println("[ERROR] Failed to create resume.json:", err)
			return
		}

		// 3. Create .gitignore (Keeping it here for now as it's not resume data)
		ignore := "*.pdf\norvexa.exe\nnode_modules/\n.DS_Store\n"
		resume.WriteGitIgnore(ignore) // I'll add this helper to internal/resume

		fmt.Println("[SUCCESS] Orvexa network initialized successfully.")
		fmt.Println("[INFO] 'resume.json' has been seeded with a professional template.")
		fmt.Println("[INFO] Run 'orvexa edit' to begin tailoring your profile.")
	},
}
