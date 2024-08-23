package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	repos := []string{"/path/to/repo1", "/path/to/repo2"}

	for _, repo := range repos {
		if err := processRepo(repo); err != nil {
			fmt.Printf("Error processing repo %s: %v\n", repo, err)
		}
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
}

func processRepo(repoPath string) error {
	gitUser := os.Getenv("GIT_USER")
	gitEmail := os.Getenv("GIT_EMAIL")

	setGitConfig(repoPath, "user.name", gitUser)
	setGitConfig(repoPath, "user.email", gitEmail)

	pullRepo(repoPath)
	updateReadme(repoPath)
	commitChanges(repoPath)
	pushChanges(repoPath)

	return nil
}

func setGitConfig(repoPath, key, value string) {
	cmd := exec.Command("git", "-C", repoPath, "config", key, value)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error setting Git config %s: %v\n", key, err)
		os.Exit(1)
	}
}

func pullRepo(repoPath string) {
	cmd := exec.Command("git", "-C", repoPath, "pull")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error pulling repo %s: %v\n", repoPath, err)
		os.Exit(1)
	}
}

func updateReadme(repoPath string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	appendText := fmt.Sprintf("\nHola mundo - %s", currentTime)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' >> %s/README.md", appendText, repoPath))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error updating README in repo %s: %v\n", repoPath, err)
		os.Exit(1)
	}
}

func commitChanges(repoPath string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	commitMsg := fmt.Sprintf("Add 'Hola mundo' with date %s", currentTime)
	cmd := exec.Command("git", "-C", repoPath, "commit", "-am", commitMsg)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error committing changes in repo %s: %v\n", repoPath, err)
		os.Exit(1)
	}
}

func pushChanges(repoPath string) {
	cmd := exec.Command("git", "-C", repoPath, "push")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error pushing changes in repo %s: %v\n", repoPath, err)
		os.Exit(1)
	}
}
