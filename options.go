package distlock

import (
	"strings"
)

type OptFunc func(dl *Distlock)

func Prefix(value string) OptFunc {
	return func(dl *Distlock) {
		dl.prefix = strings.TrimSpace(value)
	}
}
