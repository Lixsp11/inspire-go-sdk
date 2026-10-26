// Check launcher downloads, checksums, cache reuse and CLI I/O without network access.
package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func run(binary string) error {
	data, err := os.ReadFile(binary)
	if err != nil {
		return err
	}
	versionData, err := os.ReadFile("skills/inspire-go-sdk/scripts/VERSION")
	if err != nil {
		return err
	}
	version := strings.TrimSpace(string(versionData))
	root, err := os.MkdirTemp("tmp", "launcher-check-")
	if err != nil {
		return err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	fixtures := filepath.Join(root, "release")
	if err = os.MkdirAll(fixtures, 0700); err != nil {
		return err
	}
	name := "inspire"
	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		name += ".exe"
		extension = "zip"
	}
	asset := fmt.Sprintf("inspire_%s_%s_%s.%s", version, runtime.GOOS, runtime.GOARCH, extension)
	files := map[string][]byte{name: data, "LICENSE": []byte("test fixture"), "NOTICE": []byte("test fixture"), "third_party/test/LICENSE": []byte("test fixture")}
	var archive bytes.Buffer
	if runtime.GOOS == "windows" {
		z := zip.NewWriter(&archive)
		for path, content := range files {
			w, e := z.Create(path)
			if e != nil {
				return e
			}
			if _, e = w.Write(content); e != nil {
				return e
			}
		}
		if err = z.Close(); err != nil {
			return err
		}
	} else {
		g := gzip.NewWriter(&archive)
		t := tar.NewWriter(g)
		for path, content := range files {
			mode := int64(0600)
			if path == name {
				mode = 0700
			}
			if err = t.WriteHeader(&tar.Header{Name: path, Size: int64(len(content)), Mode: mode}); err != nil {
				return err
			}
			if _, err = t.Write(content); err != nil {
				return err
			}
		}
		if err = t.Close(); err != nil {
			return err
		}
		if err = g.Close(); err != nil {
			return err
		}
	}
	if err = os.WriteFile(filepath.Join(fixtures, asset), archive.Bytes(), 0600); err != nil {
		return err
	}
	checksum := fmt.Sprintf("%x  %s\n", sha256.Sum256(archive.Bytes()), asset)
	if err = os.WriteFile(filepath.Join(fixtures, "SHA256SUMS"), []byte(checksum), 0600); err != nil {
		return err
	}
	launcher, _ := filepath.Abs("skills/inspire-go-sdk/scripts/inspire.sh")
	harness := ""
	shell := "sh"
	env := os.Environ()
	env = append(env, "INSPIRE_TEST_RELEASE_DIR="+fixtures)
	if runtime.GOOS == "windows" {
		launcher, _ = filepath.Abs("skills/inspire-go-sdk/scripts/inspire.ps1")
		shell, err = exec.LookPath("powershell.exe")
		if err != nil {
			return err
		}
		harness = filepath.Join(root, "harness.ps1")
		code := `$cliArguments = @($args)
function Invoke-WebRequest {
 param([switch]$UseBasicParsing,[string]$Uri,[string]$OutFile)
 if (-not $Uri.StartsWith('https://github.com/Lixsp11/inspire-go-sdk/releases/download/')) { throw 'Unexpected release host' }
 $assetName = $Uri.Substring($Uri.LastIndexOf('/')+1)
 Copy-Item -LiteralPath (Join-Path $env:INSPIRE_TEST_RELEASE_DIR $assetName) -Destination $OutFile
 Add-Content -LiteralPath (Join-Path $env:INSPIRE_TEST_RELEASE_DIR 'requests') -Value $assetName
}
if ($env:INSPIRE_TEST_PIPELINE -eq '1') {
 '{"workspace_id":"pipeline-ws"}' | & $env:INSPIRE_TEST_LAUNCHER train ListJobs --body-file=- --cli-dry-run-json
} else {
 & $env:INSPIRE_TEST_LAUNCHER @cliArguments
}
exit $LASTEXITCODE
`
		if err = os.WriteFile(harness, []byte(code), 0600); err != nil {
			return err
		}
		env = append(env, "INSPIRE_TEST_LAUNCHER="+launcher)
	} else {
		bin := filepath.Join(root, "bin")
		if err = os.Mkdir(bin, 0700); err != nil {
			return err
		}
		code := `#!/bin/sh
set -eu
url= out=
while [ "$#" -gt 0 ]; do
 case "$1" in -o) shift; out=$1 ;; https://github.com/Lixsp11/inspire-go-sdk/releases/download/*) url=$1 ;; esac
 shift
done
[ -n "$url" ] && [ -n "$out" ]
name=${url##*/}
cp "$INSPIRE_TEST_RELEASE_DIR/$name" "$out"
printf '%s\n' "$name" >> "$INSPIRE_TEST_RELEASE_DIR/requests"
`
		if err = os.WriteFile(filepath.Join(bin, "curl"), []byte(code), 0700); err != nil {
			return err
		}
		env = append(env, "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	invoke := func(cache, input string, args ...string) (int, string, string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var command *exec.Cmd
		if runtime.GOOS == "windows" {
			options := []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", harness}
			command = exec.CommandContext(ctx, shell, append(options, args...)...)
		} else {
			command = exec.CommandContext(ctx, shell, append([]string{launcher}, args...)...)
		}
		command.Env = append(append([]string(nil), env...), "INSPIRE_CLI_CACHE_DIR="+cache)
		command.WaitDelay = 5 * time.Second
		command.Stdin = strings.NewReader(input)
		out, diagnostics := &bytes.Buffer{}, &bytes.Buffer{}
		command.Stdout = out
		command.Stderr = diagnostics
		err := command.Run()
		code := 0
		if err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				code = exit.ExitCode()
			} else {
				return 0, "", "", err
			}
		}
		return code, out.String(), diagnostics.String(), nil
	}
	cache := filepath.Join(root, "cache")
	code, out, diagnostics, err := invoke(cache, "", "version")
	if err != nil || code != 0 || diagnostics != "" || !strings.Contains(out, version) {
		return fmt.Errorf("bootstrap version failed: code=%d stdout=%s stderr=%s err=%v", code, out, diagnostics, err)
	}
	requestsBefore, err := os.ReadFile(filepath.Join(fixtures, "requests"))
	if err != nil {
		return err
	}
	code, out, diagnostics, err = invoke(cache, "", "train", "ListJobs", "--body", `{"workspace_id":"ws","keyword":"quoted \"value\" \\"}`, "--cli-dry-run-json")
	if err != nil || code != 0 || diagnostics != "" || !json.Valid([]byte(out)) {
		return fmt.Errorf("launcher argument forwarding failed: %d %s %s %v", code, out, diagnostics, err)
	}
	code, out, diagnostics, err = invoke(cache, `{"workspace_id":"stdin-ws"}`, "train", "ListJobs", "--body-file=-", "--cli-dry-run-json")
	if err != nil || code != 0 || diagnostics != "" || !strings.Contains(out, "stdin-ws") {
		return fmt.Errorf("launcher stdin forwarding failed: %d %s %s %v", code, out, diagnostics, err)
	}
	code, out, diagnostics, err = invoke(cache, "", "train", "ListJobs", "--WorkspaceId", "ws", "--PageSize", "invalid")
	if err != nil || code != 2 || out != "" || !json.Valid([]byte(diagnostics)) {
		return fmt.Errorf("launcher error forwarding failed: %d %s %s %v", code, out, diagnostics, err)
	}
	requestsAfter, _ := os.ReadFile(filepath.Join(fixtures, "requests"))
	if !bytes.Equal(requestsBefore, requestsAfter) {
		return fmt.Errorf("cache hits downloaded a release again")
	}
	code, out, diagnostics, err = invoke(cache, strings.Repeat("x", 2<<20), "train", "ListJobs", "--PageSize", "invalid", "--body-file=-")
	if err != nil || code != 2 || out != "" || !json.Valid([]byte(diagnostics)) {
		return fmt.Errorf("early CLI rejection damaged diagnostics: %d %s %s %v", code, out, diagnostics, err)
	}
	if runtime.GOOS == "windows" {
		env = append(env, "INSPIRE_TEST_PIPELINE=1")
		code, out, diagnostics, err = invoke(cache, "")
		env = append(env, "INSPIRE_TEST_PIPELINE=")
		if err != nil || code != 0 || diagnostics != "" || !strings.Contains(out, "pipeline-ws") {
			return fmt.Errorf("PowerShell pipeline forwarding failed: %d %s %s %v", code, out, diagnostics, err)
		}
	}
	if err = os.WriteFile(filepath.Join(fixtures, "SHA256SUMS"), []byte(strings.Repeat("0", 64)+"  "+asset+"\n"), 0600); err != nil {
		return err
	}
	badCache := filepath.Join(root, "bad-cache")
	code, out, diagnostics, err = invoke(badCache, "", "version")
	if err != nil || code != 1 || out != "" || !json.Valid([]byte(diagnostics)) || !strings.Contains(diagnostics, "BootstrapError") {
		return fmt.Errorf("bad checksum was not rejected: %d %s %s %v", code, out, diagnostics, err)
	}
	if _, err = os.Stat(filepath.Join(badCache, version, runtime.GOOS+"_"+runtime.GOARCH, name)); !os.IsNotExist(err) {
		return fmt.Errorf("bad checksum installed an executable")
	}
	if err = os.WriteFile(filepath.Join(cache, version, runtime.GOOS+"_"+runtime.GOARCH, name), []byte("invalid executable fixture"), 0700); err != nil {
		return err
	}
	code, out, diagnostics, err = invoke(cache, "", "version")
	if err != nil || code != 1 || out != "" || !json.Valid([]byte(diagnostics)) || !strings.Contains(diagnostics, "BootstrapError") {
		return fmt.Errorf("corrupt cache did not produce safe JSON: %d %s %s %v", code, out, diagnostics, err)
	}
	fmt.Println("Launcher valid: verified download, cache, quoting, stdin, JSON errors, exit status and checksum rejection; no network or installation.")
	return nil
}

func main() {
	binary := flag.String("binary", "", "Native CLI executable to test")
	flag.Parse()
	if *binary == "" {
		fmt.Fprintln(os.Stderr, "--binary is required")
		os.Exit(2)
	}
	if err := run(*binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
