// Package secretsplit provides utilities for partitioning a flat secrets
// map into named groups based on key-prefix rules.
//
// Each rule associates a prefix string with a group name. Keys matching the
// prefix are placed into that group; unmatched keys are collected into a
// special "_default" group.
//
// Usage:
//
//	splitter, err := secretsplit.New([]secretsplit.Rule{
//		{Name: "db", Prefix: "DB_"},
//		{Name: "aws", Prefix: "AWS_"},
//	})
//	groups, err := splitter.Apply(secrets)
package secretsplit
