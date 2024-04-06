package cierrelotedtos

import "time"

type PrismaClResultGroup struct {
	Cantidadregistro    int64
	Nroestablecimiento  string
	ExternalloteId      string
	Nombrearchivolote   string
	Monto               int64
	Fechaoperacion      string
	FechaCierre         string
	FechaAcreditacion   time.Time
	EstadoConciliacion  bool
	ExternalmediopagoId int64
	Nrocuota            int64
	BancoExternalId     int64
}
