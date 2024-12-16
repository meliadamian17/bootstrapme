package runner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Runner struct {
	ProjectName string
	Preset      PresetInfo
	LogChan     chan string
	wg          sync.WaitGroup
}

type PresetInfo struct {
	Name            string
	Files           []FileInfo
	PostInstallCmds []string
	Variables       map[string]string
}

type FileInfo struct {
	Path    string
	Content string
}

func (r *Runner) RunPresetAsync() {
	go func() {
		var failed bool

		r.log("Starting bootstrapping...")
		r.log(fmt.Sprintf("Creating project directory: %s", r.ProjectName))
		if err := os.MkdirAll(r.ProjectName, 0755); err != nil {
			r.log(fmt.Sprintf("Error creating project dir: %v", err))
			failed = true
		} else {
			// Write files
			if !failed && len(r.Preset.Files) > 0 {
				r.log("Writing project files...")
			}
			for _, f := range r.Preset.Files {
				if failed {
					break
				}
				targetPath := filepath.Join(r.ProjectName, f.Path)
				if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
					r.log(fmt.Sprintf("Error creating dirs for %s: %v", f.Path, err))
					failed = true
					break
				}
				content := substituteVariables(f.Content, r.Preset.Variables)
				if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
					r.log(fmt.Sprintf("Error writing file %s: %v", f.Path, err))
					failed = true
					break
				}
				r.log("Created file: " + f.Path)
			}

			// Run post-install commands
			if !failed && len(r.Preset.PostInstallCmds) > 0 {
				r.log("Running post-install commands...")
			}
			for _, cmdLine := range r.Preset.PostInstallCmds {
				if failed {
					break
				}
				cmdLine = substituteVariables(cmdLine, r.Preset.Variables)
				err := r.runCommand(cmdLine)
				if err != nil {
					failed = true
					break
				}
			}
		}

		r.wg.Wait()

		r.done()

	}()
}

func (r *Runner) runCommand(cmdLine string) error {
	r.log("Running: " + cmdLine)
	cmd := exec.Command("bash", "-c", cmdLine)
	cmd.Dir = r.ProjectName

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		r.log(fmt.Sprintf("Error stdout pipe: %v", err))
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		r.log(fmt.Sprintf("Error stderr pipe: %v", err))
		return err
	}

	if err := cmd.Start(); err != nil {
		r.log(fmt.Sprintf("Error starting command: %v", err))
		return err
	}

	outBytes, cmdErr := r.readCombinedOutput(stdout, stderr)
	if cmdErr != nil {
		r.log(fmt.Sprintf("Command failed: %v", cmdErr))
		if outBytes != "" {
			r.log("Output:\n" + outBytes)
		}
		return cmdErr
	}

	r.log("Command completed: " + cmdLine)
	if outBytes != "" {
		r.log("Command output:\n" + outBytes)
	}
	return nil
}

func (r *Runner) readCombinedOutput(stdout, stderr io.Reader) (string, error) {
	var buf bytes.Buffer

	r.wg.Add(2)

	go r.scanLines(stdout, &buf, "", &r.wg)
	go r.scanLines(stderr, &buf, "ERR: ", &r.wg)

	return buf.String(), nil
}

func (r *Runner) scanLines(rd io.Reader, buf *bytes.Buffer, prefix string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line := prefix + scanner.Text()
		r.log(line)
		buf.WriteString(line + "\n")
	}
}

func substituteVariables(str string, vars map[string]string) string {
	for k, v := range vars {
		placeholder := fmt.Sprintf("{{ %s }}", k)
		str = strings.ReplaceAll(str, placeholder, v)
	}
	return str
}

func (r *Runner) log(line string) {
	r.LogChan <- line
}

func (r *Runner) done() {
	close(r.LogChan)
}
