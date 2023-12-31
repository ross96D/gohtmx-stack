package output

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

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
	chanCheckAir := make(chan any)
	chanBuildFirst := make(chan any)

	go o.checkTempl(chanCheckTempl)
	go o.checkCobra(chanCheckCobra)
	go o.checkAir(chanCheckAir)
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

	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		defer w.Done()
		if err = o.cobraCliInit(); err != nil {
			return
		}
		if err = o.addServeCommand(); err != nil {
			return
		}
		if err = o.goModTidy(); err != nil {
			return
		}
	}()

	if err = o.templateGenerate(); err != nil {
		return
	}
	if err = o.airInit(); err != nil {
		return
	}
	if err = o.gitInit(); err != nil {
		return
	}
	w.Wait()
	return nil
}

func (o Output) addServeCommand() (err error) {
	f1, err := os.Open(path.Join(o.BasePath, "cmd", "root.go"))
	if err != nil {
		return err
	}
	defer f1.Close()
	b, err := io.ReadAll(f1)
	if err != nil {
		return err
	}
	s := string(b)
	s = strings.ReplaceAll(
		s,
		`rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")`,
		"rootCmd.Flags().BoolP(\"toggle\", \"t\", false, \"Help message for toggle\")\n\trootCmd.AddCommand(serveCommand)",
	)
	f2, err := os.Create(path.Join(o.BasePath, "cmd", "root.go"))
	if err != nil {
		return err
	}
	defer f2.Close()
	_, err = f2.WriteString(s)
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

func (o Output) checkAir(d chan any) {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("installing Air %w", err)
			d <- err
		} else {
			d <- true
		}
	}()
	cmd := exec.Command("air", "-v")
	err = cmd.Run()
	if err != nil {
		err = o.installAir()
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

func (o Output) installAir() error {
	println("Installing air")
	cmd := exec.Command("go", "install", "github.com/cosmtrek/air@latest")
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

func (o *Output) crateStructure() (err error) {
	if err = os.MkdirAll(path.Join(o.BasePath, "assets", "js"), fs.ModePerm); err != nil {
		return
	}
	if err = os.MkdirAll(path.Join(o.BasePath, "assets", "css"), fs.ModePerm); err != nil {
		return
	}
	if err = os.Mkdir(path.Join(o.BasePath, "handlers"), fs.ModePerm); err != nil {
		return
	}
	if err = os.Mkdir(path.Join(o.BasePath, "shared"), fs.ModePerm); err != nil {
		return
	}

	if err = o.writeIndex(); err != nil {
		return
	}

	if err = o.htmx(); err != nil {
		return
	}

	if err = o.writeServer(); err != nil {
		return
	}

	if err = o.writePackageJson(); err != nil {
		return
	}

	if err = o.writeTailwind(); err != nil {
		return
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
	if err := os.Mkdir(path.Join(o.BasePath, "views"), fs.ModePerm); err != nil {
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
	defer func() {
		if err != nil {
			err = fmt.Errorf("write-server %w", err)
		}
	}()
	err = os.MkdirAll(path.Join(o.BasePath, "cmd"), fs.ModePerm)
	if err != nil {
		return err
	}
	f1, err := os.Create(path.Join(o.BasePath, "cmd", "serve.go"))
	if err != nil {
		return err
	}
	defer f1.Close()
	_, err = f1.WriteString(fmt.Sprintf(templates.Serve, o.ProyectName))
	if err != nil {
		return err
	}

	f2, err := os.Create(path.Join(o.BasePath, "handlers", "server.go"))
	if err != nil {
		return err
	}
	defer f2.Close()
	_, err = f2.WriteString(templates.BaseServer)
	if err != nil {
		return err
	}

	f3, err := os.Create(path.Join(o.BasePath, "handlers", "index.go"))
	if err != nil {
		return err
	}
	defer f3.Close()
	_, err = f3.WriteString(fmt.Sprintf(templates.Echo, o.ProyectName, o.ProyectName))
	if err != nil {
		return err
	}
	return nil
}

func (o Output) writePackageJson() (err error) {
	f, err := os.Create(path.Join(o.BasePath, "package.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(templates.PackageJson)
	return err
}

func (o Output) writeTailwind() (err error) {
	f1, err := os.Create(path.Join(o.BasePath, "tailwind.config.js"))
	if err != nil {
		return err
	}
	defer f1.Close()

	_, err = f1.WriteString(templates.TailwindConfig)
	if err != nil {
		return err
	}

	f2, err := os.Create(path.Join(o.BasePath, "views", "input.css"))
	if err != nil {
		return err
	}
	defer f2.Close()
	_, err = f2.WriteString(templates.TailwindInput)
	if err != nil {
		return err
	}

	return err
}

func (o Output) airInit() error {
	return exec.Command("air", "init").Run()
}

func (o Output) gitInit() error {
	if err := o.writeGitIgnore(); err != nil {
		return err
	}
	return exec.Command("git", "init").Run()
}

func (o Output) writeGitIgnore() error {
	f1, err := os.Create(path.Join(o.BasePath, ".gitignore"))
	if err != nil {
		return err
	}
	defer f1.Close()
	_, err = f1.WriteString(templates.Gitignore)
	return err
}
