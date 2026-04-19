// Package secrettemplate renders Go templates using secrets as template data.
package secrettemplate

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

// Renderer renders templates populated with secret values.
type Renderer struct{}

// New returns a new Renderer.
func New() *Renderer {
	return &Renderer{}
}

// Render executes the given template string with secrets as the data map.
// Keys in the template should be referenced as {{ index . "KEY" }}.
func (r *Renderer) Render(tmpl string, secrets map[string]string) (string, error) {
	if len(secrets) == 0 {
		return "", errors.New("secrettemplate: secrets must not be empty")
	}
	if tmpl == "" {
		return "", errors.New("secrettemplate: template must not be empty")
	}

	t, err := template.New("secret").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("secrettemplate: parse error: %w", err)
	}

	data := make(map[string]any, len(secrets))
	for k, v := range secrets {
		data[k] = v
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("secrettemplate: execute error: %w", err)
	}
	return buf.String(), nil
}

// RenderAll renders multiple named templates and returns a map of results.
func (r *Renderer) RenderAll(templates map[string]string, secrets map[string]string) (map[string]string, error) {
	if len(templates) == 0 {
		return nil, errors.New("secrettemplate: templates must not be empty")
	}
	out := make(map[string]string, len(templates))
	for name, tmpl := range templates {
		result, err := r.Render(tmpl, secrets)
		if err != nil {
			return nil, fmt.Errorf("secrettemplate: template %q: %w", name, err)
		}
		out[name] = result
	}
	return out, nil
}
