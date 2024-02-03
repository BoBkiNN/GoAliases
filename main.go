package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	// Get the executable name without path
	executableName := os.Args[0]
	executableName = strings.TrimPrefix(executableName, "./")

	// Get alias file path from environment variable
	aliasFile := os.Getenv("GoAliasesFile")
	if aliasFile == "" {
		fmt.Println("No GoAliasesFile environment variable set")
		return
	}
	aliasFile = normalizePath(aliasFile)

	// Read the aliases from the file
	aliasMap, err := readAliases(aliasFile)
	if err != nil {
		fmt.Println("Error reading aliases:", err)
		os.Exit(1)
	}

	// Get the corresponding command for the executable
	command, found := aliasMap[executableName]
	if !found {
		fmt.Println("Alias not found for executable:", executableName)
		os.Exit(1)
	}

	// Build the command to run
	cmd := exec.Command(command, os.Args[1:]...)

	// Set standard input, output, and error streams
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		os.Exit(0)
	}
}

func readAliases(filename string) (map[string]string, error) {
	aliasMap := make(map[string]string)

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		// If the file doesn't exist, create an empty file
		file, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
		file.Close()
		return aliasMap, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			aliasMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return aliasMap, nil
}

func normalizePath(path string) string {
	// Expand user directory
	if strings.HasPrefix(path, "~\\") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}

	// Normalize path
	absPath, _ := filepath.Abs(path)
	return absPath
}
