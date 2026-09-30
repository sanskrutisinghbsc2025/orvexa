<div align='center'>
<img width="100" height="100" alt="WhatsApp_Image_2026-02-28_at_22 40 34-removebg-preview" src="https://github.com/user-attachments/assets/af959cbf-a738-4d9e-955e-3ff0d0e6baf5" />

# 🍄 Orvexa: The Resume Versioning Network

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Release](https://img.shields.io/github/v/release/sanskrutisinghbsc2025/orvexa)
![License](https://img.shields.io/badge/License-MIT-green)

Orvexa is a version control system designed to manage the lifecycle of professional resumes. Instead of tracking binary files, Orvexa tracks structured career data (JSON), allowing for semantic versioning, branching for specific job roles, and automated high-fidelity PDF generation.
</div>

## 🚀 Installation

### 1. The Quickest Way (Download Binary)
No coding knowledge required.
1. Go to the [Latest Releases](https://github.com/sanskrutisinghbsc2025/orvexa/releases/latest).
2. Download the `.zip` or `.tar.gz` file for your OS (e.g., `orvexa_Windows_x86_64.zip`).
3. Extract the `orvexa.exe` and move it to a folder in your **System PATH**.
4. Open your terminal and type `orvexa version`.

### 2. For Developers (Go installed)
```bash
go install github.com/sanskrutisinghbsc2025/orvexa@latest
```

---

## 🛠️ Quick Start

1. **Seed the Network**: Create a new folder and run `orvexa init`. This creates a professional John Doe template.
2. **Launch the Dashboard**: Run `orvexa edit`. A live form-based editor will open in your browser.
3. **Audit Your Profile**: Run `orvexa review --role="AI Engineer"`. Our integrated Gemini AI will critique your resume.
4. **Tailor with Branches**: Use `orvexa branch create company-name` to create specific versions of your CV.
5. **Deploy**: Run `orvexa export` to generate a perfectly formatted PDF.

---

## 🧠 Core Features

- **Semantic Diff**: Run `orvexa diff` to see human-readable changes between versions (e.g., "Changed Role from X to Y") instead of messy code lines.
- **Time Travel**: Use `orvexa restore <hash>` to instantly revert your resume to any previous state in your history.
- **Role-Specific Intelligence**: Deep integration with Google Gemini-1.5-Flash to provide specialized technical audits.

---

**1. Getting Started (The Onboarding)**
- **Installation (Non-Devs):** Step-by-step for Windows/Mac/Linux. Download from Releases -> Extract -> Add to PATH.
- **Installation (Go Devs):** Use `go install github.com/sanskrutisinghbsc2025/orvexa@latest`.
- **First Run:** Running `orvexa init` to bootstrap the network with the 'John Doe' professional template.

**2. The Management Suite (Editing & Exporting)**
- **`orvexa edit`:** Explain the visual dashboard. Mention it runs on `localhost:9090`, features a live split-screen preview, and auto-syncs with `resume.json`.
- **`orvexa export`:** Explain the headless browser orchestration. Mention it generates a pixel-perfect PDF using Jake's Resume format.

**3. Version Control (The Time-Machine)**
- **`orvexa commit -m "msg"`:** Saving a snapshot of the current state.
- **`orvexa status` & `list`:** Tracking the active branch and viewing the 40-character commit history.
- **`orvexa branch [create/switch]`:** The logic of specializing resumes. Explain the use case: creating a 'frontend-role' branch vs a 'backend-role' branch.
- **`orvexa diff`:** Explain the 'Semantic Diff' engine. Contrast it with raw Git diffs—show how Orvexa understands that a 'Role' changed, not just a line of text.
- **`orvexa restore [hash] --force`:** The 'Time Travel' command. Explain how to revert a resume to any point in the history.
- **`orvexa sync [branch]`:** The rebase logic. Explain how to pull updates from 'main' into a specialized branch.

**4. Intelligence Layer (AI Audit)**
- **`orvexa config --key [key]`:** Guide on obtaining a Google Gemini API key and storing it locally in `~/.orvexa_config.json`.
- **`orvexa review --role="[title]"`:** Detail the AI Audit feature. Explain how it provides a Role-Match score, missing keywords, and bullet point strengthening.

**5. Advanced Configuration**
- **Manual JSON Editing:** Explaining the JSON schema (Basics, Education, Experience, Skills, Projects).
- **The .gitignore File:** Why Orvexa ignores generated PDFs to keep the history clean.

---

## 🤝 Contributing & Feedback

The Orvexa network grows through community input! 

- **Found a Bug?** Open an [Issue](https://github.com/sanskrutisinghbsc2025/orvexa/issues). I will review and work on them actively.
- **Have a Feature Idea?** We would love to hear about new templates or AI features.
- **Contributing Code**: 
    1. Fork the repo.
    2. Create your feature branch.
    3. Open a Pull Request.
    *Assigning issues is not compulsory—feel free to pick up any open issue and start building!*

---

## ⚖️ License
Distributed under the MIT License. 

*Built with ❤️ by [Sanskruti Singh](https://github.com/sanskrutisinghbsc2025)*
