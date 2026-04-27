// Package secretscore provides quality scoring for secret values.
//
// A Scorer evaluates each secret against configurable criteria:
//   - Minimum value length
//   - Presence of at least one uppercase letter
//   - Presence of at least one digit
//
// Each criterion contributes one point toward a MaxPoints total.
// Use Score.Percent() to obtain a 0–100 quality percentage.
package secretscore
