package cmd

import (
	"encoding/json"
	"fmt"
	"orvexa/internal/resume"
	"reflect"
	"sort"

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
		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			fmt.Println("[ERROR] Could not read the latest version:", err)
			return
		}
		tree, err := commit.Tree()
		if err != nil {
			fmt.Println("[ERROR] Could not read the latest version:", err)
			return
		}
		file, err := tree.File(resume.ResumeFile)
		if err != nil {
			fmt.Println("[ERROR] resume.json is not in the latest version. Commit it first.")
			return
		}
		prevData, err := file.Contents()
		if err != nil {
			fmt.Println("[ERROR] Could not read resume.json from the latest version:", err)
			return
		}
		var prev resume.Resume
		if err := json.Unmarshal([]byte(prevData), &prev); err != nil {
			fmt.Println("[ERROR] resume.json in the latest version is malformed:", err)
			return
		}

		fmt.Printf("🔍 Diffing current changes against: %s\n", ref.Hash().String()[:7])
		fmt.Println("------------------------------------------------------------")
		hasChanges := false

		// field prints one line and records a change when old != new.
		field := func(tag, label, old, new string) {
			if old != new {
				fmt.Printf("[%s] %s: %s -> %s\n", tag, label, old, new)
				hasChanges = true
			}
		}
		note := func(format string, a ...any) {
			fmt.Printf(format+"\n", a...)
			hasChanges = true
		}

		// --- CHECK BASICS ---
		field("BASICS", "Name", prev.Basics.Name, current.Basics.Name)
		field("BASICS", "Phone", prev.Basics.Phone, current.Basics.Phone)
		field("BASICS", "Email", prev.Basics.Email, current.Basics.Email)
		field("BASICS", "LinkedIn", prev.Basics.LinkedIn, current.Basics.LinkedIn)
		field("BASICS", "GitHub", prev.Basics.GitHub, current.Basics.GitHub)

		// --- CHECK SECTION ORDER ---
		// resume.Read() fills in a default order, so skip the comparison when
		// the old version never had one (would be a false positive).
		if len(prev.SectionOrder) > 0 && !reflect.DeepEqual(prev.SectionOrder, current.SectionOrder) {
			note("[ORDER] Section order: %v -> %v", prev.SectionOrder, current.SectionOrder)
		}

		// --- CHECK EXPERIENCE ---
		if len(prev.Experience) != len(current.Experience) {
			note("[EXP] Number of jobs changed: %d -> %d", len(prev.Experience), len(current.Experience))
		} else {
			for i := range current.Experience {
				p, c := prev.Experience[i], current.Experience[i]
				field("EXP", "Company", p.Company, c.Company)
				field("EXP", "Role at "+c.Company, p.Role, c.Role)
				field("EXP", "Date at "+c.Company, p.Date, c.Date)
				field("EXP", "Location at "+c.Company, p.Location, c.Location)
				if !reflect.DeepEqual(p.Points, c.Points) {
					note("[EXP] Bullet points updated for %s", c.Company)
				}
			}
		}

		// --- CHECK EDUCATION ---
		if len(prev.Education) != len(current.Education) {
			note("[EDU] Number of schools changed: %d -> %d", len(prev.Education), len(current.Education))
		} else {
			for i := range current.Education {
				p, c := prev.Education[i], current.Education[i]
				field("EDU", "School", p.School, c.School)
				field("EDU", "Degree at "+c.School, p.Degree, c.Degree)
				field("EDU", "Date at "+c.School, p.Date, c.Date)
				field("EDU", "CGPA at "+c.School, p.CGPA, c.CGPA)
				field("EDU", "Location at "+c.School, p.Location, c.Location)
			}
		}

		// --- CHECK SKILLS ---
		if !reflect.DeepEqual(prev.Skills, current.Skills) {
			keys := map[string]bool{}
			for k := range prev.Skills {
				keys[k] = true
			}
			for k := range current.Skills {
				keys[k] = true
			}
			sorted := make([]string, 0, len(keys))
			for k := range keys {
				sorted = append(sorted, k)
			}
			sort.Strings(sorted)
			for _, k := range sorted {
				pv, inPrev := prev.Skills[k]
				cv, inCur := current.Skills[k]
				switch {
				case !inPrev:
					note("[SKILLS] Added %s: %s", k, cv)
				case !inCur:
					note("[SKILLS] Removed %s", k)
				case pv != cv:
					note("[SKILLS] %s: %s -> %s", k, pv, cv)
				}
			}
		}

		// --- CHECK PROJECTS ---
		if len(prev.Projects) != len(current.Projects) {
			note("[PROJECTS] Number of projects changed: %d -> %d", len(prev.Projects), len(current.Projects))
		} else {
			for i := range current.Projects {
				p, c := prev.Projects[i], current.Projects[i]
				field("PROJECTS", "Name", p.Name, c.Name)
				field("PROJECTS", "Tech for "+c.Name, p.Tech, c.Tech)
				field("PROJECTS", "Link for "+c.Name, p.Link, c.Link)
				if !reflect.DeepEqual(p.Points, c.Points) {
					note("[PROJECTS] Details updated for %s", c.Name)
				}
			}
		}

		if !hasChanges {
			fmt.Println("✨ No changes found. Your JSON on disk matches the Git history.")
		}
	},
}
