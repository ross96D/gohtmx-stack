package output

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
)

type Output struct {
	BasePath    string
	DefaultHtmx bool
	ProyectName string
}

func (o Output) Build() error {
	o.crateStructure()
	o.goModInit()
	o.cobraCliInit()
	return nil
}

func (o *Output) crateStructure() error {
	err := os.Mkdir(path.Join(o.BasePath, "assets"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "assets", "js"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "assets", "css"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "handlers"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "webserver"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "pages"), fs.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(o.BasePath, "pages", "index.html"))
	if err != nil {
		return err
	}
	f.WriteString(html5)
	f.Close()

	err = o.htmx()
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) goModInit() error {
	cmd := exec.Command("go", "mod", "init", o.ProyectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *Output) cobraCliInit() error {
	cmd := exec.Command("cobra-cli", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *Output) htmx() error {
	f, err := os.Create(path.Join(o.BasePath, "assets", "js", "htmx.min.js"))
	if err != nil {
		return err
	}
	f.WriteString(htmx)
	return f.Close()
}
