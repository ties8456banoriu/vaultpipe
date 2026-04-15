// Package transform provides post-fetch value transformation for secrets
// retrieved from HashiCorp Vault.
//
// Transformations are applied after secrets are fetched and filtered, but
// before they are written to the .env file. Each rule targets a specific
// secret key and specifies a transformation type such as "upper", "lower",
// "trim", "prefix", or "suffix".
//
// Rules can be supplied programmatically via ParseRules from a slice of
// strings in the format "KEY:TYPE" or "KEY:TYPE:ARG".
package transform
