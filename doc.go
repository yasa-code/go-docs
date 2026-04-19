/*
The go-docs server code acts as a proxy for requests for go module documentation:
  - requests for internal modules such as `/github.com/private/*` are passed to the `origin` (local pkgsite) server
  - requests for all other modules are forwarded to the `remote` (https://pkg.go.dev) server
*/
package main
