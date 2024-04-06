package cierrelotedtos

import "time"

type ApilinkRequest struct {
	FechaInicio time.Time
	FechaFin    time.Time
	Number      uint32
	Size        uint32
}
