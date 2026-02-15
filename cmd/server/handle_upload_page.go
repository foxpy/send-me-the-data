package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "embed"
)

var (
	//go:embed templates/upload.gohtml
	uploadTemplateStr string

	uploadTemplate = template.Must(template.New("").Parse(uploadTemplateStr))
)

func (s *State) handleUploadPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := doesLinkExist(s.db, id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		respond404(w)
		return nil
	}

	files, err := listFiles(s.prefix, id)
	if err != nil {
		return fmt.Errorf("failed to obtain a list of files: %w", err)
	}

	var b bytes.Buffer
	err = uploadTemplate.Execute(&b, files)
	if err != nil {
		return fmt.Errorf("failed to render a template: %w", err)
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		return fmt.Errorf("failed to write rendered template: %w", err)
	}

	return nil
}
