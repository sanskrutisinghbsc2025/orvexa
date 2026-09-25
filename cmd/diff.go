package cmd

import (
	"encoding/json"
	"fmt"
	"orvexa/internal/resume"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show all changes in your resume data",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Read Current from disk
		current, err := resume.Read()
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		// 2. Read Previous from Git
		r, err := git.PlainOpen(".")
		if err != nil {
			fmt.Println("[ERROR] Not a orvexa repo. Run 'orvexa init'")
			return
		}
		ref, err := r.Head()
		if err != nil {
			fmt.Println("[ERROR] No commit history found. Commit once first.")
			return
		}
		commit, _ := r.CommitObject(ref.Hash())
		tree, _ := commit.Tree()
		file, _ := tree.File("resume.json")
		prevData, _ := file.Contents()
		var prev resume.Resume
		json.Unmarshal([]byte(prevData), &prev)

		fmt.Printf("🔍 Diffing current changes against: %s\n", ref.Hash().String()[:7])
		fmt.Println("------------------------------------------------------------")
		hasChanges := false

		// --- CHECK BASICS ---
		if prev.Basics.Name != current.Basics.Name {
			fmt.Printf("[BASICS] Name: %s -> %s\n", prev.Basics.Name, current.Basics.Name)
			hasChanges = true
		}
		if prev.Basics.Phone != current.Basics.Phone {
			fmt.Printf("[BASICS] Phone: %s -> %s\n", prev.Basics.Phone, current.Basics.Phone)
			hasChanges = true
		}
		if prev.Basics.Email != current.Basics.Email {
			fmt.Printf("[BASICS] Email: %s -> %s\n", prev.Basics.Email, current.Basics.Email)
			hasChanges = true
		}

		// --- CHECK EXPERIENCE ---
		if len(prev.Experience) != len(current.Experience) {
			fmt.Printf("[EXP] Number of jobs changed: %d -> %d\n", len(prev.Experience), len(current.Experience))
			hasChanges = true
		} else {
			for i := range current.Experience {
				p, c := prev.Experience[i], current.Experience[i]
				if p.Company != c.Company {
					fmt.Printf("[EXP] Company: %s -> %s\n", p.Company, c.Company)
					hasChanges = true
				}
				if p.Role != c.Role {
					fmt.Printf("[EXP] Role at %s: %s -> %s\n", c.Company, p.Role, c.Role)
					hasChanges = true
				}
				if !reflect.DeepEqual(p.Points, c.Points) {
					fmt.Printf("[EXP] Bullet points updated for %s\n", c.Company)
					hasChanges = true
				}
			}
		}

		// --- CHECK EDUCATION ---
		if len(prev.Education) != len(current.Education) {
			fmt.Printf("[EDU] Schools changed.\n")
			hasChanges = true
		} else {
			for i := range current.Education {
				if prev.Education[i].School != current.Education[i].School {
					fmt.Printf("[EDU] School: %s -> %s\n", prev.Education[i].School, current.Education[i].School)
					hasChanges = true
				}
				if prev.Education[i].Degree != current.Education[i].Degree {
					fmt.Printf("[EDU] Degree at %s: %s -> %s\n", current.Education[i].School, prev.Education[i].Degree, current.Education[i].Degree)
					hasChanges = true
				}
			}
		}

		// --- CHECK SKILLS ---
		if !reflect.DeepEqual(prev.Skills, current.Skills) {
			fmt.Println("[SKILLS] Skills list was modified.")
			hasChanges = true
		}

		// --- CHECK PROJECTS ---
		if len(prev.Projects) != len(current.Projects) {
			fmt.Printf("[PROJECTS] Number of projects changed.\n")
			hasChanges = true
		} else {
			for i := range current.Projects {
				p, c := prev.Projects[i], current.Projects[i]
				if p.Name != c.Name {
					fmt.Printf("[PROJECTS] Name: %s -> %s\n", p.Name, c.Name)
					hasChanges = true
				}
				if !reflect.DeepEqual(p.Points, c.Points) {
					fmt.Printf("[PROJECTS] Details updated for %s\n", c.Name)
					hasChanges = true
				}
			}
		}

		if !hasChanges {
			fmt.Println("✨ No changes found. Your JSON on disk matches the Git history.")
		}
	},
}
