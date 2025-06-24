package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <module-path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Outputs --ignore arguments for go-licenses based on .go-licenses-ignore\n")
		os.Exit(1)
	}

	modulePath := os.Args[1]

	// Convert module path to absolute path
	absModulePath, err := filepath.Abs(modulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving module path: %v\n", err)
		os.Exit(1)
	}

	// Find repository root (search upward from module path)
	repoRoot, err := findRepoRootFromPath(absModulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding repo root: %v\n", err)
		os.Exit(1)
	}

	// Build disallowed types and ignore arguments
	disallowedArgs, ignoreArgs, err := buildArgs(repoRoot, absModulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building args: %v\n", err)
		os.Exit(1)
	}

	// Combine all arguments
	allArgs := append(disallowedArgs, ignoreArgs...)

	// Output all arguments (space-separated for easy shell consumption)
	fmt.Print(strings.Join(allArgs, " "))
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func findRepoRootFromPath(startPath string) (string, error) {
	dir := startPath
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func buildArgs(repoRoot, modulePath string) ([]string, []string, error) {
	ignoreFile := filepath.Join(repoRoot, ".go-licenses-ignore")
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		return nil, nil, nil
	}

	file, err := os.Open(ignoreFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var disallowedArgs []string
	var disallowedTypes []string
	var ignorePackages []string

	scanner := bufio.NewScanner(file)

	// Get relative module path from repo root
	relModulePath, err := filepath.Rel(repoRoot, modulePath)
	if err != nil {
		return nil, nil, err
	}
	if relModulePath == "." {
		relModulePath = ""
	}

	var currentPrefix string
	var currentPrefixMatches bool

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)

		// Strip inline comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "!") {
			// Disallowed type: !forbidden
			currentPrefix = ""
			currentPrefixMatches = false
			disallowedType := line[1:]
			disallowedTypes = append(disallowedTypes, disallowedType)
		} else {
			// Check if line has a prefix
			var packagesPart string
			if colonIndex := strings.Index(line, ":"); colonIndex != -1 {
				// Line has prefix: update current prefix
				prefix := line[:colonIndex]
				packagesPart = strings.TrimSpace(line[colonIndex+1:])

				currentPrefix = prefix
				matched, err := doublestar.Match(prefix, relModulePath)
				if err != nil {
					currentPrefixMatches = false
				} else {
					currentPrefixMatches = matched
				}
			} else {
				// No prefix: use whole line as packages
				packagesPart = line
			}

			// Process packages if we should (either no current prefix, or current prefix matches)
			if currentPrefix == "" || currentPrefixMatches {
				if packagesPart != "" {
					packages := strings.Fields(packagesPart)
					for _, pkg := range packages {
						resolvedPackages, err := resolvePackageNames(pkg, repoRoot, relModulePath)
						if err != nil {
							continue // Skip unresolvable packages
						}
						ignorePackages = append(ignorePackages, resolvedPackages...)
					}
				}
			}
		}
	}

	// Build disallowed types argument
	if len(disallowedTypes) > 0 {
		disallowedArgs = append(disallowedArgs, "--disallowed_types", strings.Join(disallowedTypes, ","))
	}

	// Sort and deduplicate ignore packages
	slices.Sort(ignorePackages)
	ignorePackages = slices.Compact(ignorePackages)

	var ignoreArgs []string
	if len(ignorePackages) > 0 {
		ignoreArgs = append(ignoreArgs, "--ignore", strings.Join(ignorePackages, ","))
	}

	return disallowedArgs, ignoreArgs, scanner.Err()
}

func resolvePackageNames(packageName, repoRoot, currentModulePath string) ([]string, error) {
	// Handle self-reference
	if packageName == "." {
		name, err := getCurrentModuleName(repoRoot, currentModulePath)
		if err != nil {
			return nil, err
		}
		return []string{name}, nil
	}

	// Handle repository root relative paths starting with /
	if strings.HasPrefix(packageName, "/") {
		return expandLocalGlob(packageName[1:], repoRoot, "")
	}

	// Handle module relative paths starting with ./
	if strings.HasPrefix(packageName, "./") {
		return expandLocalGlob(packageName[2:], repoRoot, currentModulePath)
	}

	return []string{packageName}, nil
}

func expandLocalGlob(pattern, repoRoot, currentModulePath string) ([]string, error) {
	var basePath string
	if currentModulePath != "" {
		basePath = filepath.Join(repoRoot, currentModulePath)
	} else {
		basePath = repoRoot
	}

	// Find all directories matching the pattern
	fsys := os.DirFS(basePath)
	matches, err := doublestar.Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, match := range matches {
		// Convert relative match back to absolute path
		fullPath := filepath.Join(basePath, match)

		// Check if this is a directory
		info, err := os.Stat(fullPath)
		if err != nil || !info.IsDir() {
			continue
		}

		// Check if directory contains Go files
		if !containsGoFiles(fullPath) {
			continue
		}

		// Check if this directory has a go.mod file (is a module)
		goModPath := filepath.Join(fullPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Read module name from go.mod
			moduleName, err := readModuleName(goModPath)
			if err != nil {
				continue
			}
			results = append(results, moduleName)
		} else {
			// No go.mod found, treat as package within parent module
			var parentModuleName string
			if currentModulePath != "" {
				parentModuleName, err = getCurrentModuleName(repoRoot, currentModulePath)
			} else {
				parentModuleName, err = getCurrentModuleName(repoRoot, "")
			}
			if err != nil {
				continue
			}

			// Calculate relative path from module root to this directory
			var moduleRoot string
			if currentModulePath != "" {
				moduleRoot = filepath.Join(repoRoot, currentModulePath)
			} else {
				moduleRoot = repoRoot
			}

			relPath, err := filepath.Rel(moduleRoot, fullPath)
			if err != nil {
				continue
			}

			packagePath := parentModuleName + "/" + relPath
			results = append(results, packagePath)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no Go packages found matching pattern %s", pattern)
	}

	return results, nil
}

func containsGoFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			return true
		}
	}
	return false
}

func getCurrentModuleName(repoRoot, currentModulePath string) (string, error) {
	// Check if current module has its own go.mod
	if currentModulePath != "" {
		currentModuleDir := filepath.Join(repoRoot, currentModulePath)
		goModPath := filepath.Join(currentModuleDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return readModuleName(goModPath)
		}
	}

	// Fall back to root module
	rootGoMod := filepath.Join(repoRoot, "go.mod")
	if _, err := os.Stat(rootGoMod); err == nil {
		baseModuleName, err := readModuleName(rootGoMod)
		if err != nil {
			return "", err
		}
		if currentModulePath != "" {
			return baseModuleName + "/" + currentModulePath, nil
		}
		return baseModuleName, nil
	}

	return "", fmt.Errorf("no go.mod found")
}

func readModuleName(goModPath string) (string, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[7:]), nil
		}
	}

	return "", fmt.Errorf("module declaration not found in %s", goModPath)
}
