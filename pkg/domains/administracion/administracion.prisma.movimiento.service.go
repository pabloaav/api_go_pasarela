package administracion

import (
	"errors"
	"strconv"

	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

func (s *service) BuildPrismaMovimiento(reversion bool) (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

	// obtengo desde la BD los cierre de lotes
	listaPrismaCierreLote, err := s.repository.GetPrismaCierreLotes(reversion)
	// valido si existe error al obtner los datos desde DB
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	// valido si la lista esta vacia
	if len(listaPrismaCierreLote) == 0 {
		erro = errors.New(ERROR_LISTA_CIERRE_LOTE_VACIA)
		return
	}
	var transaccionesIds = make([]string, len(listaPrismaCierreLote))
	var ticketIds = make([]string, len(listaPrismaCierreLote))
	var codigosAutIds = make([]string, len(listaPrismaCierreLote))
	// recorro listaprismacierrelote para construir un slice con los pagosuuid del cierre de lote
	for k, v := range listaPrismaCierreLote {
		transaccionesIds[k] = v.PagosUuid
		ticketIds[k] = strconv.FormatInt(v.Nroticket, 10)
		codigosAutIds[k] = v.Codigoautorizacion
		//codigosAutIds[k] = strconv.FormatInt(v.Codigoautorizacion, 10)
	}
	filtroPagoIntentos := filtros.PagoIntentoFiltro{
		//TransaccionesId:      transaccionesIds,
		TicketNumber:         ticketIds,
		CodigoAutorizacion:   codigosAutIds,
		Channel:              true,
		CargarPago:           true,
		CargarPagoTipo:       true,
		CargarCuenta:         true,
		CargarCliente:        true,
		CargarPagoEstado:     true,
		ExternalId:           true,
		CargarMovimientos:    true,
		CargarCuentaComision: true,
		CargarImpuestos:      true,
	}
	// obtengo todos los pagos intentos relacionados con los tickets y codigos de autorizacion
	resultPagoIntentos, err := s.repository.GetPagosIntentos(filtroPagoIntentos)
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	for _, valueCierreLote := range listaPrismaCierreLote {
		filtroEstadoExterno := filtros.PagoEstadoExternoFiltro{Nombre: string(valueCierreLote.Tipooperacion)}

		estadosExternos, err := s.repository.GetPagosEstadosExternos(filtroEstadoExterno)
		/// se agrego el contro del  error durante la etapa del test unit
		if err != nil {
			erro = errors.New(err.Error())
			return
		}
		filtroEstadoPasarela := filtros.PagoEstadoFiltro{EstadoId: uint(estadosExternos[0].PagoestadosId)}
		estadosPasarela, err := s.repository.GetPagosEstados(filtroEstadoPasarela)
		if err != nil {
			erro = errors.New(err.Error())
		}

		for _, valuePagoIntento := range resultPagoIntentos {
			// transaccionid de pago intento es el mismo que pagouuid de cierrelote
			// FIXME se debe comprobar si external_id de pago intento es distinto de 0
			//if valuePagoIntento.TransactionID == valueCierreLote.PagosUuid {
			nroTicket := strconv.FormatInt(valueCierreLote.Nroticket, 10)
			codigoAutorizacion := valueCierreLote.Codigoautorizacion //strconv.FormatInt(valueCierreLote.Codigoautorizacion, 10)
			transaccionIdRecuperado := valuePagoIntento.TransactionID[0:15]
			if valuePagoIntento.TicketNumber == nroTicket && valuePagoIntento.AuthorizationCode == codigoAutorizacion && transaccionIdRecuperado == valueCierreLote.PagosUuid {
				// FIXME modifico pago estado id
				valuePagoIntento.Pago.PagoestadosID = int64(estadosExternos[0].PagoestadosId)
				// guardo el pago relacionado con el pago intento en movimientoCierreLote.ListaPagos
				//movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, valuePagoIntento.Pago)

				// creo un movimiento dependiendo del tipo de operacion
				movimiento := entities.Movimiento{}
				// se crea reversion
				reversion := entities.Reversione{}
				//////////////////////////////
				/*
					calcular comisiones e impuesto
					como paramentro se debe pasar:
					- id del channel.
					- lista de movimientos.
					- lista de cuenta comisiones y lista de impuestos.
				*/
				/*
					modificado 15-07-2022
					se agregar el id de cuenta
				*/
				var pagoCuotas bool
				if valueCierreLote.Nrocuota > 1 {
					pagoCuotas = true
				}
				var idMedioPago uint
				if valuePagoIntento.MediopagosID == 30 {
					idMedioPago = uint(valuePagoIntento.MediopagosID)
					pagoCuotas = true
				}

				filtroComisionChannel := filtros.CuentaComisionFiltro{
					CargarCuenta:      true,
					CuentaId:          valuePagoIntento.Pago.PagosTipo.Cuenta.ID,
					ChannelId:         valuePagoIntento.Mediopagos.Channel.ID,
					Mediopagoid:       idMedioPago,
					ExaminarPagoCuota: true,
					PagoCuota:         pagoCuotas,
					Channelarancel:    true,
					FechaPagoVigencia: valuePagoIntento.PaidAt,
				}
				cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
				if err != nil {
					erro = errors.New(err.Error() + " de comisiones")
					return
				}
				movimiento.Enobservacion = valueCierreLote.Enobservacion
				/////////////////////////////
				// controlar si el importe se obtiene de amount o valor cupon
				var importe entities.Monto
				importe = valuePagoIntento.Amount
				if valuePagoIntento.Valorcupon != 0 {
					importe = valuePagoIntento.Valorcupon
				}
				if valueCierreLote.Tipooperacion == "C" && !valueCierreLote.Reversion && !valueCierreLote.Conciliado {
					// genero un movimiento tipo debito, y le paso cuentaid, pagointentoid y el monto
					// modifico en pago tipo los campos availableat y revertedat
					movimiento.AddCredito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), importe) //valuePagoIntento.Amount) // valueCierreLote.Monto)
					movimiento.Reversion = false
					valuePagoIntento.AvailableAt = valueCierreLote.FechaCierre
					valuePagoIntento.RevertedAt = time.Time{}
					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva, valuePagoIntento.Amount) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions

				}
				if (valueCierreLote.Tipooperacion == "D" && !valueCierreLote.Reversion && !valueCierreLote.Conciliado) || (valueCierreLote.Tipooperacion == "C" && valueCierreLote.Reversion && valueCierreLote.Conciliado) {
					/*
					   FIXME SE DEBE REVISAR QUE LA OPERACIÓN DE REVERSIÓN DEL CIERRE DE LOTE, tenga generado su movimiento
					*/
					// obtengo un movimiento relacionado con un pago intento
					//if len(valuePagoIntento.Movimientos) != 0 {
					/*
						si existe movimintos ralacionados con el pago intento,
						genero un movimiento negativo y calculo la comision e impuesto en valor negativo
					*/

					movimiento.AddCredito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), -1.00*importe) //-1.00*valuePagoIntento.Amount)
					// movimiento.AddDebito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), -1.00*valuePagoIntento.Amount) //valueCierreLote.Monto)
					movimiento.Reversion = true
					valuePagoIntento.RevertedAt = valueCierreLote.FechaCierre
					valuePagoIntento.AvailableAt = time.Time{}
					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva, valuePagoIntento.Amount) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions
					status := estadosPasarela[0].Estado
					//actualizo el estado a revertido
					valuePagoIntento.Pago.PagoestadosID = int64(estadosPasarela[0].ID)
					reversion.AddReversion(valuePagoIntento.ID, -1.00*valueCierreLote.Monto.Int64(), valueCierreLote.ExternalclienteID, string(status))
					// se guarda el reversiones en Listareversiones
					if movimiento.MotivoBaja == "" {
						movimientoCierreLote.ListaReversiones = append(movimientoCierreLote.ListaReversiones, reversion)
					}
				}

				if movimiento.MotivoBaja == "" {
					if movimiento.PagointentosId != 0 {
						// se guarda el movimiento en Listamovimientos
						movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
					}
					// se guarda pagointento en listapagointentos
					movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, valuePagoIntento)
					// se crea el objeto pago estadolog
					pagoEstadoLog := entities.Pagoestadologs{
						PagosID:       int64(valuePagoIntento.Pago.ID),
						PagoestadosID: int64(estadosExternos[0].PagoestadosId),
					}
					// se guarda pago pagoestadolog en listapagosestadologs
					movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)
					// guardo el pago relacionado con el pago intento en movimientoCierreLote.ListaPagos
					movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, valuePagoIntento.Pago)

				}
			}
		}
	}
	movimientoCierreLote.ListaCLPrisma = listaPrismaCierreLote
	return
}

//////////old 22-09-2022//////////////
// func (s *service) BuildPrismaMovimiento() (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

// 	// obtengo desde la BD los cierre de lotes
// 	listaPrismaCierreLote, err := s.repository.GetPrismaCierreLotes()
// 	// valido si existe error al obtner los datos desde DB
// 	if err != nil {
// 		erro = errors.New(err.Error())
// 		return
// 	}
// 	// valido si la lista esta vacia
// 	if len(listaPrismaCierreLote) == 0 {
// 		erro = errors.New(ERROR_LISTA_CIERRE_LOTE_VACIA)
// 		return
// 	}
// 	var transaccionesIds = make([]string, len(listaPrismaCierreLote))
// 	var ticketIds = make([]string, len(listaPrismaCierreLote))
// 	var codigosAutIds = make([]string, len(listaPrismaCierreLote))
// 	// recorro listaprismacierrelote para construir un slice con los pagosuuid del cierre de lote
// 	for k, v := range listaPrismaCierreLote {
// 		transaccionesIds[k] = v.PagosUuid
// 		ticketIds[k] = strconv.FormatInt(v.Nroticket, 10)
// 		codigosAutIds[k] = v.Codigoautorizacion
// 		//codigosAutIds[k] = strconv.FormatInt(v.Codigoautorizacion, 10)
// 	}
// 	filtroPagoIntentos := filtros.PagoIntentoFiltro{
// 		//TransaccionesId:      transaccionesIds,
// 		TicketNumber:         ticketIds,
// 		CodigoAutorizacion:   codigosAutIds,
// 		Channel:              true,
// 		CargarPago:           true,
// 		CargarPagoTipo:       true,
// 		CargarCuenta:         true,
// 		CargarCliente:        true,
// 		CargarPagoEstado:     true,
// 		ExternalId:           true,
// 		CargarMovimientos:    true,
// 		CargarCuentaComision: true,
// 		CargarImpuestos:      true,
// 	}
// 	// obtengo todos los pagos intentos relacionados con los tickets y codigos de autorizacion
// 	resultPagoIntentos, err := s.repository.GetPagosIntentos(filtroPagoIntentos)
// 	if err != nil {
// 		erro = errors.New(err.Error())
// 		return
// 	}
// 	for _, valueCierreLote := range listaPrismaCierreLote {
// 		filtroEstadoExterno := filtros.PagoEstadoExternoFiltro{Nombre: string(valueCierreLote.Tipooperacion)}

// 		estadosExternos, err := s.repository.GetPagosEstadosExternos(filtroEstadoExterno)
// 		/// se agrego el contro del  error durante la etapa del test unit
// 		if err != nil {
// 			erro = errors.New(err.Error())
// 			return
// 		}
// 		filtroEstadoPasarela := filtros.PagoEstadoFiltro{EstadoId: uint(estadosExternos[0].PagoestadosId)}
// 		estadosPasarela, err := s.repository.GetPagosEstados(filtroEstadoPasarela)
// 		if err != nil {
// 			erro = errors.New(err.Error())
// 		}

// 		for _, valuePagoIntento := range resultPagoIntentos {
// 			// transaccionid de pago intento es el mismo que pagouuid de cierrelote
// 			// FIXME se debe comprobar si external_id de pago intento es distinto de 0
// 			//if valuePagoIntento.TransactionID == valueCierreLote.PagosUuid {
// 			nroTicket := strconv.FormatInt(valueCierreLote.Nroticket, 10)
// 			codigoAutorizacion := valueCierreLote.Codigoautorizacion //strconv.FormatInt(valueCierreLote.Codigoautorizacion, 10)
// 			transaccionIdRecuperado := valuePagoIntento.TransactionID[0:15]
// 			if valuePagoIntento.TicketNumber == nroTicket && valuePagoIntento.AuthorizationCode == codigoAutorizacion && transaccionIdRecuperado == valueCierreLote.PagosUuid {
// 				// FIXME modifico pago estado id
// 				valuePagoIntento.Pago.PagoestadosID = int64(estadosExternos[0].PagoestadosId)
// 				// guardo el pago relacionado con el pago intento en movimientoCierreLote.ListaPagos
// 				//movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, valuePagoIntento.Pago)

// 				// creo un movimiento dependiendo del tipo de operacion
// 				movimiento := entities.Movimiento{}
// 				// se crea reversion
// 				reversion := entities.Reversione{}
// 				//////////////////////////////
// 				/*
// 					calcular comisiones e impuesto
// 					como paramentro se debe pasar:
// 					- id del channel.
// 					- lista de movimientos.
// 					- lista de cuenta comisiones y lista de impuestos.
// 				*/
// 				/*
// 					modificado 15-07-2022
// 					se agregar el id de cuenta
// 				*/
// 				var pagoCuotas bool
// 				if valueCierreLote.Nrocuota > 1 {
// 					pagoCuotas = true
// 				}
// 				var idMedioPago uint
// 				if valuePagoIntento.MediopagosID == 30 {
// 					idMedioPago = uint(valuePagoIntento.MediopagosID)
// 					pagoCuotas = true
// 				}

// 				filtroComisionChannel := filtros.CuentaComisionFiltro{
// 					CargarCuenta:      true,
// 					CuentaId:          valuePagoIntento.Pago.PagosTipo.Cuenta.ID,
// 					ChannelId:         valuePagoIntento.Mediopagos.Channel.ID,
// 					Mediopagoid:       idMedioPago,
// 					ExaminarPagoCuota: true,
// 					PagoCuota:         pagoCuotas,
// 				}
// 				cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
// 				if err != nil {
// 					erro = errors.New(err.Error() + " de comisiones")
// 					return
// 				}
// 				/////////////////////////////
// 				if valueCierreLote.Tipooperacion == "C" {
// 					// genero un movimiento tipo debito, y le paso cuentaid, pagointentoid y el monto
// 					// modifico en pago tipo los campos availableat y revertedat
// 					movimiento.AddCredito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), valuePagoIntento.Amount) // valueCierreLote.Monto)
// 					valuePagoIntento.AvailableAt = valueCierreLote.FechaCierre
// 					valuePagoIntento.RevertedAt = time.Time{}
// 					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
// 					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions

// 				} else {
// 					/*
// 					   FIXME SE DEBE REVISAR QUE LA OPERACIÓN DE REVERSIÓN DEL CIERRE DE LOTE, tenga generado su movimiento
// 					*/
// 					// obtengo un movimiento relacionado con un pago intento
// 					//if len(valuePagoIntento.Movimientos) != 0 {
// 					/*
// 						si existe movimintos ralacionados con el pago intento,
// 						genero un movimiento negativo y calculo la comision e impuesto en valor negativo
// 					*/
// 					movimiento.AddDebito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), -1.00*valuePagoIntento.Amount) //valueCierreLote.Monto)
// 					valuePagoIntento.RevertedAt = valueCierreLote.FechaCierre
// 					valuePagoIntento.AvailableAt = time.Time{}
// 					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
// 					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions
// 					status := estadosPasarela[0].Estado
// 					//actualizo el estado a revertido
// 					valuePagoIntento.Pago.PagoestadosID = int64(estadosPasarela[0].ID)
// 					reversion.AddReversion(valuePagoIntento.ID, -1.00*valueCierreLote.Monto.Int64(), valueCierreLote.ExternalclienteID, string(status))
// 					// se guarda el reversiones en Listareversiones
// 					movimientoCierreLote.ListaReversiones = append(movimientoCierreLote.ListaReversiones, reversion)
// 				}
// 				if movimiento.PagointentosId != 0 {
// 					// se guarda el movimiento en Listamovimientos
// 					movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
// 				}
// 				// se guarda pagointento en listapagointentos
// 				movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, valuePagoIntento)
// 				// se crea el objeto pago estadolog
// 				pagoEstadoLog := entities.Pagoestadologs{
// 					PagosID:       int64(valuePagoIntento.Pago.ID),
// 					PagoestadosID: int64(estadosExternos[0].PagoestadosId),
// 				}
// 				// se guarda pago pagoestadolog en listapagosestadologs
// 				movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)
// 				// guardo el pago relacionado con el pago intento en movimientoCierreLote.ListaPagos
// 				movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, valuePagoIntento.Pago)
// 			}
// 		}
// 	}
// 	movimientoCierreLote.ListaCLPrisma = listaPrismaCierreLote
// 	return
// }

///////////////////most old////////////////
// func (s *service) BuildPrismaMovimiento() (movimientoCierreLote administraciondtos.MovimientoCierreLoteResponse, erro error) {

// 	// obtengo desde la BD los cierre de lotes
// 	listaPrismaCierreLote, err := s.repository.GetPrismaCierreLotes()
// 	// valido si existe error al obtner los datos desde DB
// 	if err != nil {
// 		erro = errors.New(err.Error())
// 		return
// 	}
// 	// valido si la lista esta vacia
// 	if len(listaPrismaCierreLote) == 0 {
// 		erro = errors.New(ERROR_LISTA_CIERRE_LOTE_VACIA)
// 		return
// 	}
// 	var transaccionesIds = make([]string, len(listaPrismaCierreLote))
// 	var ticketIds = make([]string, len(listaPrismaCierreLote))
// 	var codigosAutIds = make([]string, len(listaPrismaCierreLote))
// 	// recorro listaprismacierrelote para construir un slice con los pagosuuid del cierre de lote
// 	for k, v := range listaPrismaCierreLote {
// 		transaccionesIds[k] = v.PagosUuid
// 		ticketIds[k] = strconv.FormatInt(v.Nroticket, 10)
// 		codigosAutIds[k] = v.Codigoautorizacion
// 		//codigosAutIds[k] = strconv.FormatInt(v.Codigoautorizacion, 10)
// 	}
// 	filtroPagoIntentos := filtros.PagoIntentoFiltro{
// 		//TransaccionesId:      transaccionesIds,
// 		TicketNumber:         ticketIds,
// 		CodigoAutorizacion:   codigosAutIds,
// 		Channel:              true,
// 		CargarPago:           true,
// 		CargarPagoTipo:       true,
// 		CargarCuenta:         true,
// 		CargarCliente:        true,
// 		CargarPagoEstado:     true,
// 		ExternalId:           true,
// 		CargarMovimientos:    true,
// 		CargarCuentaComision: true,
// 		CargarImpuestos:      true,
// 	}
// 	// obtengo todos los pagos intentos relacionados con los tickets y codigos de autorizacion
// 	resultPagoIntentos, err := s.repository.GetPagosIntentos(filtroPagoIntentos)
// 	if err != nil {
// 		erro = errors.New(err.Error())
// 		return
// 	}
// 	for _, valueCierreLote := range listaPrismaCierreLote {
// 		filtroEstadoExterno := filtros.PagoEstadoExternoFiltro{Nombre: string(valueCierreLote.Tipooperacion)}

// 		estadosExternos, err := s.repository.GetPagosEstadosExternos(filtroEstadoExterno)
// 		/// se agrego el contro del  error durante la etapa del test unit
// 		if err != nil {
// 			erro = errors.New(err.Error())
// 			return
// 		}
// 		filtroEstadoPasarela := filtros.PagoEstadoFiltro{EstadoId: uint(estadosExternos[0].PagoestadosId)}
// 		estadosPasarela, err := s.repository.GetPagosEstados(filtroEstadoPasarela)
// 		if err != nil {
// 			erro = errors.New(err.Error())
// 		}

// 		for _, valuePagoIntento := range resultPagoIntentos {
// 			// transaccionid de pago intento es el mismo que pagouuid de cierrelote
// 			// FIXME se debe comprobar si external_id de pago intento es distinto de 0
// 			//if valuePagoIntento.TransactionID == valueCierreLote.PagosUuid {
// 			nroTicket := strconv.FormatInt(valueCierreLote.Nroticket, 10)
// 			codigoAutorizacion := valueCierreLote.Codigoautorizacion //strconv.FormatInt(valueCierreLote.Codigoautorizacion, 10)
// 			transaccionIdRecuperado := valuePagoIntento.TransactionID[0:15]
// 			if valuePagoIntento.TicketNumber == nroTicket && valuePagoIntento.AuthorizationCode == codigoAutorizacion && transaccionIdRecuperado == valueCierreLote.PagosUuid {
// 				// FIXME modifico pago estado id
// 				valuePagoIntento.Pago.PagoestadosID = int64(estadosExternos[0].PagoestadosId)
// 				// guardo el pago relacionado con el pago intento en movimientoCierreLote.ListaPagos
// 				movimientoCierreLote.ListaPagos = append(movimientoCierreLote.ListaPagos, valuePagoIntento.Pago)

// 				// creo un movimiento dependiendo del tipo de operacion
// 				movimiento := entities.Movimiento{}
// 				// se crea reversion
// 				reversion := entities.Reversione{}
// 				//////////////////////////////
// 				/*
// 					calcular comisiones e impuesto
// 					como paramentro se debe pasar:
// 					- id del channel.
// 					- lista de movimientos.
// 					- lista de cuenta comisiones y lista de impuestos.
// 				*/
// 				/*
// 					modificado 15-07-2022
// 					se agregar el id de cuenta
// 				*/
// 				var pagoCuotas bool
// 				if valueCierreLote.Nrocuota > 1 {
// 					pagoCuotas = true
// 				}
// 				var idMedioPago uint
// 				if valuePagoIntento.MediopagosID == 30 {
// 					idMedioPago = uint(valuePagoIntento.MediopagosID)
// 					pagoCuotas = true
// 				}

// 				filtroComisionChannel := filtros.CuentaComisionFiltro{
// 					CargarCuenta:      true,
// 					CuentaId:          valuePagoIntento.Pago.PagosTipo.Cuenta.ID,
// 					ChannelId:         valuePagoIntento.Mediopagos.Channel.ID,
// 					Mediopagoid:       idMedioPago,
// 					ExaminarPagoCuota: true,
// 					PagoCuota:         pagoCuotas,
// 				}
// 				cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
// 				if err != nil {
// 					erro = errors.New(err.Error() + " de comisiones")
// 					return
// 				}
// 				/////////////////////////////
// 				if valueCierreLote.Tipooperacion == "C" {
// 					// genero un movimiento tipo debito, y le paso cuentaid, pagointentoid y el monto
// 					// modifico en pago tipo los campos availableat y revertedat
// 					movimiento.AddCredito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), valuePagoIntento.Amount) // valueCierreLote.Monto)
// 					valuePagoIntento.AvailableAt = valueCierreLote.FechaCierre
// 					valuePagoIntento.RevertedAt = time.Time{}
// 					// filtroMedioPago := filtros.FiltroMedioPago{
// 					// 	IdMedioPago: valuePagoIntento.MediopagosID,
// 					// }
// 					// medioPago, err := s.repository.GetMedioPagoRepository(filtroMedioPago)
// 					// if err != nil {
// 					// 	erro = errors.New(err.Error())
// 					// 	return
// 					// }
// 					// var amexEstado bool
// 					// if medioPago.ID == 30 {
// 					// 	amexEstado = true
// 					// }
// 					/*
// 						se verifica si el medio de pago es Amex
// 					*/
// 					// var idMedioPago uint
// 					// if valuePagoIntento.MediopagosID == 30 {
// 					// 	idMedioPago = uint(valuePagoIntento.MediopagosID)
// 					// }
// 					// /*
// 					// 	calcular comisiones e impuesto
// 					// 	como paramentro se debe pasar:
// 					// 	- id del channel.
// 					// 	- lista de movimientos.
// 					// 	- lista de cuenta comisiones y lista de impuestos.
// 					// */
// 					// /*
// 					// 	modificado 15-07-2022
// 					// 	se agregar el id de cuenta
// 					// */
// 					// filtroComisionChannel := filtros.CuentaComisionFiltro{
// 					// 	CargarCuenta: true,
// 					// 	CuentaId:     valuePagoIntento.Pago.PagosTipo.Cuenta.ID,
// 					// 	ChannelId:    valuePagoIntento.Mediopagos.Channel.ID,
// 					// 	Mediopagoid:  idMedioPago,
// 					// }
// 					// cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
// 					// if err != nil {
// 					// 	erro = errors.New(err.Error() + " de comisiones")
// 					// 	return
// 					// }
// 					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
// 					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions

// 				} else {
// 					/*
// 					   FIXME SE DEBE REVISAR QUE LA OPERACIÓN DE REVERSIÓN DEL CIERRE DE LOTE, tenga generado su movimiento
// 					*/
// 					// obtengo un movimiento relacionado con un pago intento
// 					//if len(valuePagoIntento.Movimientos) != 0 {
// 					/*
// 						si existe movimintos ralacionados con el pago intento,
// 						genero un movimiento negativo y calculo la comision e impuesto en valor negativo
// 					*/
// 					movimiento.AddDebito(uint64(valuePagoIntento.Pago.PagosTipo.CuentasID), uint64(valuePagoIntento.ID), -1.00*valuePagoIntento.Amount) //valueCierreLote.Monto)
// 					valuePagoIntento.RevertedAt = valueCierreLote.FechaCierre
// 					valuePagoIntento.AvailableAt = time.Time{}
// 					// filtroComisionChannel := filtros.CuentaComisionFiltro{
// 					// 	CargarCuenta: true,
// 					// 	CuentaId:     valuePagoIntento.Pago.PagosTipo.Cuenta.ID,
// 					// 	ChannelId:    valuePagoIntento.Mediopagos.Channel.ID,
// 					// }
// 					// cuentaComision, err := s.repository.GetCuentaComision(filtroComisionChannel)
// 					// if err != nil {
// 					// 	erro = errors.New(err.Error())
// 					// 	return
// 					// }
// 					comisiones := append([]entities.Cuentacomision{}, cuentaComision)
// 					s.utilService.BuildComisiones(&movimiento, &comisiones, valuePagoIntento.Pago.PagosTipo.Cuenta.Cliente.Iva) //valuePagoIntento.Pago.PagosTipo.Cuenta.Cuentacomisions
// 					//}
// 					status := estadosPasarela[0].Estado
// 					reversion.AddReversion(valuePagoIntento.ID, -1.00*valueCierreLote.Monto.Int64(), valueCierreLote.ExternalclienteID, string(status))
// 					// se guarda el reversiones en Listareversiones
// 					movimientoCierreLote.ListaReversiones = append(movimientoCierreLote.ListaReversiones, reversion)
// 				}
// 				if movimiento.PagointentosId != 0 {
// 					// se guarda el movimiento en Listamovimientos
// 					movimientoCierreLote.ListaMovimientos = append(movimientoCierreLote.ListaMovimientos, movimiento)
// 				}
// 				// se guarda pagointento en listapagointentos
// 				movimientoCierreLote.ListaPagoIntentos = append(movimientoCierreLote.ListaPagoIntentos, valuePagoIntento)
// 				// se crea el objeto pago estadolog
// 				pagoEstadoLog := entities.Pagoestadologs{
// 					PagosID:       int64(valuePagoIntento.Pago.ID),
// 					PagoestadosID: int64(estadosExternos[0].PagoestadosId),
// 				}
// 				// se guarda pago pagoestadolog en listapagosestadologs
// 				movimientoCierreLote.ListaPagosEstadoLogs = append(movimientoCierreLote.ListaPagosEstadoLogs, pagoEstadoLog)
// 			}
// 		}
// 	}
// 	movimientoCierreLote.ListaCLPrisma = listaPrismaCierreLote
// 	return
// }
//////////////////////////////////
