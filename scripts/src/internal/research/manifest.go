package research

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Dependency represents a single dependency extracted from a manifest file.
type Dependency struct {
	Name           string
	CurrentVersion string
	Ecosystem      string
}

// ParseManifest parses a supported manifest file and returns its dependencies.
func ParseManifest(path string) ([]Dependency, error) {
	base := filepath.Base(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasSuffix(base, "go.mod"):
		return parseGoMod(string(data)), nil
	case strings.HasSuffix(base, "package.json"):
		return parsePackageJSON(data)
	case strings.HasSuffix(base, "requirements.txt"):
		return parseRequirementsTxt(string(data)), nil
	case strings.HasSuffix(base, "pyproject.toml"):
		return parsePyprojectTOML(string(data)), nil
	case strings.HasSuffix(base, "Cargo.toml"):
		return parseCargoTOML(string(data)), nil
	default:
		return nil, fmt.Errorf("unsupported manifest type: %s", path)
	}
}

func parseGoMod(content string) []Dependency {
	var deps []Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if inRequire {
			if strings.HasPrefix(line, ")") {
				inRequire = false
				continue
			}
			if strings.Contains(line, "// indirect") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				deps = append(deps, Dependency{
					Name:           fields[0],
					CurrentVersion: fields[1],
					Ecosystem:      "go",
				})
			}
			continue
		}

		if strings.HasPrefix(line, "require") {
			if strings.Contains(line, "(") {
				inRequire = true
				continue
			}
			if strings.Contains(line, "// indirect") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				deps = append(deps, Dependency{
					Name:           fields[1],
					CurrentVersion: fields[2],
					Ecosystem:      "go",
				})
			}
		}
	}

	return deps
}

func parsePackageJSON(data []byte) ([]Dependency, error) {
	var pkg struct {
		Dependencies    json.RawMessage `json:"dependencies"`
		DevDependencies json.RawMessage `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []Dependency
	parsed, err := parseNPMObject(pkg.Dependencies)
	if err != nil {
		return nil, err
	}
	deps = append(deps, parsed...)

	parsed, err = parseNPMObject(pkg.DevDependencies)
	if err != nil {
		return nil, err
	}
	deps = append(deps, parsed...)

	return deps, nil
}

func parseNPMObject(raw json.RawMessage) ([]Dependency, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '{' {
		return nil, fmt.Errorf("expected JSON object for npm dependencies")
	}

	var deps []Dependency
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return nil, err
		}
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key in npm dependencies")
		}
		var value string
		if err := dec.Decode(&value); err != nil {
			return nil, err
		}
		deps = append(deps, Dependency{
			Name:           keyStr,
			CurrentVersion: value,
			Ecosystem:      "npm",
		})
	}

	// consume closing }
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return deps, nil
}

func parseRequirementsTxt(content string) []Dependency {
	var deps []Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, version, ok := splitNameVersion(line)
		if ok {
			deps = append(deps, Dependency{
				Name:           name,
				CurrentVersion: version,
				Ecosystem:      "pypi",
			})
		}
	}

	return deps
}

func parsePyprojectTOML(content string) []Dependency {
	var deps []Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	inProject := false
	inDeps := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !inProject {
			if trimmed == "[project]" {
				inProject = true
			}
			continue
		}

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if inDeps {
				break
			}
			continue
		}

		if !inDeps {
			if strings.HasPrefix(trimmed, "dependencies = [") {
				inDeps = true
			}
			continue
		}

		if trimmed == "]" || trimmed == "]," {
			break
		}

		if idx := strings.Index(trimmed, "]"); idx >= 0 {
			trimmed = trimmed[:idx]
			if dep, ok := parsePyprojectDepLine(trimmed); ok {
				deps = append(deps, dep)
			}
			break
		}

		if dep, ok := parsePyprojectDepLine(trimmed); ok {
			deps = append(deps, dep)
		}
	}

	return deps
}

func parsePyprojectDepLine(line string) (Dependency, bool) {
	entry := strings.Trim(line, `", `)
	name, version, ok := splitNameVersion(entry)
	if !ok {
		return Dependency{}, false
	}
	return Dependency{
		Name:           name,
		CurrentVersion: version,
		Ecosystem:      "pypi",
	}, true
}

func parseCargoTOML(content string) []Dependency {
	var deps []Dependency
	scanner := bufio.NewScanner(strings.NewReader(content))
	inDeps := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !inDeps {
			if line == "[dependencies]" {
				inDeps = true
			}
			continue
		}

		if strings.HasPrefix(line, "[") {
			break
		}

		name, version, ok := parseCargoDepLine(line)
		if ok {
			deps = append(deps, Dependency{
				Name:           name,
				CurrentVersion: version,
				Ecosystem:      "crates",
			})
		}
	}

	return deps
}

func parseCargoDepLine(line string) (string, string, bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	name := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if strings.HasPrefix(value, "{") {
		idx := strings.Index(value, "version")
		if idx < 0 {
			return "", "", false
		}
		after := value[idx+len("version"):]
		after = strings.TrimSpace(after)
		if !strings.HasPrefix(after, "=") {
			return "", "", false
		}
		after = strings.TrimSpace(after[1:])
		if !strings.HasPrefix(after, `"`) {
			return "", "", false
		}
		after = after[1:]
		end := strings.Index(after, `"`)
		if end < 0 {
			return "", "", false
		}
		return name, after[:end], true
	}

	version := strings.Trim(value, `"`)
	return name, version, true
}

func splitNameVersion(s string) (string, string, bool) {
	// Try each known operator in order of specificity (longest first).
	operators := []string{"~=", "!=", "<=", ">=", "==", "<", ">"}
	var bestSep string
	var bestIdx = -1

	for _, sep := range operators {
		if idx := strings.Index(s, sep); idx >= 0 {
			if bestIdx < 0 || idx < bestIdx {
				bestSep = sep
				bestIdx = idx
			}
		}
	}

	if bestIdx < 0 {
		return "", "", false
	}

	name := strings.TrimSpace(s[:bestIdx])
	version := strings.TrimSpace(s[bestIdx+len(bestSep):])
	if name == "" || version == "" {
		return "", "", false
	}
	return name, version, true
}
