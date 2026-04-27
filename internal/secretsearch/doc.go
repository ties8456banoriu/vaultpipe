// Package secretsearch provides key and value search over a secrets map.
//
// Three search modes are supported:
//   - exact: the query must exactly match the target field
//   - prefix: the target field must begin with the query
//   - regex: the query is compiled as a regular expression
//
// Either keys, values, or both may be searched in a single pass.
package secretsearch
