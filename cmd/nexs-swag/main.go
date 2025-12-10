// Package main is the entry point for nexs-swag CLI.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	pkgformat "github.com/fsvxavier/nexs-swag/pkg/format"
	"github.com/fsvxavier/nexs-swag/pkg/generator"
	"github.com/fsvxavier/nexs-swag/pkg/parser"
)

const (
	version = "1.0.0"
)

func main() {
	app := &cli.App{
		Name:    "nexs-swag",
		Usage:   "Generate OpenAPI 3.1.x specification from Go code comments",
		Version: version,
		Authors: []*cli.Author{
			{
				Name: "Fabricio Xavier",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize and generate OpenAPI documentation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Value:   "./",
						Usage:   "Directory to search for Go files",
					},
					&cli.StringFlag{
						Name:    "generalInfo",
						Aliases: []string{"g"},
						Value:   "",
						Usage:   "Go file path with general API annotations (e.g., main.go)",
					},
					&cli.StringFlag{
						Name:  "exclude",
						Value: "",
						Usage: "Exclude directories and files (comma-separated)",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "./docs",
						Usage:   "Output directory for generated files",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Value:   "json,yaml,go",
						Usage:   "Output formats (comma-separated): json, yaml, go",
					},
					&cli.StringFlag{
						Name:    "outputTypes",
						Aliases: []string{"ot"},
						Value:   "",
						Usage:   "Output types (alias for --format): go, json, yaml",
					},
					&cli.StringFlag{
						Name:    "propertyStrategy",
						Aliases: []string{"p"},
						Value:   "camelcase",
						Usage:   "Property naming strategy: snakecase, camelcase, pascalcase",
					},
					&cli.BoolFlag{
						Name:  "requiredByDefault",
						Value: false,
						Usage: "Set all fields as required by default",
					},
					&cli.BoolFlag{
						Name:  "parseInternal",
						Value: false,
						Usage: "Parse internal packages",
					},
					&cli.BoolFlag{
						Name:    "parseDependency",
						Aliases: []string{"pd"},
						Value:   false,
						Usage:   "Parse go dependencies in vendor folder",
					},
					&cli.IntFlag{
						Name:  "parseDepth",
						Value: 100,
						Usage: "Dependency parse depth",
					},
					&cli.StringFlag{
						Name:    "markdownFiles",
						Aliases: []string{"md"},
						Value:   "",
						Usage:   "Parse folder containing markdown files for descriptions",
					},
					&cli.StringFlag{
						Name:  "overridesFile",
						Value: ".swaggo",
						Usage: "File to read global type overrides from",
					},
					&cli.StringFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "Filter by tags (comma-separated, use ! to exclude)",
					},
					&cli.BoolFlag{
						Name:  "parseFuncBody",
						Value: false,
						Usage: "Parse annotations in function body",
					},
					&cli.BoolFlag{
						Name:  "parseVendor",
						Value: false,
						Usage: "Parse vendor folder",
					},
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},
						Value:   false,
						Usage:   "Suppress output",
					},
					&cli.BoolFlag{
						Name:  "validate",
						Value: true,
						Usage: "Validate OpenAPI specification after generation",
					},
					&cli.IntFlag{
						Name:    "parseDependencyLevel",
						Aliases: []string{"pdl"},
						Value:   0,
						Usage:   "Dependency parse level (0=disabled, 1=models, 2=operations, 3=all)",
					},
					&cli.StringFlag{
						Name:    "codeExampleFilesDir",
						Aliases: []string{"cef"},
						Value:   "",
						Usage:   "Directory containing code example files for x-codeSamples",
					},
					&cli.BoolFlag{
						Name:  "generatedTime",
						Value: false,
						Usage: "Generate timestamp in output",
					},
					&cli.StringFlag{
						Name:  "instanceName",
						Value: "swagger",
						Usage: "Name of the swagger instance",
					},
					&cli.BoolFlag{
						Name:  "parseGoList",
						Value: false,
						Usage: "Use 'go list' to parse dependencies",
					},
					&cli.StringFlag{
						Name:    "templateDelims",
						Aliases: []string{"td"},
						Value:   "",
						Usage:   "Custom template delimiters (format: 'left,right')",
					},
					&cli.StringFlag{
						Name:    "collectionFormat",
						Aliases: []string{"cf"},
						Value:   "csv",
						Usage:   "Default collection format (csv, multi, pipes, tsv, ssv)",
					},
					&cli.StringFlag{
						Name:  "parseExtension",
						Value: "",
						Usage: "Filter operations by extension prefix (e.g., x-)",
					},
					&cli.StringFlag{
						Name:  "state",
						Value: "",
						Usage: "State file for @HostState annotation",
					},
				},
				Action: initAction,
			},
			{
				Name:    "generate",
				Aliases: []string{"gen"},
				Usage:   "Generate OpenAPI documentation (alias for init)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Value:   "./",
						Usage:   "Directory to search for Go files",
					},
					&cli.StringFlag{
						Name:    "generalInfo",
						Aliases: []string{"g"},
						Value:   "",
						Usage:   "Go file path with general API annotations (e.g., main.go)",
					},
					&cli.StringFlag{
						Name:  "exclude",
						Value: "",
						Usage: "Exclude directories and files (comma-separated)",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "./docs",
						Usage:   "Output directory for generated files",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Value:   "json,yaml,go",
						Usage:   "Output formats (comma-separated): json, yaml, go",
					},
					&cli.StringFlag{
						Name:    "outputTypes",
						Aliases: []string{"ot"},
						Value:   "",
						Usage:   "Output types (alias for --format): go, json, yaml",
					},
					&cli.StringFlag{
						Name:    "propertyStrategy",
						Aliases: []string{"p"},
						Value:   "camelcase",
						Usage:   "Property naming strategy: snakecase, camelcase, pascalcase",
					},
					&cli.BoolFlag{
						Name:  "requiredByDefault",
						Value: false,
						Usage: "Set all fields as required by default",
					},
					&cli.BoolFlag{
						Name:  "parseInternal",
						Value: false,
						Usage: "Parse internal packages",
					},
					&cli.BoolFlag{
						Name:    "parseDependency",
						Aliases: []string{"pd"},
						Value:   false,
						Usage:   "Parse go dependencies in vendor folder",
					},
					&cli.IntFlag{
						Name:  "parseDepth",
						Value: 100,
						Usage: "Dependency parse depth",
					},
					&cli.StringFlag{
						Name:    "markdownFiles",
						Aliases: []string{"md"},
						Value:   "",
						Usage:   "Parse folder containing markdown files for descriptions",
					},
					&cli.StringFlag{
						Name:  "overridesFile",
						Value: ".swaggo",
						Usage: "File to read global type overrides from",
					},
					&cli.StringFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "Filter by tags (comma-separated, use ! to exclude)",
					},
					&cli.BoolFlag{
						Name:  "parseFuncBody",
						Value: false,
						Usage: "Parse annotations in function body",
					},
					&cli.BoolFlag{
						Name:  "parseVendor",
						Value: false,
						Usage: "Parse vendor folder",
					},
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},
						Value:   false,
						Usage:   "Suppress output",
					},
					&cli.BoolFlag{
						Name:  "validate",
						Value: true,
						Usage: "Validate OpenAPI specification after generation",
					},
					&cli.IntFlag{
						Name:    "parseDependencyLevel",
						Aliases: []string{"pdl"},
						Value:   0,
						Usage:   "Dependency parse level (0=disabled, 1=models, 2=operations, 3=all)",
					},
					&cli.StringFlag{
						Name:    "codeExampleFilesDir",
						Aliases: []string{"cef"},
						Value:   "",
						Usage:   "Directory containing code example files for x-codeSamples",
					},
					&cli.BoolFlag{
						Name:  "generatedTime",
						Value: false,
						Usage: "Generate timestamp in output",
					},
					&cli.StringFlag{
						Name:  "instanceName",
						Value: "swagger",
						Usage: "Name of the swagger instance",
					},
					&cli.BoolFlag{
						Name:  "parseGoList",
						Value: false,
						Usage: "Use 'go list' to parse dependencies",
					},
					&cli.StringFlag{
						Name:    "templateDelims",
						Aliases: []string{"td"},
						Value:   "",
						Usage:   "Custom template delimiters (format: 'left,right')",
					},
					&cli.StringFlag{
						Name:    "collectionFormat",
						Aliases: []string{"cf"},
						Value:   "csv",
						Usage:   "Default collection format (csv, multi, pipes, tsv, ssv)",
					},
					&cli.StringFlag{
						Name:  "parseExtension",
						Value: "",
						Usage: "Filter operations by extension prefix (e.g., x-)",
					},
					&cli.StringFlag{
						Name:  "state",
						Value: "",
						Usage: "State file for @HostState annotation",
					},
				},
				Action: initAction,
			},
			{
				Name:    "fmt",
				Aliases: []string{"f"},
				Usage:   "Format swagger comments in Go source files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Value:   "./",
						Usage:   "Directory to search for Go files",
					},
					&cli.BoolFlag{
						Name:  "quiet",
						Value: false,
						Usage: "Suppress output messages",
					},
				},
				Action: fmtAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initAction(c *cli.Context) error {
	searchDir := c.String("dir")
	generalInfo := c.String("generalInfo")
	exclude := c.String("exclude")
	outputDir := c.String("output")
	formatStr := c.String("format")
	outputTypes := c.String("outputTypes")
	propertyStrategy := c.String("propertyStrategy")
	requiredByDefault := c.Bool("requiredByDefault")
	parseInternal := c.Bool("parseInternal")
	parseDependency := c.Bool("parseDependency")
	parseDepth := c.Int("parseDepth")
	markdownFiles := c.String("markdownFiles")
	overridesFile := c.String("overridesFile")
	tags := c.String("tags")
	parseFuncBody := c.Bool("parseFuncBody")
	parseVendor := c.Bool("parseVendor")
	quiet := c.Bool("quiet")
	validate := c.Bool("validate")
	parseDependencyLevel := c.Int("parseDependencyLevel")
	codeExampleFilesDir := c.String("codeExampleFilesDir")
	generatedTime := c.Bool("generatedTime")
	instanceName := c.String("instanceName")
	parseGoList := c.Bool("parseGoList")
	templateDelims := c.String("templateDelims")
	collectionFormat := c.String("collectionFormat")
	parseExtension := c.String("parseExtension")
	state := c.String("state")

	// Use outputTypes if format is not explicitly set
	if outputTypes != "" && formatStr == "json,yaml,go" {
		formatStr = outputTypes
	}

	// Parse formats
	formats := strings.Split(formatStr, ",")
	for i, format := range formats {
		formats[i] = strings.TrimSpace(format)
	}

	// Parse exclude patterns
	var excludePatterns []string
	if exclude != "" {
		excludePatterns = strings.Split(exclude, ",")
		for i, pattern := range excludePatterns {
			excludePatterns[i] = strings.TrimSpace(pattern)
		}
	}

	// Parse tag filters
	var includeTags, excludeTags []string
	if tags != "" {
		tagList := strings.Split(tags, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if strings.HasPrefix(tag, "!") {
				excludeTags = append(excludeTags, tag[1:])
			} else {
				includeTags = append(includeTags, tag)
			}
		}
	}

	if !quiet {
		fmt.Printf("nexs-swag v%s - OpenAPI 3.1.x Generator\n\n", version)
		fmt.Printf("Parsing directory: %s\n", searchDir)
		if generalInfo != "" {
			fmt.Printf("General info file: %s\n", generalInfo)
		}
		if len(excludePatterns) > 0 {
			fmt.Printf("Excluding: %v\n", excludePatterns)
		}
		if propertyStrategy != "camelcase" {
			fmt.Printf("Property strategy: %s\n", propertyStrategy)
		}
	}

	// Create parser with options
	p := parser.New()

	// Set parser options
	p.SetGeneralInfoFile(generalInfo)
	p.SetExcludePatterns(excludePatterns)
	p.SetPropertyStrategy(propertyStrategy)
	p.SetRequiredByDefault(requiredByDefault)
	p.SetParseInternal(parseInternal)
	p.SetParseDependency(parseDependency)
	p.SetParseDepth(parseDepth)
	p.SetMarkdownFilesDir(markdownFiles)
	p.SetOverridesFile(overridesFile)
	p.SetTagFilters(includeTags, excludeTags)
	p.SetParseFuncBody(parseFuncBody)
	p.SetParseVendor(parseVendor)
	p.SetParseDependencyLevel(parseDependencyLevel)
	p.SetCodeExampleFilesDir(codeExampleFilesDir)
	p.SetGeneratedTime(generatedTime)
	p.SetInstanceName(instanceName)
	p.SetParseGoList(parseGoList)
	p.SetTemplateDelims(templateDelims)
	p.SetCollectionFormat(collectionFormat)
	p.SetParseExtension(parseExtension)
	p.SetState(state)

	// Parse directory
	if err := p.ParseDir(searchDir); err != nil {
		return fmt.Errorf("failed to parse directory: %w", err)
	}

	// Validate if requested
	if validate {
		if !quiet {
			fmt.Println("Validating OpenAPI specification...")
		}
		if err := p.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		if !quiet {
			fmt.Println("✓ Validation passed")
		}
	}

	// Get OpenAPI specification
	spec := p.GetOpenAPI()

	// Generate output files
	if !quiet {
		fmt.Printf("Generating documentation in: %s\n", outputDir)
	}
	gen := generator.New(spec, outputDir, formats)
	gen.SetInstanceName(instanceName)
	gen.SetGeneratedTime(generatedTime)
	gen.SetTemplateDelims(templateDelims)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	if !quiet {
		fmt.Println("\n✓ Documentation generated successfully!")
	}
	return nil
}

func fmtAction(c *cli.Context) error {
	searchDir := c.String("dir")
	quiet := c.Bool("quiet")

	if !quiet {
		fmt.Printf("Formatting swagger comments in: %s\n", searchDir)
	}

	// Use the complete format package
	formatter := pkgformat.New()

	config := &pkgformat.Config{
		SearchDir: searchDir,
		Excludes:  "",
	}

	if err := formatter.Build(config); err != nil {
		return fmt.Errorf("failed to format files: %w", err)
	}

	if !quiet {
		fmt.Println("✓ All files formatted successfully!")
	}
	return nil
}
