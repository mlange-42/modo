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
	proc := document.NewProcessor(formatter, t, useExports, shortLinks)
	err = proc.ProcessLinks(docs)
	if err != nil {
		return err
	}

	err = renderPackage(docs.Decl, dir, &proc)
	if err != nil {
		return err
	}

	if err := proc.Formatter.WriteAuxiliary(docs.Decl, dir, t); err != nil {
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
	var summary string
	if d, ok := data.(document.Summarized); ok {
		summary = d.GetSummary()
	}
	return proc.Formatter.ProcessMarkdown(data.GetName(), summary, b.String())
}

func renderPackage(p *document.Package, dir string, proc *document.Processor) error {
	if proc.UseExports {
		return renderPackageExports(p, dir, proc, &members{Members: []member{{Include: true}}})
	} else {
		return renderPackageNoExports(p, dir, proc)
	}
}

type members struct {
	Members []member
}

type member struct {
	Include bool
	RelPath []string
}

func renderPackageExports(p *document.Package, dir string, proc *document.Processor, parentMembers *members) error {
	selfIncluded, toCrawl := collectExportMembers(parentMembers)
	for _, ex := range p.Exports {
		var newMember member
		if len(ex.Short) == 1 {
			newMember = member{Include: true}
		} else {
			newMember = member{RelPath: ex.Short[1:]}
		}
		if members, ok := toCrawl[ex.Short[0]]; ok {
			members.Members = append(members.Members, newMember)
			continue
		}
		toCrawl[ex.Short[0]] = &members{Members: []member{newMember}}
	}

	var pkgPath, pkgFile string
	if selfIncluded {
		pkgPath = path.Join(dir, p.GetFileName())
		if err := mkDirs(pkgPath); err != nil {
			return err
		}
		var err error
		pkgFile, err = proc.Formatter.ToFilePath(pkgPath, "package")
		if err != nil {
			return err
		}
	} else {
		pkgPath = dir
	}

	for _, pkg := range p.Packages {
		if mems, ok := toCrawl[pkg.Name]; ok {
			if err := renderPackageExports(pkg, pkgPath, proc, mems); err != nil {
				return err
			}
		}
	}

	for _, mod := range p.Modules {
		if mems, ok := toCrawl[mod.Name]; ok {
			if err := renderModuleExports(mod, pkgPath, proc, mems); err != nil {
				return err
			}
		}
	}

	if selfIncluded {
		text, err := renderElement(p, proc)
		if err != nil {
			return err
		}
		if err := os.WriteFile(pkgFile, []byte(text), 0666); err != nil {
			return err
		}
	}

	return nil
}

func collectExportMembers(parentMembers *members) (bool, map[string]*members) {
	selfIncluded := false
	toCrawl := map[string]*members{}
	for _, mem := range parentMembers.Members {
		if mem.Include {
			selfIncluded = true
			continue
		}
		var newMember member
		if len(mem.RelPath) == 1 {
			newMember = member{Include: true}
		} else {
			newMember = member{RelPath: mem.RelPath[1:]}
		}
		if members, ok := toCrawl[mem.RelPath[0]]; ok {
			members.Members = append(members.Members, newMember)
			continue
		}
		toCrawl[mem.RelPath[0]] = &members{Members: []member{newMember}}
	}
	return selfIncluded, toCrawl
}

func renderPackageNoExports(p *document.Package, dir string, proc *document.Processor) error {
	pkgPath := path.Join(dir, p.GetFileName())
	if err := mkDirs(pkgPath); err != nil {
		return err
	}

	pkgFile, err := proc.Formatter.ToFilePath(pkgPath, "package")
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := renderPackageNoExports(pkg, pkgPath, proc); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		modPath := path.Join(pkgPath, mod.GetFileName())
		if err := renderModuleNoExports(mod, modPath, proc); err != nil {
			return err
		}
	}

	text, err := renderElement(p, proc)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pkgFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderModuleExports(mod *document.Module, dir string, proc *document.Processor, parentMembers *members) error {
	selfIncluded, toCrawl := collectExportMembers(parentMembers)

	var modPath, modFile string
	if selfIncluded {
		modPath = path.Join(dir, mod.GetFileName())
		if err := mkDirs(modPath); err != nil {
			return err
		}
		var err error
		modFile, err = proc.Formatter.ToFilePath(modPath, "module")
		if err != nil {
			return err
		}
	} else {
		modPath = dir
	}

	if err := renderListExports(mod.Structs, modPath, proc, toCrawl); err != nil {
		return err
	}
	if err := renderListExports(mod.Traits, modPath, proc, toCrawl); err != nil {
		return err
	}
	if err := renderListExports(mod.Functions, modPath, proc, toCrawl); err != nil {
		return err
	}

	if selfIncluded {
		text, err := renderElement(mod, proc)
		if err != nil {
			return err
		}
		if err := os.WriteFile(modFile, []byte(text), 0666); err != nil {
			return err
		}
	}

	return nil
}

func renderModuleNoExports(mod *document.Module, dir string, proc *document.Processor) error {
	modPath := path.Join(dir, mod.GetFileName())
	if err := mkDirs(modPath); err != nil {
		return err
	}

	modFile, err := proc.Formatter.ToFilePath(modPath, "module")
	if err != nil {
		return err
	}

	if err := renderListNoExports(mod.Structs, modPath, proc); err != nil {
		return err
	}
	if err := renderListNoExports(mod.Traits, modPath, proc); err != nil {
		return err
	}
	if err := renderListNoExports(mod.Functions, modPath, proc); err != nil {
		return err
	}

	text, err := renderElement(mod, proc)
	if err != nil {
		return err
	}
	if err := os.WriteFile(modFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderListExports[T interface {
	document.Named
	document.Kinded
}](list []T, dir string, proc *document.Processor, parentMembers map[string]*members) error {
	for _, elem := range list {
		if _, ok := parentMembers[elem.GetName()]; !ok {
			continue
		}
		text, err := renderElement(elem, proc)
		if err != nil {
			return err
		}
		memberPath := path.Join(dir, elem.GetFileName())

		memberFile, err := proc.Formatter.ToFilePath(memberPath, elem.GetKind())
		if err != nil {
			return err
		}

		if err := os.WriteFile(memberFile, []byte(text), 0666); err != nil {
			return err
		}
	}
	return nil
}

func renderListNoExports[T interface {
	document.Named
	document.Kinded
}](list []T, dir string, proc *document.Processor) error {
	for _, elem := range list {
		text, err := renderElement(elem, proc)
		if err != nil {
			return err
		}
		memberPath := path.Join(dir, elem.GetFileName())

		memberFile, err := proc.Formatter.ToFilePath(memberPath, elem.GetKind())
		if err != nil {
			return err
		}

		if err := os.WriteFile(memberFile, []byte(text), 0666); err != nil {
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

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
