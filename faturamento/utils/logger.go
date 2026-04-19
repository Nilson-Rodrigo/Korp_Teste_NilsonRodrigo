package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger inicializa o logger estruturado
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Info().Msg("Logger inicializado com sucesso")
}

// Logger retorna o logger global
func Logger() zerolog.Logger {
	return log.Logger
}
