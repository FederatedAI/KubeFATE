package service

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func debug(format string, v ...interface{}) {
	s := fmt.Sprintf("helm debug %s", format)
	log.Debug().Msgf(s, v...)
	//log.Debug().Object()
	//log.Debug().Msgf(format, v)
}
