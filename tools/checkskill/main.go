// Check plugin metadata and validate skill examples through the CLI without network calls.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	commands "github.com/Lixsp11/inspire-go-sdk/src/openapi"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
)

var linkPattern = regexp.MustCompile(`\]\(([^)]+)\)`)
var codePattern = regexp.MustCompile("(?s)```shell\\r?\\n(.*?)\\r?\\n```")
var bodyPattern = regexp.MustCompile("(?s)```json\\r?\\n(.*?)\\r?\\n```")

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
	version, err := os.ReadFile(filepath.Join("skills", "inspire-go-sdk", "scripts", "VERSION"))
	if err != nil || strings.TrimSpace(string(version)) != manifest.Version {
		return fmt.Errorf("launcher and plugin versions must match")
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
	jobs, err := os.ReadFile(filepath.Join(root, "references", "jobs.md"))
	if err != nil {
		return err
	}
	bodyMatch := bodyPattern.FindStringSubmatch(string(jobs))
	if len(bodyMatch) != 2 {
		return fmt.Errorf("missing job request example")
	}
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
		if strings.Contains(text, "```go") {
			return fmt.Errorf("skill examples must use the prebuilt CLI")
		}
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
			for _, line := range strings.Split(match[1], "\n") {
				args := strings.Fields(line)
				if len(args) == 0 {
					continue
				}
				if args[0] != "inspire" {
					return fmt.Errorf("unexpected CLI in %s", path)
				}
				args = args[1:]
				examples++
				for i := 0; i < len(args)-1; i++ {
					if args[i] == "--body-file" {
						args[i+1] = "-"
					}
				}
				help := len(args) == 0 || args[0] == "configure" || args[0] == "version"
				dry := false
				for _, arg := range args {
					if arg == "--help" {
						help = true
					}
					if arg == "--cli-dry-run-json" {
						dry = true
					}
				}
				if help {
					if !strings.Contains(strings.Join(args, " "), "--help") {
						args = append(args, "--help")
					}
				} else if !dry {
					args = append(args, "--cli-dry-run-json")
				}
				ctx := cli.NewCommandContext(&bytes.Buffer{}, &bytes.Buffer{})
				ctx.Context = context.Background()
				ctx.Stdin = strings.NewReader(bodyMatch[1])
				ctx.Getenv = func(string) string { return "" }
				if e := commands.NewCommand().Execute(ctx, args); e != nil {
					return fmt.Errorf("invalid CLI example in %s: %w", path, e)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if examples == 0 {
		return fmt.Errorf("skill must include CLI examples")
	}
	fmt.Printf("Skill valid: %d Markdown files, %d CLI examples checked without platform access.\n", files, examples)
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
