// Package output provides a lightweight formatter for rendering Vault secrets
// to an io.Writer in dotenv, shell-export, or JSON format.
//
// Usage:
//
//	w := output.New(os.Stdout, output.FormatExport)
//	if err := w.Write(secrets); err != nil { ... }
package output
