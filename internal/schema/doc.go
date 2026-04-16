// Package schema validates secrets against a user-defined schema.
//
// Rules are declared as strings in the form "KEY:modifier:modifier",
// where modifiers include "required", "string", "int", and "bool".
//
// Example:
//
//	rules, _ := schema.ParseRules([]string{"DB_HOST:required", "PORT:required:int", "DEBUG:bool"})
//	v := schema.NewValidator(rules)
//	if err := v.Validate(secrets); err != nil {
//		log.Fatal(err)
//	}
package schema
