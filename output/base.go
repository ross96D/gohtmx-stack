package output

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ross96D/gohtmx-stack/output/templates"
)

type Output struct {
	BasePath    string
	ProyectName string
	Htmx        []byte
}

func (o Output) Build() (err error) {
	chanCheckTempl := make(chan any)
	chanCheckCobra := make(chan any)
	chanBuildFirst := make(chan any)

	go o.checkTempl(chanCheckTempl)
	go o.checkCobra(chanCheckCobra)
	go o.buildFirst(chanBuildFirst)

	v := <-chanCheckTempl
	if err, ok := v.(error); ok {
		return err
	}
	v = <-chanCheckCobra
	if err, ok := v.(error); ok {
		return err
	}
	v = <-chanBuildFirst
	if err, ok := v.(error); ok {
		return err
	}

	if err = o.cobraCliInit(); err != nil {
		return
	}
	if err = o.goModTidy(); err != nil {
		return
	}

	if err = o.templateGenerate(); err != nil {
		return
	}
	if err = o.addServeCommand(); err != nil {
		return
	}
	return nil
}

func (o Output) addServeCommand() (err error) {
	f, err := os.Open(path.Join(o.BasePath, "cmd", "root.go"))
	if err != nil {
		return err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	s := string(b)
	s = strings.ReplaceAll(
		s,
		`rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")`,
		"rootCmd.Flags().BoolP(\"toggle\", \"t\", false, \"Help message for toggle\")\nrootCmd.AddCommand(serveCommand)",
	)
	f.Close()
	f, err = os.Create(path.Join(o.BasePath, "cmd", "root.go"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	return err
}

func (o Output) checkTempl(d chan any) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("installing templ %w", err)
			d <- err
		} else {
			d <- true
		}
	}()
	cmd := exec.Command("templ", "")
	err = cmd.Run()
	if err != nil {
		err = o.installTempl()
	}
}

func (o Output) checkCobra(d chan any) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("installing cobra %w", err)
			d <- err
		} else {
			d <- true
		}
	}()
	cmd := exec.Command("cobra-cli", "")
	err = cmd.Run()
	if err != nil {
		err = o.installCobra()
	}
}

func (o Output) installTempl() error {
	println("Installing templ")
	cmd := exec.Command("go", "install", "github.com/a-h/templ/cmd/templ@latest")
	return cmd.Run()
}

func (o Output) installCobra() error {
	println("Installing cobra-cli")
	cmd := exec.Command("go", "install", "github.com/spf13/cobra-cli@latest")
	return cmd.Run()
}

func (o Output) buildFirst(d chan any) {
	var err error
	defer func() {
		if err != nil {
			d <- err
		} else {
			d <- true
		}
	}()
	if err = o.crateStructure(); err != nil {
		return
	}
	if err = o.goModInit(); err != nil {
		return
	}
}

func (o *Output) crateStructure() error {
	err := os.MkdirAll(path.Join(o.BasePath, "assets", "js"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Join(o.BasePath, "assets", "css"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "handlers"), fs.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(o.BasePath, "shared"), fs.ModePerm)
	if err != nil {
		return err
	}

	err = o.writeIndex()
	if err != nil {
		return err
	}

	err = o.htmx()
	if err != nil {
		return err
	}

	err = o.writeServer()
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) templateGenerate() error {
	cmd := exec.Command("templ", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *Output) goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	defer f.Close()
	_, err = f.Write(o.Htmx)
	return err
}

func (o *Output) writeIndex() error {
	err := os.Mkdir(path.Join(o.BasePath, "views"), fs.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(o.BasePath, "views", "index.templ"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf(templates.Templ, o.ProyectName))
	if err != nil {
		return err
	}
	return err
}

func (o *Output) writeServer() (err error) {
	var f *os.File
	defer func() {
		f.Close()
		if err != nil {
			err = fmt.Errorf("write-server %w", err)
		}
	}()
	err = os.MkdirAll(path.Join(o.BasePath, "cmd"), fs.ModePerm)
	if err != nil {
		return err
	}
	f, err = os.Create(path.Join(o.BasePath, "cmd", "serve.go"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(templates.Serve, o.ProyectName))
	if err != nil {
		return err
	}
	f.Close()

	f, err = os.Create(path.Join(o.BasePath, "handlers", "server.go"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(templates.BaseServer)
	if err != nil {
		return err
	}
	f.Close()

	f, err = os.Create(path.Join(o.BasePath, "handlers", "index.go"))
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf(templates.Echo, o.ProyectName, o.ProyectName))
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
