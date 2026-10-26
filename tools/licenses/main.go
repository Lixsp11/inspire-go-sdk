// Collect dependency and Go runtime license texts for distribution archives.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type module struct {
	Path, Version, Dir string
	Main               bool
}
type listedPackage struct{ Module *module }
type notice struct {
	Path    string   `json:"path"`
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

func legalName(name string) bool {
	upper := strings.ToUpper(name)
	for _, prefix := range []string{"LICENSE", "LICENCE", "COPYING", "NOTICE", "PATENTS"} {
		if upper == prefix || strings.HasPrefix(upper, prefix+".") {
			return true
		}
	}
	return false
}
func copyLegal(source, destination string) ([]string, error) {
	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !legalName(entry.Name()) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			return nil, err
		}
		if err = os.MkdirAll(destination, 0755); err != nil {
			return nil, err
		}
		if err = os.WriteFile(filepath.Join(destination, entry.Name()), data, 0644); err != nil {
			return nil, err
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	return files, nil
}
func run(output, targetOS, targetArch string) error {
	binary := "go"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	command := exec.Command(filepath.Join(runtime.GOROOT(), "bin", binary), "list", "-deps", "-json", "./src/main")
	command.Env = os.Environ()
	if targetOS != "" {
		command.Env = append(command.Env, "GOOS="+targetOS)
	}
	if targetArch != "" {
		command.Env = append(command.Env, "GOARCH="+targetArch)
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	data, err := command.Output()
	if err != nil {
		return fmt.Errorf("list linked modules: %w: %s", err, stderr.String())
	}
	dependencies := map[string]module{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var item listedPackage
		err = decoder.Decode(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if item.Module != nil && !item.Module.Main {
			dependencies[item.Module.Path] = *item.Module
		}
	}
	if err = os.MkdirAll(output, 0755); err != nil {
		return err
	}
	var records []notice
	keys := make([]string, 0, len(dependencies))
	for key := range dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		item := dependencies[key]
		directory := strings.ReplaceAll(item.Path, "/", "__") + "@" + item.Version
		files, e := copyLegal(item.Dir, filepath.Join(output, directory))
		if e != nil {
			return e
		}
		if len(files) == 0 {
			return fmt.Errorf("dependency has no top-level license: %s", item.Path)
		}
		for i := range files {
			files[i] = filepath.ToSlash(filepath.Join(directory, files[i]))
		}
		records = append(records, notice{item.Path, item.Version, files})
	}
	files, err := copyLegal(runtime.GOROOT(), filepath.Join(output, "go"))
	if err != nil {
		return err
	}
	for i := range files {
		files[i] = filepath.ToSlash(filepath.Join("go", files[i]))
	}
	records = append(records, notice{"Go standard library", runtime.Version(), files})
	vendor := filepath.Join(runtime.GOROOT(), "src", "vendor")
	err = filepath.WalkDir(vendor, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !legalName(entry.Name()) {
			return nil
		}
		relative, e := filepath.Rel(vendor, path)
		if e != nil {
			return e
		}
		destination := filepath.Join(output, "go-vendor", relative)
		if e = os.MkdirAll(filepath.Dir(destination), 0755); e != nil {
			return e
		}
		content, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		return os.WriteFile(destination, content, 0644)
	})
	if err != nil {
		return err
	}
	manifest, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(output, "MODULES.json"), append(manifest, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Collected license texts for %d linked modules and the Go standard library.\n", len(dependencies))
	return nil
}
func main() {
	output := flag.String("output", "third_party/licenses", "Output directory")
	targetOS := flag.String("os", "", "Target GOOS for dependency selection")
	targetArch := flag.String("arch", "", "Target GOARCH for dependency selection")
	flag.Parse()
	if err := run(*output, *targetOS, *targetArch); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
