// Package parser - Go list integration for dependency parsing
package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

// GoListPackage represents a package returned by 'go list -json'
type GoListPackage struct {
	Dir        string   `json:"Dir"`
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	GoFiles    []string `json:"GoFiles"`
	CgoFiles   []string `json:"CgoFiles"`
	Goroot     bool     `json:"Goroot"`
}

// listPackagesWithGoList uses 'go list' command to discover packages and dependencies
func (p *Parser) listPackagesWithGoList(ctx context.Context, dirs []string, args ...string) ([]*GoListPackage, error) {
	if !p.parseGoList {
		return nil, nil
	}

	var allPackages []*GoListPackage

	for _, dir := range dirs {
		pkgs, err := p.listOnePackage(ctx, dir, args...)
		if err != nil {
			return nil, err
		}
		allPackages = append(allPackages, pkgs...)
	}

	return allPackages, nil
}

// listOnePackage runs 'go list' for a single directory
func (p *Parser) listOnePackage(ctx context.Context, dir string, args ...string) ([]*GoListPackage, error) {
	cmdArgs := []string{"list", "-json"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, "./...")

	cmd := exec.CommandContext(ctx, "go", cmdArgs...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list failed: %v, stderr: %s", err, stderr.String())
	}

	// Parse JSON output (one JSON object per line)
	var packages []*GoListPackage
	decoder := json.NewDecoder(&stdout)

	for decoder.More() {
		var pkg GoListPackage
		if err := decoder.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode go list output: %v", err)
		}
		packages = append(packages, &pkg)
	}

	return packages, nil
}

// parsePackageFromGoList parses a package discovered by go list
func (p *Parser) parsePackageFromGoList(pkg *GoListPackage) error {
	// Skip internal packages unless parseInternal is true
	if pkg.Goroot && !p.parseInternal {
		return nil
	}

	srcDir := pkg.Dir

	// Parse .go files
	for _, file := range pkg.GoFiles {
		path := filepath.Join(srcDir, file)
		if err := p.ParseFile(path); err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}
	}

	// Parse .go files that import "C"
	for _, file := range pkg.CgoFiles {
		path := filepath.Join(srcDir, file)
		if err := p.ParseFile(path); err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}
	}

	return nil
}

// parseWithGoList is the main entry point for go list based parsing
func (p *Parser) parseWithGoList(dir string) error {
	if !p.parseGoList {
		return nil
	}

	ctx := context.Background()
	args := []string{}

	// Add -deps flag if parsing dependencies
	if p.parseDependency {
		args = append(args, "-deps")
	}

	packages, err := p.listPackagesWithGoList(ctx, []string{dir}, args...)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if err := p.parsePackageFromGoList(pkg); err != nil {
			return err
		}
	}

	return nil
}
