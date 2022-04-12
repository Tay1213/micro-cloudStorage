package v1

import (
	"errors"
	"micro-cloudStorage/gateway/pkg/logging"
)

func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("userService--" + err.Error())
		logging.Info(err)
		panic(err)
	}
}

func PanicIfFileError(err error) {
	if err != nil {
		err = errors.New("taskService--" + err.Error())
		logging.Info(err)
		panic(err)
	}
}
