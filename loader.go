package tplmgr

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/oxtoacart/bpool"
)

type Config struct {
	layoutPath  string
	includePath string
}

var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var config *Config

var leftDelim = "{{"
var rightDelim = "}}"

const mainTmpl = `{{ define "main" }} {{ template "base" . }} {{ end }}`

func SetConfig(layoutPath string, includePath string) {
	if !strings.HasSuffix(layoutPath, "/") {
		layoutPath += "/"
	}
	if !strings.HasSuffix(includePath, "/") {
		includePath += "/"
	}
	config = &Config{layoutPath, includePath}
}

func SetDelimiters(leftDelimiters string, rightDelimiters string) {
	leftDelim = leftDelimiters
	rightDelim = rightDelimiters
}

func MustLoad() {
	funcMap := template.FuncMap{}
	MustLoadWithFuncs(funcMap)
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
	mainTemplate := mustGetMain(funcMap, leftDelim, rightDelim)

	for _, file := range includeFiles {
		filename := filepath.Base(file)
		files := append(layoutFiles, file)
		var err error
		templates[filename], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[filename] = template.Must(templates[filename].Delims(leftDelim, rightDelim).Funcs(funcMap).ParseFiles(files...))
	}
	log.Println("Successfully loaded templates!")

	bufpool = bpool.NewBufferPool(64)
	log.Println("Successfully allocated buffers!")
}
