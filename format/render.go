package format

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
)

func Render(docs *document.Docs, dir string, rFormat Format, useExports bool, shortLinks bool) error {
	formatter := GetFormatter(rFormat)
	t, err := loadTemplates(formatter)
	if err != nil {
		return err
	}
	proc := document.NewProcessor(docs, formatter, t, useExports, shortLinks)
	if err := proc.PrepareDocs(); err != nil {
		return err
	}
	if err := renderPackage(proc.ExportDocs.Decl, []string{dir}, &proc); err != nil {
		return err
	}
	if err := proc.Formatter.WriteAuxiliary(proc.ExportDocs.Decl, dir, &proc); err != nil {
		return err
	}
	return nil
}

func renderElement(data interface {
	document.Named
	document.Kinded
}, proc *document.Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return proc.Formatter.ProcessMarkdown(data, b.String(), proc)
}

func renderPackage(p *document.Package, dir []string, proc *document.Processor) error {
	newDir := document.AppendNew(dir, p.GetFileName())
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

func renderModule(mod *document.Module, dir []string, proc *document.Processor) error {
	newDir := document.AppendNew(dir, mod.GetFileName())
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
	document.Named
	document.Kinded
}](list []T, dir []string, proc *document.Processor) error {
	for _, elem := range list {
		newDir := document.AppendNew(dir, elem.GetFileName())
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

func loadTemplates(f document.Formatter) (*template.Template, error) {
	allTemplates, err := findTemplates()
	if err != nil {
		return nil, err
	}
	templ, err := template.New("all").Funcs(template.FuncMap{
		"toLink": f.ToLinkPath,
	}).ParseFS(assets.Templates, allTemplates...)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

func findTemplates() ([]string, error) {
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

func linkAndWrite(text string, dir []string, modElems int, kind string, proc *document.Processor) error {
	text, err := proc.ReplacePlaceholders(text, dir[1:], modElems-1)
	if err != nil {
		return err
	}
	outFile, err := proc.Formatter.ToFilePath(path.Join(dir...), kind)
	if err != nil {
		return err
	}
	return os.WriteFile(outFile, []byte(text), 0666)
}

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
