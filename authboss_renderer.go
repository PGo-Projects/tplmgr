package tplmgr

import (
	"context"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
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

func AuthbossSAHTMLRenderer(w http.ResponseWriter, r *http.Request, name string, extension string data authboss.HTMLData) {
	var htmlData authboss.HTMLData
	contextData := r.Context().Value(authboss.CTXKeyData)
	if contextData == nil {
		htmlData = authboss.HTMLData{}
	} else {
		htmlData = contextData.(authboss.HTMLData)
	}

	htmlData.MergeKV("csrf_token", nosurf.Token(r))
	htmlData.Merge(data)

	Render(w, name+extension, htmlData)
}
