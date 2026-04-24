// Package output provides a lightweight formatter for rendering Vault secrets
// to an io.Writer in dotenv, shell-export, or JSON format.
//
// Supported formats:
//
//   - [FormatDotenv]: writes KEY=VALUE pairs suitable for .env files
//   - [FormatExport]: writes export KEY=VALUE pairs suitable for shell sourcing
//   - [FormatJSON]: writes a JSON object mapping keys to values
//
// Usage:
//
//	w := output.New(os.Stdout, output.FormatExport)
//	if err := w.Write(secrets); err != nil { ... }
package output
