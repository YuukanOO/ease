package main

import (
	"github.com/YuukanOO/ease/pkg/generator"
	"github.com/YuukanOO/ease/pkg/parser"
)

type (
	options struct {
		packages   []string
		outputDir  string
		parsers    []parser.Extension
		generators []generator.Extension
	}

	Option func(*options)
)

// Run ease with the following options.
func Run(opts ...Option) error {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	parseResult, err := parser.
		New(o.parsers...).
		Parse(o.packages...)

	if err != nil {
		return err
	}

	return generator.
		New(o.outputDir, o.generators...).
		Generate(parseResult)
}

// Add packages to be parsed.
func WithPackages(packages ...string) Option {
	return func(o *options) {
		o.packages = packages
	}
}

// WithParsers set the parsers to be used.
func WithParsers(parsers ...parser.Extension) Option {
	return func(o *options) {
		o.parsers = parsers
	}
}

// WithGenerators set the generators to be used and the output directory.
func WithGenerators(outputDir string, generators ...generator.Extension) Option {
	return func(o *options) {
		o.outputDir = outputDir
		o.generators = generators
	}
}
