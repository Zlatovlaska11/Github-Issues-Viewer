# 🐛 Issue Viewer

A simple command-line tool written in Go for viewing GitHub issues from any public repository.
## Features

    🔍 Search and display issues by repository

    📄 View issue titles, numbers, states (open/closed), and creation dates

    🧵 View issue body (optional)

    ⚡ Fast and lightweight

## Installation

    Clone the repo:
```bash
git clone https://github.com/yourusername/issue-viewer.git
cd issue-viewer
```
    Build the binary:
```bash
go build -o issue-viewer
```
Usage
```bash
./issue-viewer -repo YourRepo
``` 
Example:
```bash
./issue-viewer go
```
Flags
Flag Description
-repo 'the name of the repo'
-user 'username'

Example:

./issue-viewer go 

Requirements

    Go 1.18+

    Internet connection (uses GitHub API)
