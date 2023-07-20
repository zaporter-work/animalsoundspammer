package templategomodule

import (
	"errors"
	"strings"
)

type Config struct {
	Message string `json:"message"`
    Animals []string `json:"animals"`
}

// Validate takes the current location in the config (useful for good error messages).
// It should return a []string which contains all of the implicit
// dependencies of a module. (or nil,err if the config does not pass validation)
func (cfg *Config) Validate(path string) ([]string, error) {
	if strings.HasPrefix(cfg.Message, "failure") {
		return nil, errors.New(path + " Not permitted to have a message that begins with failure")
	}
	return make([]string, 0), nil
}
