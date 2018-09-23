package tplmgr

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/volatiletech/authboss"
)

type AuthbossHTMLRenderer struct {
	extension string
}

func NewAuthbossHTMLRenderer() *AuthbossHTMLRenderer {
	return &AuthbossHTMLRenderer{}
}

func NewAuthbossHTMLRendererWithExt(extension string) *AuthbossHTMLRenderer {
	return &AuthbossHTMLRenderer{
		extension: extension,
	}
}

func (abhr *AuthbossHTMLRenderer) SetExtension(extension string) {
	abhr.extension = extension
}

func (abhr *AuthbossHTMLRenderer) Load(names ...string) error {
	return nil
}

func (abhr *AuthbossHTMLRenderer) Render(ctx context.Context, name string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	if !strings.HasSuffix(name, abhr.extension) {
		name += abhr.extension
	}
	template, ok := templates[name]
	if !ok {
		return nil, "", errors.Errorf("Template for page %s not found", name)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err = template.Execute(buf, data)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to render template for page %s", name)
	}

	return buf.Bytes(), "text/html", nil
}
