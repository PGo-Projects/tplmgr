package tplmgr

import (
	"html/template"
	"log"
	"path/filepath"
)

func mustGetLayoutFiles() (layoutFiles []string) {
	layoutFiles, err := filepath.Glob(config.layoutPath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func mustGetIncludeFiles() (includeFiles []string) {
	includeFiles, err := filepath.Glob(config.includePath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func mustGetMain(funcMap template.FuncMap, leftDelim string, rightDelim string) (mainTempl *template.Template) {
	mainTempl, err := template.New("main").Funcs(funcMap).Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	return
}
