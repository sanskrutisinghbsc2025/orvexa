# Orvexa Technical Specification

This document outlines the internal architecture and engineering decisions behind the Orvexa Resume Versioning Network.

## 1. System Architecture
Orvexa is a Command Line Interface (CLI) application developed in Golang. The project follows a modular design pattern:
- **cmd/**: CLI command definitions utilizing the Cobra framework.
- **VCS Layer**: A programmatic wrapper for the `go-git` library.
- **UI Layer**: A Go-based HTTP server serving an interactive Vanilla JS form-to-JSON editor.
- **Production Layer**: Headless Chrome orchestration via the `go-rod` library.

## 2. VCS Implementation
Orvexa leverages the Git internal database to manage state. 
- **Object Mapping**: Structured JSON data is unmarshaled into Go structs, validated for schema compliance, and committed as Blobs to a hidden repository.
- **Ref Management**: Commands such as `branch` and `sync` manipulate Git Reference (Ref) pointers directly.
- **Sync Logic**: The `sync` command implements `git rebase` logic to allow specialized branches to be updated with the growth of the master branch while preserving local modifications.

## 3. Semantic Diff Engine
Standard Git diffs compare lines of text. Orvexa’s `diff` command performs a Field-Level Comparison:
- It utilizes Go’s `reflect` package to perform deep-equality checks on JSON objects.
- It identifies changes in specific nested objects (e.g., detecting a modified bullet point in an Experience array) rather than reporting standard line additions/deletions.

## 4. PDF Orchestration
To achieve a professional LaTeX-style aesthetic without requiring a LaTeX installation:
1. Orvexa initializes a local HTTP server.
2. It renders JSON data into a CSS-hardened HTML template.
3. It launches a headless browser instance (Chrome/Edge).
4. It executes a `PagePrintToPDF` protocol with 0.0 margins and A4 scaling to produce a print-ready document.

## 5. Intelligence Layer (AI)
The `review` command integrates the **Google Gemini-1.5-Flash** model. 
- **Prompt Engineering**: System prompts simulate a Senior Technical Recruiter at a Tier-1 tech firm.
- **Context Injection**: The AI is provided the structured JSON along with a `--role` flag to perform a targeted gap analysis between the candidate's skills and the target role's expectations.

## 6. Development Dependencies
- **CLI Framework**: `github.com/spf13/cobra`
- **Git Internals**: `github.com/go-git/go-git/v5`
- **Browser Engine**: `github.com/go-rod/rod`
- **AI Client**: `github.com/google/generative-ai-go`