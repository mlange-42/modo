package document

import (
	"fmt"
	"path"
	"strings"
)

// Render generates documentation for the given docs and writes it to the output directory.
func Render(docs *Docs, config *Config, form Formatter, subdir string) error {
	t, err := LoadTemplates(form, config.SourceURLs[strings.ToLower(docs.Decl.Name)], config.TemplateDirs...)
	if err != nil {
		return err
	}
	if !config.DryRun {
		proc := NewProcessor(docs, form, t, config)
		return renderWith(config, proc, subdir)
	}

	files := []string{}
	proc := NewProcessorWithWriter(docs, form, t, config, func(file, text string) error {
		files = append(files, file)
		return nil
	})
	err = renderWith(config, proc, subdir)
	if err != nil {
		return err
	}

	fmt.Println("Dry-run. Would write these files:")
	for _, f := range files {
		fmt.Println(f)
	}
	return nil
}

// ExtractTests extracts tests from the documentation.
func ExtractTests(docs *Docs, config *Config, form Formatter, subdir string) error {
	caseSensitiveSystem = !config.CaseInsensitive
	t, err := LoadTemplates(form, config.SourceURLs[strings.ToLower(docs.Decl.Name)], config.TemplateDirs...)
	if err != nil {
		return err
	}
	var proc *Processor
	if config.DryRun {
		proc = NewProcessorWithWriter(docs, form, t, config, func(file, text string) error {
			return nil
		})
	} else {
		proc = NewProcessor(docs, form, t, config)
	}
	return proc.ExtractTests(subdir)
}

// ExtractTestsMarkdown extracts tests from markdown files.
func ExtractTestsMarkdown(config *Config, form Formatter, baseDir string, build bool) error {
	caseSensitiveSystem = !config.CaseInsensitive

	t, err := LoadTemplates(form, "", config.TemplateDirs...)
	if err != nil {
		return err
	}
	var proc *Processor
	if config.DryRun {
		proc = NewProcessorWithWriter(nil, form, t, config, func(file, text string) error {
			return nil
		})
	} else {
		proc = NewProcessor(nil, form, t, config)
	}
	return proc.extractDocTestsMarkdown(baseDir, build)
}

func renderWith(config *Config, proc *Processor, subdir string) error {
	caseSensitiveSystem = !config.CaseInsensitive

	if err := proc.PrepareDocs(subdir); err != nil {
		return err
	}
	var missing []missingDocs
	var stats missingStats
	if config.ReportMissing {
		missing = proc.Docs.Decl.checkMissing("", &stats)
	}

	outPath := path.Join(config.OutputDir, subdir)
	if err := renderPackage(proc.ExportDocs.Decl, []string{outPath}, proc); err != nil {
		return err
	}
	if err := proc.Formatter.WriteAuxiliary(proc.ExportDocs.Decl, outPath, proc); err != nil {
		return err
	}
	if config.ReportMissing {
		if err := reportMissing(proc.Docs.Decl.Name, missing, stats, config.Strict); err != nil {
			return err
		}
	}
	return nil
}

func renderElement(data interface {
	Named
	Kinded
}, proc *Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return proc.Formatter.ProcessMarkdown(data, b.String(), proc)
}

func renderPackage(p *Package, dir []string, proc *Processor) error {
	newDir := appendNew(dir, p.GetFileName())
	pkgPath := path.Join(newDir...)
	if err := proc.mkDirs(pkgPath); err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := renderPackage(pkg, newDir, proc); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		if err := renderModule(mod, newDir, proc); err != nil {
			return err
		}
	}

	if err := renderList(p.Structs, newDir, proc); err != nil {
		return err
	}
	if err := renderList(p.Traits, newDir, proc); err != nil {
		return err
	}
	if err := renderList(p.Functions, newDir, proc); err != nil {
		return err
	}

	text, err := renderElement(p, proc)
	if err != nil {
		return err
	}
	if err := linkAndWrite(text, newDir, len(newDir), "package", proc); err != nil {
		return err
	}

	return nil
}

func renderModule(mod *Module, dir []string, proc *Processor) error {
	newDir := appendNew(dir, mod.GetFileName())
	if err := proc.mkDirs(path.Join(newDir...)); err != nil {
		return err
	}

	if err := renderList(mod.Structs, newDir, proc); err != nil {
		return err
	}
	if err := renderList(mod.Traits, newDir, proc); err != nil {
		return err
	}
	if err := renderList(mod.Functions, newDir, proc); err != nil {
		return err
	}

	text, err := renderElement(mod, proc)
	if err != nil {
		return err
	}
	if err := linkAndWrite(text, newDir, len(newDir), "module", proc); err != nil {
		return err
	}

	return nil
}

func renderList[T interface {
	Named
	Kinded
}](list []T, dir []string, proc *Processor) error {
	for _, elem := range list {
		newDir := appendNew(dir, elem.GetFileName())
		text, err := renderElement(elem, proc)
		if err != nil {
			return err
		}
		if err := linkAndWrite(text, newDir, len(dir), elem.GetKind(), proc); err != nil {
			return err
		}
	}
	return nil
}

func linkAndWrite(text string, dir []string, modElems int, kind string, proc *Processor) error {
	text, err := proc.ReplacePlaceholders(text, dir[1:], modElems-1)
	if err != nil {
		return err
	}
	outFile := proc.Formatter.ToFilePath(path.Join(dir...), kind)
	return proc.writeFile(outFile, text)
}

func reportMissing(pkg string, missing []missingDocs, stats missingStats, strict bool) error {
	if len(missing) == 0 {
		fmt.Printf("Docstring coverage of package %s: 100%%\n", pkg)
		return nil
	}
	for _, m := range missing {
		fmt.Printf("WARNING: missing %s in %s\n", m.What, m.Who)
	}
	fmt.Printf("Docstring coverage package %s: %.1f%%\n", pkg, 100.0*float64(stats.Total-stats.Missing)/float64(stats.Total))
	if strict {
		return fmt.Errorf("missing docstrings in strict mode")
	}
	return nil
}
