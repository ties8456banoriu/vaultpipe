// Package secrettemplate provides template rendering for secret values.
//
// Use New to create a Renderer, then call Render with a Go template string
// and a map of secrets. Keys are accessed via {{ index . "KEY" }}.
// RenderAll renders multiple named templates in one call.
package secrettemplate
