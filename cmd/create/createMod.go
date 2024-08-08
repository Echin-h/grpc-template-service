package create

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"grpc-template-service/pkg/colorful"
	"grpc-template-service/pkg/fsx"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	daoTemplatePath     = "pkg/template/dao.template"
	eTemplatePath       = "pkg/template/e.template"
	serviceTemplatePath = "pkg/template/service.template"
	modelTemplatePath   = "pkg/template/model.template"
	modTemplatePath     = "pkg/template/mod.template"
)

var (
	modName  string
	dir      string
	force    bool
	StartCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a new mod",
		Example: "go run main.go create -n users",
		Run: func(cmd *cobra.Command, args []string) {
			err := load()
			if err != nil {
				fmt.Println(colorful.Red(err.Error()))
				os.Exit(1)
			}
			fmt.Println(colorful.Green("Module " + modName + " generate success under " + dir))
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&modName, "name", "n", "", "Create a new mod with provided name")
	StartCmd.PersistentFlags().StringVarP(&dir, "path", "p", "internal/mod", "New file will generate under provided path")
	StartCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force generate the mod")
}

func load() error {
	if modName == "" {
		return errors.New("mod name should not be empty, use -n")
	}

	if err := loadTemplate(); err != nil {
		fmt.Println(colorful.Red("loadTemplate error:" + err.Error()))
		return err
	}

	return nil

}

func loadTemplate() error {

	dao := path.Join(dir, modName, "dao")
	e := path.Join(dir, modName, "e")
	service := path.Join(dir, modName, "service")
	model := path.Join(dir, modName, "model")
	mod := path.Join(dir, modName)

	_ = fsx.IsNotExistMkDir(dao)
	_ = fsx.IsNotExistMkDir(e)
	_ = fsx.IsNotExistMkDir(service)
	_ = fsx.IsNotExistMkDir(model)
	_ = fsx.IsNotExistMkDir(mod)

	// data map for template
	data := make(map[string]string)
	data["modName"] = strings.ToLower(modName[:1]) + modName[1:]

	dao += "/" + "dao.go"
	e += "/" + "e.go"
	service += "/" + "service.go"
	model += "/" + "model.go"
	mod += "/" + "mod.go"

	if !force && (fsx.FileExist(dao) || fsx.FileExist(e) || fsx.FileExist(service) || fsx.FileExist(model) || fsx.FileExist(mod)) {
		return errors.New("mod already exists, use -f to force generate")
	}

	if rt, err := template.ParseFiles(daoTemplatePath); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, data)
		fsx.FileCreate(b, dao)
	}

	if rt, err := template.ParseFiles(eTemplatePath); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, data)
		fsx.FileCreate(b, e)
	}

	if rt, err := template.ParseFiles(serviceTemplatePath); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, data)
		fsx.FileCreate(b, service)
	}

	if rt, err := template.ParseFiles(modelTemplatePath); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, data)
		fsx.FileCreate(b, model)
	}

	if rt, err := template.ParseFiles(modTemplatePath); err != nil {
		return err
	} else {
		var b bytes.Buffer
		err = rt.Execute(&b, data)
		fsx.FileCreate(b, mod)
	}

	return nil
}
