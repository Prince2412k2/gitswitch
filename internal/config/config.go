package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Profile represents a single git identity profile.
type Profile struct {
	Name       string `toml:"name"`
	Email      string `toml:"email"`
	SSHKey     string `toml:"ssh_key"`
	GitHubUser string `toml:"github_user,omitempty"`
	GitHubURL  string `toml:"github_url,omitempty"`
	Notes      string `toml:"notes,omitempty"`
}

// Dir returns the path to ~/.config/git_conf/
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "git_conf"), nil
}

// EnsureDir creates the config dir if it doesn't exist.
func EnsureDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("could not create config dir: %w", err)
	}
	return dir, nil
}

// LoadAll reads all .toml profiles from the config dir.
func LoadAll() ([]Profile, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read config dir: %w", err)
	}

	var profiles []Profile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		p, err := Load(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// Load reads a single profile from a TOML file.
func Load(path string) (Profile, error) {
	var p Profile
	if _, err := toml.DecodeFile(path, &p); err != nil {
		return Profile{}, fmt.Errorf("could not parse profile %s: %w", path, err)
	}
	return p, nil
}

// Save writes a profile to ~/.config/git_conf/<slug>.toml
func Save(p Profile) error {
	dir, err := EnsureDir()
	if err != nil {
		return err
	}

	slug := toSlug(p.GitHubUser)
	if slug == "" {
		slug = toSlug(p.Name)
	}
	path := filepath.Join(dir, slug+".toml")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("could not write profile: %w", err)
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(p)
}

// ProfilePath returns the expected path for a profile.
func ProfilePath(p Profile) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	slug := toSlug(p.GitHubUser)
	if slug == "" {
		slug = toSlug(p.Name)
	}
	return filepath.Join(dir, slug+".toml"), nil
}

func toSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
