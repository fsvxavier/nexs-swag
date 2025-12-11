package parser

import (
	"regexp"
	"strings"

	openapi "github.com/fsvxavier/nexs-swag/pkg/openapi/v3"
)

// GeneralInfoProcessor processes general API information annotations.
type GeneralInfoProcessor struct {
	openapi *openapi.OpenAPI
	lastTag *openapi.Tag
}

// NewGeneralInfoProcessor creates a new general info processor.
func NewGeneralInfoProcessor(spec *openapi.OpenAPI) *GeneralInfoProcessor {
	return &GeneralInfoProcessor{
		openapi: spec,
	}
}

var (
	// General Info regex patterns.
	titleRegex       = regexp.MustCompile(`^@title\s+(.+)$`)
	versionRegex     = regexp.MustCompile(`^@version\s+(.+)$`)
	descriptionRegex = regexp.MustCompile(`^@description\s+(.+)$`)
	summaryRegex     = regexp.MustCompile(`^@summary\s+(.+)$`)
	tosRegex         = regexp.MustCompile(`^@termsOfService\s+(.+)$`)

	// Contact regex patterns.
	contactNameRegex  = regexp.MustCompile(`^@contact\.name\s+(.+)$`)
	contactURLRegex   = regexp.MustCompile(`^@contact\.url\s+(.+)$`)
	contactEmailRegex = regexp.MustCompile(`^@contact\.email\s+(.+)$`)

	// License regex patterns.
	licenseNameRegex = regexp.MustCompile(`^@license\.name\s+(.+)$`)
	licenseURLRegex  = regexp.MustCompile(`^@license\.url\s+(.+)$`)
	licenseIDRegex   = regexp.MustCompile(`^@license\.identifier\s+(.+)$`)

	// Server regex patterns.
	hostRegex       = regexp.MustCompile(`^@host\s+(\S+)$`)
	basePathRegex   = regexp.MustCompile(`^@basePath\s+(\S+)$`)
	schemesRegex    = regexp.MustCompile(`^@schemes\s+(.+)$`)
	serverRegex     = regexp.MustCompile(`^@server\s+(\S+)\s*(.*)$`)
	serverDescRegex = regexp.MustCompile(`^@server\.description\s+(.+)$`)

	// Tag regex patterns.
	tagNameRegex     = regexp.MustCompile(`^@tag\.name\s+(.+)$`)
	tagDescRegex     = regexp.MustCompile(`^@tag\.description\s+(.+)$`)
	tagDocsURLRegex  = regexp.MustCompile(`^@tag\.docs\.url\s+(.+)$`)
	tagDocsDescRegex = regexp.MustCompile(`^@tag\.docs\.description\s+(.+)$`)

	// Security regex patterns.
	securityRegex       = regexp.MustCompile(`^@securityDefinitions\.(\w+)\s+(\S+)\s+(.+)$`)
	securityBasicRegex  = regexp.MustCompile(`^@securityDefinitions\.basic\s+(\S+)\s*(.*)$`)
	securityAPIKeyRegex = regexp.MustCompile(`^@securityDefinitions\.apikey\s+(\S+)\s+(\w+)\s+(\w+)\s*(.*)$`)
	securityOAuth2Regex = regexp.MustCompile(`^@securityDefinitions\.oauth2\.(\w+)\s+(\S+)\s*(.*)$`)

	// External docs regex patterns.
	externalDocsURLRegex  = regexp.MustCompile(`^@externalDocs\.url\s+(.+)$`)
	externalDocsDescRegex = regexp.MustCompile(`^@externalDocs\.description\s+(.+)$`)

	// Webhook regex patterns.
	webhookRegex = regexp.MustCompile(`^@webhook\s+(\S+)\s+(.+)$`)
)

// Process processes a single general API annotation.
//
//nolint:gocyclo // Complex switch for many annotation types is acceptable
func (g *GeneralInfoProcessor) Process(text string) error {
	switch {
	// Info object
	case titleRegex.MatchString(text):
		matches := titleRegex.FindStringSubmatch(text)
		g.openapi.Info.Title = matches[1]

	case versionRegex.MatchString(text):
		matches := versionRegex.FindStringSubmatch(text)
		g.openapi.Info.Version = matches[1]

	case descriptionRegex.MatchString(text):
		matches := descriptionRegex.FindStringSubmatch(text)
		if g.openapi.Info.Description == "" {
			g.openapi.Info.Description = matches[1]
		} else {
			g.openapi.Info.Description += "\n" + matches[1]
		}

	case summaryRegex.MatchString(text):
		matches := summaryRegex.FindStringSubmatch(text)
		g.openapi.Info.Summary = matches[1]

	case tosRegex.MatchString(text):
		matches := tosRegex.FindStringSubmatch(text)
		g.openapi.Info.TermsOfService = matches[1]

	// Contact object
	case contactNameRegex.MatchString(text):
		matches := contactNameRegex.FindStringSubmatch(text)
		if g.openapi.Info.Contact == nil {
			g.openapi.Info.Contact = &openapi.Contact{}
		}
		g.openapi.Info.Contact.Name = matches[1]

	case contactURLRegex.MatchString(text):
		matches := contactURLRegex.FindStringSubmatch(text)
		if g.openapi.Info.Contact == nil {
			g.openapi.Info.Contact = &openapi.Contact{}
		}
		g.openapi.Info.Contact.URL = matches[1]

	case contactEmailRegex.MatchString(text):
		matches := contactEmailRegex.FindStringSubmatch(text)
		if g.openapi.Info.Contact == nil {
			g.openapi.Info.Contact = &openapi.Contact{}
		}
		g.openapi.Info.Contact.Email = matches[1]

	// License object
	case licenseNameRegex.MatchString(text):
		matches := licenseNameRegex.FindStringSubmatch(text)
		if g.openapi.Info.License == nil {
			g.openapi.Info.License = &openapi.License{}
		}
		g.openapi.Info.License.Name = matches[1]

	case licenseURLRegex.MatchString(text):
		matches := licenseURLRegex.FindStringSubmatch(text)
		if g.openapi.Info.License == nil {
			g.openapi.Info.License = &openapi.License{}
		}
		g.openapi.Info.License.URL = matches[1]

	case licenseIDRegex.MatchString(text):
		matches := licenseIDRegex.FindStringSubmatch(text)
		if g.openapi.Info.License == nil {
			g.openapi.Info.License = &openapi.License{}
		}
		g.openapi.Info.License.Identifier = matches[1]

	// Server object
	case hostRegex.MatchString(text):
		matches := hostRegex.FindStringSubmatch(text)
		// @host defines the base server URL
		// If no servers exist yet, create one with the host
		if len(g.openapi.Servers) == 0 {
			g.openapi.Servers = append(g.openapi.Servers, openapi.Server{
				URL: "https://" + matches[1],
			})
		} else {
			// Update the first server's URL
			g.openapi.Servers[0].URL = "https://" + matches[1]
		}

	case basePathRegex.MatchString(text):
		matches := basePathRegex.FindStringSubmatch(text)
		// Append basePath to existing server URL
		if len(g.openapi.Servers) > 0 {
			g.openapi.Servers[0].URL += matches[1]
		}

	case schemesRegex.MatchString(text):
		// @schemes is handled differently in OpenAPI 3.x
		// Schemes are part of the server URL
		matches := schemesRegex.FindStringSubmatch(text)
		schemes := strings.Fields(matches[1])
		if len(schemes) > 0 && len(g.openapi.Servers) > 0 {
			// Use the first scheme to update server URL scheme
			scheme := schemes[0]
			if strings.HasPrefix(g.openapi.Servers[0].URL, "https://") {
				g.openapi.Servers[0].URL = strings.Replace(g.openapi.Servers[0].URL, "https://", scheme+"://", 1)
			} else if strings.HasPrefix(g.openapi.Servers[0].URL, "http://") {
				g.openapi.Servers[0].URL = strings.Replace(g.openapi.Servers[0].URL, "http://", scheme+"://", 1)
			}
		}

	case serverRegex.MatchString(text):
		matches := serverRegex.FindStringSubmatch(text)
		server := openapi.Server{
			URL:         matches[1],
			Description: strings.TrimSpace(matches[2]),
		}
		g.openapi.Servers = append(g.openapi.Servers, server)

	// Tag object
	case tagNameRegex.MatchString(text):
		matches := tagNameRegex.FindStringSubmatch(text)
		tag := openapi.Tag{
			Name: matches[1],
		}
		g.openapi.Tags = append(g.openapi.Tags, tag)
		g.lastTag = &g.openapi.Tags[len(g.openapi.Tags)-1]

	case tagDescRegex.MatchString(text):
		matches := tagDescRegex.FindStringSubmatch(text)
		if g.lastTag != nil {
			g.lastTag.Description = matches[1]
		}

	case tagDocsURLRegex.MatchString(text):
		matches := tagDocsURLRegex.FindStringSubmatch(text)
		if g.lastTag != nil {
			if g.lastTag.ExternalDocs == nil {
				g.lastTag.ExternalDocs = &openapi.ExternalDocs{}
			}
			g.lastTag.ExternalDocs.URL = matches[1]
		}

	case tagDocsDescRegex.MatchString(text):
		matches := tagDocsDescRegex.FindStringSubmatch(text)
		if g.lastTag != nil {
			if g.lastTag.ExternalDocs == nil {
				g.lastTag.ExternalDocs = &openapi.ExternalDocs{}
			}
			g.lastTag.ExternalDocs.Description = matches[1]
		}

	// External documentation
	case externalDocsURLRegex.MatchString(text):
		matches := externalDocsURLRegex.FindStringSubmatch(text)
		if g.openapi.ExternalDocs == nil {
			g.openapi.ExternalDocs = &openapi.ExternalDocs{}
		}
		g.openapi.ExternalDocs.URL = matches[1]

	case externalDocsDescRegex.MatchString(text):
		matches := externalDocsDescRegex.FindStringSubmatch(text)
		if g.openapi.ExternalDocs == nil {
			g.openapi.ExternalDocs = &openapi.ExternalDocs{}
		}
		g.openapi.ExternalDocs.Description = matches[1]

	// Security definitions
	case securityBasicRegex.MatchString(text):
		matches := securityBasicRegex.FindStringSubmatch(text)
		g.openapi.Components.SecuritySchemes[matches[1]] = &openapi.SecurityScheme{
			Type:        "http",
			Scheme:      "basic",
			Description: strings.TrimSpace(matches[2]),
		}

	case securityAPIKeyRegex.MatchString(text):
		matches := securityAPIKeyRegex.FindStringSubmatch(text)
		g.openapi.Components.SecuritySchemes[matches[1]] = &openapi.SecurityScheme{
			Type:        "apiKey",
			Name:        matches[2],
			In:          matches[3],
			Description: strings.TrimSpace(matches[4]),
		}
	}

	return nil
}
