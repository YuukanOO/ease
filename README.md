# ease : generating the boring stuff

_This is an early experimentation, missing a lot of what I have in mind. See my [blog post](https://julien.leicher.me/writes/static-analysis-ftw) (in French) if you wish to know more._

## Motivation

Exposing a service as an HTTP API should be easy. By leveraging minimal annotations and your code types, **ease** aims to generate the _glue_ code so you can focus on delivering business value.

Instead of generating _ease specific_ code, it only generates code using well maintained packages.

It should also handle the usage of a **remote module** annotated with **ease** to integrate it without effort.

## Example

The [Todo example](/examples/todo/) demonstrates how **ease** could ease (got it?) the development process. For example, the package `todo` only contains use cases (defined in [service.go](/examples/todo/service.go)) and use `go generate` to write the needed stuff to expose an HTTP server exposing needed endpoints, instantiating service dependencies as needed without configuring anything by using static analysis.

It also demonstrates how [an external module](https://github.com/YuukanOO/ease-external-example) can be integrated easily and added to the generated output.
