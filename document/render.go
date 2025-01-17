package document

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
)

type Config struct {
	OutputDir     string
	TemplateDirs  []string
	RenderFormat  Formatter
	UseExports    bool
	ShortLinks    bool
	CaseSensitive bool
}

func Render(docs *Docs, config *Config) error {
	t, err := loadTemplates(config.RenderFormat, config.TemplateDirs...)
	if err != nil {
		return err
	}
	proc := NewProcessor(docs, config.RenderFormat, t, config.UseExports, config.ShortLinks)
	return renderWith(config, proc)
}

func renderWith(config *Config, proc *Processor) error {
	caseSensitiveSystem = config.CaseSensitive
	if err := proc.PrepareDocs(); err != nil {
		return err
	}
	if err := renderPackage(proc.ExportDocs.Decl, []string{config.OutputDir}, proc); err != nil {
		return err
	}
	if err := proc.Formatter.WriteAuxiliary(proc.ExportDocs.Decl, config.OutputDir, proc); err != nil {
		return err
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
	if err := mkDirs(pkgPath); err != nil {
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
	if err := mkDirs(path.Join(newDir...)); err != nil {
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

func loadTemplates(f Formatter, additional ...string) (*template.Template, error) {
	allTemplates, err := findTemplatesFS()
	if err != nil {
		return nil, err
	}
	templ, err := template.New("all").Funcs(template.FuncMap{
		"toLink": f.ToLinkPath,
	}).ParseFS(assets.Templates, allTemplates...)
	if err != nil {
		return nil, err
	}
	for _, dir := range additional {
		moreTemplates, err := findTemplates(dir)
		if err != nil {
			return nil, err
		}
		templ, err = templ.ParseFiles(moreTemplates...)
		if err != nil {
			return nil, err
		}
	}
	return templ, nil
}

func findTemplatesFS() ([]string, error) {
	allTemplates := []string{}
	err := fs.WalkDir(assets.Templates, ".",
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				allTemplates = append(allTemplates, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}

func findTemplates(dir string) ([]string, error) {
	allTemplates := []string{}
	err := filepath.WalkDir(dir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				allTemplates = append(allTemplates, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}

func linkAndWrite(text string, dir []string, modElems int, kind string, proc *Processor) error {
	text, err := proc.ReplacePlaceholders(text, dir[1:], modElems-1)
	if err != nil {
		return err
	}
	outFile, err := proc.Formatter.ToFilePath(path.Join(dir...), kind)
	if err != nil {
		return err
	}
	return proc.WriteFile(outFile, text)
}

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
