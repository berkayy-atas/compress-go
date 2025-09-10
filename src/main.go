package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func runCommand(name string, args []string, description string) error {
	start := time.Now()
	log.Printf("🚀 %s...", description)
	
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s failed: %v", description, err)
	}
	
	duration := time.Since(start).Round(time.Millisecond)
	log.Printf("✅ %s completed in %v", description, duration)
	return nil
}

func cleanupFile(path string) {
	if _, err := os.Stat(path); err == nil {
		log.Printf("🧹 Cleaning up: %s", path)
		os.RemoveAll(path)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	log.Println("🎯 Starting Repository Mirror Archive Action")
	
	// Environment variables
	githubToken := os.Getenv("GITHUB_TOKEN")
	repository := os.Getenv("GITHUB_REPOSITORY")
	
	if githubToken == "" {
		log.Fatal("❌ GITHUB_TOKEN environment variable is required")
	}
	if repository == "" {
		log.Fatal("❌ GITHUB_REPOSITORY environment variable is required")
	}
	
	// Configuration
	repoUrl := fmt.Sprintf("https://x-access-token:%s@github.com/%s.git", githubToken, repository)
	cloneDir := "repo-mirror"
	tarFile := "repo.tar"
	zstFile := "repo.tar.zst"
	
	// Cleanup previous runs
	cleanupFile(cloneDir)
	cleanupFile(tarFile)
	cleanupFile(zstFile)
	
	// Step 1: Clone repository with --mirror
	err := runCommand("git", []string{"clone", "--mirror", repoUrl, cloneDir}, 
		"Cloning repository with --mirror")
	if err != nil {
		log.Fatal(err)
	}
	
	// Step 2: Create tar archive
	err = runCommand("tar", []string{"-cf", tarFile, "-C", cloneDir, "."}, 
		"Creating tar archive")
	if err != nil {
		log.Fatal(err)
	}
	
	// Step 3: Compress with zstd
	err = runCommand("zstd", []string{"--rm", "-19", "-T0", tarFile, "-o", zstFile}, 
		"Compressing with zstd")
	if err != nil {
		log.Fatal(err)
	}
	
	// Step 4: Get file info
	if fileExists(zstFile) {
		err = runCommand("ls", []string{"-lh", zstFile}, 
			"Getting file information")
		if err != nil {
			log.Printf("⚠️ Could not get file info: %v", err)
		}
	}
	
	log.Printf("🎉 Action completed successfully!")
	log.Printf("💾 Final archive: %s", zstFile)
	
	// Check if file was created
	if fileExists(zstFile) {
		log.Printf("📦 Archive created successfully!")
	} else {
		log.Printf("❌ Archive file was not created")
	}
}