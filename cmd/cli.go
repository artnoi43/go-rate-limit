package main

import (
	"flag"

	"github.com/artnoi43/go-rate-limit/config"
)

type flags struct {
	maxGuard int
	URL      string
}

func (f *flags) parse(conf *config.Config) {
	flag.IntVar(&f.maxGuard, "c", conf.MaxGuard, "MaxGuard")
	flag.StringVar(&f.URL, "u", conf.URL, "URL")
	flag.Parse()
}
