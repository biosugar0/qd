package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

const (
	name                     = "qd"
	helpName                 = "Quick Daily Notes"
	usage                    = "qd [section title]"
	templateDailyNoteContent = `[[daily notes]]

[[{{.Before}}]] | [[{{.After}}]]
`
)

func command() int {
	app := cli.NewApp()
	app.Name = name
	app.HelpName = helpName
	app.Usage = helpName
	app.UsageText = usage
	app.Version = version
	app.Action = run
	return handler(app.Run(os.Args))
}

type config struct {
	DailyNoteDir string `toml:"dailynotedir"`
	Editor       string `toml:"editor"`
}

func (cfg *config) load() error {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "qd")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	file := filepath.Join(dir, "config.toml")

	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		if len(cfg.DailyNoteDir) == 0 {
			return fmt.Errorf("dailynotedir is not found: %v", err)
		}
		cfg.DailyNoteDir = expandPath(cfg.DailyNoteDir)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.Editor = os.Getenv("EDITOR")
	if cfg.Editor == "" {
		cfg.Editor = "vim"
	}

	return toml.NewEncoder(f).Encode(cfg)
}

func run(c *cli.Context) error {
	var cfg config
	err := cfg.load()
	if err != nil {
		return err
	}

	now := time.Now()
	const dateLayout = "2006-01-02"
	dailyNoteName, err := getDailyNoteName(now)
	if err != nil {
		return err
	}

	before := now.AddDate(0, 0, -1).Format(dateLayout)
	after := now.AddDate(0, 0, 1).Format(dateLayout)

	file := filepath.Join(cfg.DailyNoteDir, dailyNoteName)
	var title string
	if c.Args().Present() {
		title = c.Args().First()
	}

	nowString := now.Format("15:04")

	if fileExists(file) {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		if len(title) > 0 {
			fmt.Fprintf(f, "\n### %s %s\n", nowString, title)
			f.Close()
		} else {
			fmt.Fprintf(f, "\n### %s\n", nowString)
			f.Close()
		}
		return openEditor(cfg.Editor, file)
	}

	var tmplString string
	if len(title) > 0 {
		tmplString = fmt.Sprintf("%s\n### %s %s\n", templateDailyNoteContent, nowString, title)
	} else {
		tmplString = fmt.Sprintf("%s\n### %s\n", templateDailyNoteContent, nowString)
	}

	t := template.Must(template.New("qd").Parse(tmplString))
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	err = t.Execute(f, struct {
		Before, After string
	}{
		before, after,
	})

	if err != nil {
		return err
	}
	f.Close()

	return openEditor(cfg.Editor, file)
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		s = filepath.Join(os.Getenv("HOME"), s[2:])
	}
	return os.Expand(s, os.Getenv)
}

func openEditor(command, file string) error {
	cmdargs := file
	command += " " + cmdargs

	cmd := exec.Command("sh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func handler(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func getDailyNoteName(now time.Time) (string, error) {
	const dateLayout = "2006-01-02"
	fname := fmt.Sprintf("%s.md", now.Format(dateLayout))
	return fname, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {
	os.Exit(command())
}
