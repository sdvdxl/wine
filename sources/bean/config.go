package bean

import "github.com/sdvdxl/go-tools/collections"

type AuthPageConfig struct {
	LoginPages *collections.Set
	LoginPaths *collections.Set
}

type AuthPageConfigError struct {
}

func (this AuthPageConfig) Error() string {
	return "auth page config error"
}
