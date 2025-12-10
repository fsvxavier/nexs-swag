// Package parser - Code examples support for x-codeSamples
package parser

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CodeExample represents a code sample for x-codeSamples extension
type CodeExample struct {
	Lang   string `json:"lang"`
	Source string `json:"source"`
}

// codeExamplesCache stores loaded code examples
var codeExamplesCache map[string]string
var codeExamplesCacheMutex sync.RWMutex

// loadCodeExamplesFromDir loads code example files from directory
func (p *Parser) loadCodeExamplesFromDir() error {
	if p.codeExampleFilesDir == "" {
		return nil
	}

	codeExamplesCacheMutex.Lock()
	codeExamplesCache = make(map[string]string)
	codeExamplesCacheMutex.Unlock()

	return filepath.Walk(p.codeExampleFilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Use relative path as key
		relPath, err := filepath.Rel(p.codeExampleFilesDir, path)
		if err != nil {
			relPath = filepath.Base(path)
		}

		codeExamplesCacheMutex.Lock()
		codeExamplesCache[relPath] = string(content)
		codeExamplesCacheMutex.Unlock()
		return nil
	})
}

// GetCodeExample returns a code example by filename
func (p *Parser) GetCodeExample(filename string) string {
	codeExamplesCacheMutex.RLock()
	defer codeExamplesCacheMutex.RUnlock()
	if codeExamplesCache == nil {
		return ""
	}
	return codeExamplesCache[filename]
}

// detectLanguageFromExtension returns language identifier from file extension
func detectLanguageFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	languageMap := map[string]string{
		".go":     "go",
		".js":     "javascript",
		".ts":     "typescript",
		".py":     "python",
		".java":   "java",
		".rb":     "ruby",
		".php":    "php",
		".cs":     "csharp",
		".cpp":    "cpp",
		".c":      "c",
		".sh":     "bash",
		".json":   "json",
		".yaml":   "yaml",
		".yml":    "yaml",
		".xml":    "xml",
		".html":   "html",
		".css":    "css",
		".sql":    "sql",
		".swift":  "swift",
		".kt":     "kotlin",
		".rs":     "rust",
		".dart":   "dart",
		".scala":  "scala",
		".groovy": "groovy",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}
	return "text"
}
