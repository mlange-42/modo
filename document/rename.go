package document

import "fmt"

func (proc *Processor) renameAll(p *Package, parentPath string) {
	myPath := parentPath + "." + p.Name

	fmt.Println(myPath)

	for _, pp := range p.Packages {
		proc.renameAll(pp, myPath)
	}
}
