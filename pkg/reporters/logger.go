package reporters

import (
	"fmt"
	"go.uber.org/zap"
)

const production = "production"

func getLogger(env string) *zap.Logger {
	var err error
	var lgr *zap.Logger

	if env == production {
		lgr, err = zap.NewProduction()
	} else {
		lgr, err = zap.NewDevelopment()
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return lgr
}
