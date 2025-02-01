package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlange-42/modo/document"
	"github.com/rjeczalik/notify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFile = "modo"
const setExitOnError = "set -e"

const initFileText = "__init__.mojo"
const initFileEmoji = "__init__.🔥"

var watchExtensions = []string{".md", ".mojo", ".🔥"}

func runCommand(command string) error {
	commandWithExit := fmt.Sprintf("%s\n%s", setExitOnError, command)
	cmd := exec.Command("bash", "-c", commandWithExit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommands(commands []string) error {
	for _, command := range commands {
		err := runCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func readDocs(file string) (*document.Docs, error) {
	data, err := read(file)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		return document.FromYaml(data)
	}

	return document.FromJson(data)
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}

func fileExists(file string) (exists, isDir bool, err error) {
	var s os.FileInfo
	if s, err = os.Stat(file); err == nil {
		exists = true
		isDir = s.IsDir()
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		return
	}
	err = nil
	return
}

func isPackage(dir string) (isPackage bool, err error) {
	pkgFile := path.Join(dir, initFileText)
	initExists, initIsDir, err := fileExists(pkgFile)
	if err != nil {
		return
	}
	if initExists && !initIsDir {
		isPackage = true
		return
	}

	pkgFile = path.Join(dir, initFileEmoji)
	initExists, initIsDir, err = fileExists(pkgFile)
	if err != nil {
		return
	}
	if initExists && !initIsDir {
		isPackage = true
		return
	}

	return
}

func mountProject(v *viper.Viper, paths []string) error {
	withConfig := len(paths) > 0
	p := "."
	if withConfig {
		p = paths[0]
		if err := os.Chdir(p); err != nil {
			return err
		}
	}

	filePath := configFile + ".yaml"
	exists, isDir, err := fileExists(filePath)
	if err != nil {
		return err
	}
	if !exists || isDir {
		if withConfig {
			return fmt.Errorf("no config file '%s' found in path '%s'", filePath, p)
		}
		return nil
	}

	v.SetConfigName(configFile)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if !notFound {
			return err
		}
		if withConfig {
			return err
		}
	}
	return nil
}

type command = func(file string, args *document.Config, form document.Formatter, subdir string, isFile, isDir bool) error

func runFilesOrDir(cmd command, args *document.Config, form document.Formatter) error {
	if form != nil {
		if err := form.Accepts(args.InputFiles); err != nil {
			return err
		}
	}

	if len(args.InputFiles) == 0 || (len(args.InputFiles) == 1 && args.InputFiles[0] == "") {
		if err := cmd("", args, form, "", false, false); err != nil {
			return err
		}
	}

	stats := make([]struct {
		file bool
		dir  bool
	}, 0, len(args.InputFiles))

	for _, file := range args.InputFiles {
		if s, err := os.Stat(file); err == nil {
			if s.IsDir() && len(args.InputFiles) > 1 {
				return fmt.Errorf("only a single directory at a time can be processed")
			}
			stats = append(stats, struct {
				file bool
				dir  bool
			}{!s.IsDir(), s.IsDir()})
		} else {
			return err
		}
	}

	for i, file := range args.InputFiles {
		s := stats[i]
		if err := cmd(file, args, form, "", s.file, s.dir); err != nil {
			return err
		}
	}
	return nil
}

func runDir(baseDir string, args *document.Config, form document.Formatter, runFile command) error {
	baseDir = filepath.Clean(baseDir)

	err := filepath.WalkDir(baseDir,
		func(p string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(p), ".json") {
				return nil
			}
			cleanDir, _ := filepath.Split(path.Clean(p))
			relDir := filepath.Clean(strings.TrimPrefix(cleanDir, baseDir))
			return runFile(p, args, form, relDir, true, false)
		})
	return err
}

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func GetCwdName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return cwd, err
	}
	return filepath.Base(cwd), nil
}

func commandError(commandType string, err error) error {
	return fmt.Errorf("in script %s: %s\nTo skip pre- and post-processing scripts, use flag '--bare'", commandType, err)
}

// bindFlags binds flags to Viper, filtering out the `--watch` flag.
func bindFlags(v *viper.Viper, flags *pflag.FlagSet) error {
	newFlags := pflag.NewFlagSet("root", pflag.ExitOnError)
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "watch" {
			return
		}
		newFlags.AddFlag(f)
	})
	return v.BindPFlags(newFlags)
}

func watchAndRun(args *document.Config, command func(*document.Config) error) error {
	args.RemovePostScripts()

	c := make(chan notify.EventInfo, 32)
	collected := make(chan []notify.EventInfo, 1)

	toWatch, err := getWatchPaths(args)
	if err != nil {
		return err
	}
	for _, w := range toWatch {
		if err := notify.Watch(w, c, notify.All); err != nil {
			log.Fatal(err)
		}
	}
	defer notify.Stop(c)

	fmt.Printf("Watching for changes: %s\n", strings.Join(toWatch, ", "))
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		var events []notify.EventInfo
		for {
			select {
			case evt := <-c:
				events = append(events, evt)
			case <-ticker.C:
				if len(events) > 0 {
					collected <- events
					events = nil
				} else {
					collected <- nil
				}
			}
		}
	}()

	for events := range collected {
		if events == nil {
			continue
		}
		trigger := false
		for _, e := range events {
			for _, ext := range watchExtensions {
				if strings.HasSuffix(e.Path(), ext) {
					trigger = true
					break
				}
			}
		}
		if trigger {
			if err := command(args); err != nil {
				return err
			}
			fmt.Printf("Watching for changes: %s\n", strings.Join(toWatch, ", "))
		}
	}
	return nil
}

func getWatchPaths(args *document.Config) ([]string, error) {
	toWatch := append([]string{}, args.Sources...)
	toWatch = append(toWatch, args.InputFiles...)
	for i, w := range toWatch {
		p := w
		exists, isDir, err := fileExists(p)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("file or directory '%s' to watch does not exist", p)
		}
		if isDir {
			p = path.Join(w, "...")
		}
		toWatch[i] = p
	}
	return toWatch, nil
}
