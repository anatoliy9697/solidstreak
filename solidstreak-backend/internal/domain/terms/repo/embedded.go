package repo

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/i18n"
	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

type embeddedRepo struct {
	currentTermsData map[string][]byte
	termsData        map[string]map[string][]byte
}

//go:embed content
var content embed.FS

var mdConverter = goldmark.New(
	goldmark.WithExtensions(extension.Linkify),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func renderHTML(md []byte) ([]byte, error) {
	var body bytes.Buffer

	if err := mdConverter.Convert(md, &body); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}

const termsPageTemplate = `<!doctype html>
<html lang="%s">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
body{font-family: system-ui, sans-serif;flex: 1 0 auto;margin: 0 auto;width: 100%%;max-width: 32rem;padding: 0 0.5rem;}
</style>
</head>
<body>
%s
</body>
</html>`

func wrapPage(lang string, body []byte) []byte {
	return fmt.Appendf(nil, termsPageTemplate, lang, i18n.MESSAGES[lang]["termsOfService"], body)
}

func initEmbeddedRepo() (Repo, error) {
	var err error

	var raw []byte
	if raw, err = content.ReadFile("content/current.json"); err != nil {
		return nil, err
	}

	var langDirs []fs.DirEntry
	if langDirs, err = fs.ReadDir(content, "content"); err != nil {
		return nil, err
	}

	termsData := make(map[string]map[string][]byte)

	for _, d := range langDirs {
		if !d.IsDir() {
			continue
		}

		lang := d.Name()

		termsData[lang] = make(map[string][]byte)

		var files []fs.DirEntry
		if files, err = fs.ReadDir(content, "content/"+lang); err != nil {
			return nil, err
		}

		for _, f := range files {
			version := strings.TrimSuffix(f.Name(), ".md")

			var md []byte
			if md, err = content.ReadFile("content/" + lang + "/" + f.Name()); err != nil {
				return nil, err
			}

			var html []byte
			if html, err = renderHTML(md); err != nil {
				return nil, err
			}

			termsData[lang][version] = wrapPage(lang, html)
		}
	}

	currentVersions := make(map[string]string)
	if err = json.Unmarshal(raw, &currentVersions); err != nil {
		return nil, err
	}

	currentTermsData := make(map[string][]byte)
	for lang, version := range currentVersions {
		if terms, ok := termsData[lang][version]; ok {
			currentTermsData[lang] = terms
		} else {
			return nil, fmt.Errorf("terms not found for language %v and version %v", lang, version)
		}
	}

	return &embeddedRepo{
		currentTermsData: currentTermsData,
		termsData:        termsData,
	}, nil
}

func (r embeddedRepo) GetByLangAndVersion(lang, version string) ([]byte, error) {
	if version == "" {
		if terms, ok := r.currentTermsData[lang]; ok {
			return terms, nil
		}

		return nil, apperrors.NewNotFoundErr(fmt.Sprintf("current terms not found for language %v", lang))
	}

	if terms, ok := r.termsData[lang][version]; ok {
		return terms, nil
	}

	return nil, apperrors.NewNotFoundErr(fmt.Sprintf("terms not found for language %v and version %v", lang, version))
}
