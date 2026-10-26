// Check plugin metadata and the bundled skill's links and Go example syntax.
package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
)

var linkPattern = regexp.MustCompile(`\]\(([^)]+)\)`)
var codePattern = regexp.MustCompile("(?s)```go\\r?\\n(.*?)\\r?\\n```")
var importsPattern = regexp.MustCompile(`(?s)import\s*\(.*?\)`)

func checkPlugin() error {
	var manifest struct {
		Schema      string `json:"$schema"`
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Extensions  struct {
			OpenAI struct {
				RequiresLocalExecutor bool `json:"requires_local_executor"`
			} `json:"com.openai"`
		} `json:"extensions"`
	}
	data, err := os.ReadFile("plugin.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &manifest); err != nil {
		return err
	}
	if manifest.Schema != "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json" ||
		manifest.Name != "inspire-go-sdk" || manifest.Version != openapi.Version ||
		strings.TrimSpace(manifest.Description) == "" || !manifest.Extensions.OpenAI.RequiresLocalExecutor {
		return fmt.Errorf("invalid plugin identity, version or local executor setting")
	}
	var marketplace struct {
		Name    string `json:"name"`
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"source"`
		} `json:"plugins"`
	}
	data, err = os.ReadFile(filepath.Join(".agents", "plugins", "marketplace.json"))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &marketplace); err != nil {
		return err
	}
	if marketplace.Name != manifest.Name || len(marketplace.Plugins) != 1 {
		return fmt.Errorf("invalid marketplace name or plugin count")
	}
	entry := marketplace.Plugins[0]
	if entry.Name != manifest.Name || entry.Source.Source != "local" || entry.Source.Path == "" || filepath.Clean(entry.Source.Path) != "." {
		return fmt.Errorf("marketplace must resolve the plugin at the repository root")
	}
	fmt.Printf("Plugin valid: %s %s.\n", manifest.Name, manifest.Version)
	return nil
}

func check() error {
	root := filepath.Join("skills", "inspire-go-sdk")
	entry, err := os.ReadFile(filepath.Join(root, "SKILL.md"))
	if err != nil {
		return err
	}
	content := strings.ReplaceAll(string(entry), "\r\n", "\n")
	if !strings.HasPrefix(content, "---\n") {
		return fmt.Errorf("missing skill frontmatter")
	}
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid skill frontmatter")
	}
	metadata := map[string]string{}
	for _, line := range strings.Split(parts[1], "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		pair := strings.SplitN(line, ":", 2)
		if len(pair) != 2 {
			return fmt.Errorf("invalid skill metadata")
		}
		metadata[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	if metadata["name"] != "inspire-go-sdk" || metadata["description"] == "" || len([]rune(metadata["description"])) > 1024 {
		return fmt.Errorf("invalid skill name or description")
	}
	files, examples := 0, 0
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("skill cannot contain symlinks")
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		files++
		data, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		text := string(data)
		for _, match := range linkPattern.FindAllStringSubmatch(text, -1) {
			target := strings.SplitN(match[1], "#", 2)[0]
			if target == "" || strings.Contains(target, "://") {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			relative, e := filepath.Rel(root, resolved)
			if e != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				return fmt.Errorf("skill link escapes its directory")
			}
			if _, e = os.Stat(resolved); e != nil {
				return fmt.Errorf("broken skill link: %s", target)
			}
		}
		for _, match := range codePattern.FindAllStringSubmatch(text, -1) {
			examples++
			source := match[1]
			if !strings.HasPrefix(strings.TrimSpace(source), "package ") {
				imports := importsPattern.FindAllString(source, -1)
				source = importsPattern.ReplaceAllString(source, "")
				source = "package example\n" + strings.Join(imports, "\n") + "\nfunc example(){\n" + source + "\n}"
			}
			if _, e := parser.ParseFile(token.NewFileSet(), path, source, parser.AllErrors); e != nil {
				return fmt.Errorf("invalid Go example in %s: %w", path, e)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Skill valid: %d Markdown files, %d Go examples parsed (syntax only).\n", files, examples)
	return nil
}
func main() {
	if err := checkPlugin(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := check(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
