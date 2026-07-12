// Package config provides injectable configuration for the Spacecraft CLI.
// It replaces the global ROOT/SPACE_DIR vars in the original package main.
package config

import (
	"fmt"
	"path/filepath"
)

// Config holds all path configuration for the Spacecraft mission system.
type Config struct {
	root        string
	spaceDir    string
	missionsDir string
	archiveDir  string
	currentFile string
}

// ConfigOption customizes a Config after creation.
type ConfigOption func(*Config)

// WithSpaceDir overrides the default .space directory path.
func WithSpaceDir(dir string) ConfigOption {
	return func(c *Config) { c.spaceDir = dir }
}

// WithMissionsDir overrides the default missions directory path.
func WithMissionsDir(dir string) ConfigOption {
	return func(c *Config) { c.missionsDir = dir }
}

// WithArchiveDir overrides the default archive directory path.
func WithArchiveDir(dir string) ConfigOption {
	return func(c *Config) { c.archiveDir = dir }
}

// WithCurrentFile overrides the default current file path.
func WithCurrentFile(file string) ConfigOption {
	return func(c *Config) { c.currentFile = file }
}

// NewConfig creates a Config with paths derived from root.
// root must be an absolute directory path.
func NewConfig(root string, opts ...ConfigOption) (*Config, error) {
	if root == "" {
		return nil, fmt.Errorf("config: root directory must not be empty")
	}
	if !filepath.IsAbs(root) {
		return nil, fmt.Errorf("config: root must be absolute: %s", root)
	}
	c := &Config{
		root:        root,
		spaceDir:    filepath.Join(root, ".space"),
		missionsDir: filepath.Join(root, ".space", "missions"),
		archiveDir:  filepath.Join(root, ".space", "archive"),
		currentFile: filepath.Join(root, ".space", "current"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// Root returns the project root directory.
func (c *Config) Root() string { return c.root }

// SpaceDir returns the .space directory path.
func (c *Config) SpaceDir() string { return c.spaceDir }

// MissionsDir returns the missions directory path.
func (c *Config) MissionsDir() string { return c.missionsDir }

// ArchiveDir returns the archive directory path.
func (c *Config) ArchiveDir() string { return c.archiveDir }

// CurrentFile returns the .space/current file path.
func (c *Config) CurrentFile() string { return c.currentFile }

// MissionDir returns the mission directory for the given mission id.
func (c *Config) MissionDir(id string) string {
	return filepath.Join(c.missionsDir, id)
}
