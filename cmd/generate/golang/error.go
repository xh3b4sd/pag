package golang

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var commandExecutionFailedError = &tracer.Error{
	Kind: "commandExecutionFailedError",
}

func IsCommandExecutionFailed(err error) bool {
	return errors.Is(err, commandExecutionFailedError)
}

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}

var invalidFlagError = &tracer.Error{
	Kind: "invalidFlagError",
}

func IsInvalidFlag(err error) bool {
	return errors.Is(err, invalidFlagError)
}
