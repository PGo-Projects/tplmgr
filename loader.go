package tplmgr

import (
	"html/template"


	"github.com/oxtoacart/bpool"
)

type Config struct {
	layoutPath  string
	includePath string
}

var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var config *Config

var mainTmpl = `{{ define "main" }} {{ template "base" . }} {{ end }}`

func SetConfig(layoutPath string, includePath string) {
	config = &Config{layoutPath, includePath}
}

func MustLoad() {
	funcMap := template.FuncMap{}
	MustLoadTemplatesWithFuncs(funcMap)
}

func MustLoadWithFuncs(funcMap template.FuncMap) {
	if config == nil {
		panic("Error: Template config was not loaded")
	}
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles := mustGetLayoutFiles()
	includeFiles := mustGetIncludeFiles()
	mainTemplate := mustGetMain(funcMap)

	for _, file := range includeFiles {
		filename := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[filename], err := mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[filename] = template.Must(templates[filename].Funcs(funcMap).ParseFiles(files...))
	}
	log.Println("Successfully loaded templates!")

	bufpool = bpool.NewBufferPool(64)
	log.Println("Successfully allocated buffers!")
}
