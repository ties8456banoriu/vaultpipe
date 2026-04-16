// Package scopedenv provides environment-scoped secret namespacing for
// vaultpipe. It allows secrets fetched from Vault to be tagged with a
// named scope (e.g. "dev", "staging", "prod") and later filtered so
// only the relevant subset is written to a .env file.
//
// Typical usage:
//
//	s, _ := scopedenv.New("dev")
//	tagged, _ := s.Tag(allSecrets)   // prefix keys with DEV_
//	filtered, _ := s.Filter(tagged)  // strip prefix, keep only DEV_ keys
package scopedenv
