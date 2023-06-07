# ease : generating the boring stuff

## Motivation

Exposing a service as an HTTP API should be easy. By leveraging minimal annotations and your code types, **ease** aims to generate the _glue_ code so you can focus on delivering business value.

Instead of generating _ease specific_ code, it only generates code using well maintained packages.

It should also handle the usage of a remote module annotated with **ease** to integrate it without effort.

## First version goal

Annotations should be **required only if there is no other way** to retrieve the information.

- [ ] Generate router stuff (Gin for now) to call the appropriate service on endpoint
- [ ] Recursively resolve service dependencies (by finding appropriate creator functions)
- [ ] Generate OpenAPI Specs directly from code

## Other ideas

- [ ] Handle configuration via env file and inject options structures
- [ ] Handle database migrations
