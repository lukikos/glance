package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Pages    []PageConfig   `yaml:"pages"`
	Branding BrandingConfig `yaml:"branding"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host       string        `yaml:"host"`
	Port       uint16        `yaml:"port"`
	AssetsPath string        `yaml:"assets-path"`
	BaseURL    string        `yaml:"base-url"`
	Timeout    time.Duration `yaml:"timeout"`
}

// BrandingConfig holds UI branding/customization settings.
type BrandingConfig struct {
	HideFooter     bool   `yaml:"hide-footer"`
	CustomFooter   string `yaml:"custom-footer"`
	LogoURL        string `yaml:"logo-url"`
	FaviconURL     string `yaml:"favicon-url"`
	SiteName       string `yaml:"site-name"`
	CustomCSS      string `yaml:"custom-css"`
}

// PageConfig represents a single dashboard page.
type PageConfig struct {
	Name    string         `yaml:"name"`
	Slug    string         `yaml:"slug"`
	Columns []ColumnConfig `yaml:"columns"`
	HideDesktopNavigation bool `yaml:"hide-desktop-navigation"`
	ExpandMobilePageNavigation bool `yaml:"expand-mobile-page-navigation"`
}

// ColumnConfig represents a column within a page.
type ColumnConfig struct {
	Size    string        `yaml:"size"`
	Widgets []WidgetConfig `yaml:"widgets"`
}

// WidgetConfig holds configuration for an individual widget.
// The Type field determines which widget implementation is used.
type WidgetConfig struct {
	Type  string `yaml:"type"`
	Title string `yaml:"title"`
	// Additional fields are parsed dynamically per widget type.
	Properties map[string]interface{} `yaml:",inline"`
}

// defaultServerConfig returns a ServerConfig populated with sensible defaults.
func defaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:    "0.0.0.0",
		Port:    8080,
		Timeout: 30 * time.Second,
	}
}

// Load reads and parses the YAML configuration file at the given path.
// Missing optional fields are filled with defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := &Config{
		Server: defaultServerConfig(),
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate performs basic sanity checks on the parsed configuration.
func (c *Config) validate() error {
	if c.Server.Port == 0 {
		return fmt.Errorf("server.port must be a non-zero value")
	}

	slugs := make(map[string]struct{}, len(c.Pages))
	for i, page := range c.Pages {
		if page.Name == "" {
			return fmt.Errorf("pages[%d]: name is required", i)
		}
		slug := page.Slug
		if slug == "" {
			slug = page.Name
		}
		if _, dup := slugs[slug]; dup {
			return fmt.Errorf("pages[%d]: duplicate slug %q", i, slug)
		}
		slugs[slug] = struct{}{}
	}

	return nil
}
