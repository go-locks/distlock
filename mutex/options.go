package mutex

import (
	"time"
)

type options struct {
	name         string
	expiry       time.Duration
	factor       float64
	defaultWait  time.Duration
	costTopLimit time.Duration
}

func newOptions(name string, optFuncs ...OptFunc) options {
	opts := options{
		name:   name,
		expiry: 5 * time.Second,
		factor: 0.12,
	}
	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}
	factValue := float64(opts.expiry) * opts.factor
	opts.defaultWait = time.Duration(factValue)
	opts.costTopLimit = opts.expiry - time.Duration(factValue)
	return opts
}

type OptFunc func(opts *options)

func Expiry(value time.Duration) OptFunc {
	return func(opts *options) {
		opts.expiry = value
	}
}

func Factor(value float64) OptFunc {
	return func(opts *options) {
		opts.factor = value
	}
}
