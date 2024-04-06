package administracion

import (
	"fmt"
	"time"

	ribcradtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/administraciondtos/ribcra"
)

func (r *repository) BuildRICuentasCliente(request ribcradtos.RICuentasClienteRequest) (ri []ribcradtos.RiCuentaCliente, erro error) {
	type RiCC struct {
		Cantidad   uint64
		Saldo      float64
		Dd         uint16
		Created_at time.Time
	}
	var result []RiCC
	fechaInicio := request.FechaInicio.Format("2006-01-02")
	fechaFin := request.FechaFin.Format("2006-01-02")

	response := r.SQLClient.Raw("select count(x.cuentas_id) as cantidad, sum(x.saldo) as saldo, DATE_FORMAT(x.created_at, '%d') as dd, cast( x.created_at as date) as created_at  from (select cuentas_id, sum(monto) as saldo, created_at from movimientos group by cast(created_at as date), cuentas_id) as x where cast(x.created_at AS DATE) between ? and ? group by cast(x.created_at as date)", fechaInicio, fechaFin).Scan(&result)

	if response.Error != nil {
		erro = response.Error
		return
	}

	for i := range result {

		ri10000 := ribcradtos.RiCuentaCliente{
			CodigoPartida: fmt.Sprintf("10000%d", result[i].Dd),
			Saldo:         fmt.Sprint(result[i].Saldo),
			Cantidad:      fmt.Sprint(result[i].Cantidad),
			CBU:           "NULO",
			Orden:         i,
		}
		ri = append(ri, ri10000)

		var cuentasTotal uint64
		fecha := result[i].Created_at.Format("2006-01-02")

		response = r.SQLClient.Raw("SELECT count(id) as CuentasTotal FROM pasarela.cuentas where deleted_at is null and cast(created_at as date) <= cast(? as date)", fecha).Scan(&cuentasTotal)

		if response.Error != nil {
			erro = response.Error
			return
		}

		ri50000 := ribcradtos.RiCuentaCliente{
			CodigoPartida: fmt.Sprintf("50000%d", result[i].Dd),
			Saldo:         "NULO",
			Cantidad:      fmt.Sprint(cuentasTotal),
			CBU:           "NULO",
			Orden:         i,
		}
		ri = append(ri, ri50000)

		ri20000 := ribcradtos.RiCuentaCliente{
			CodigoPartida: fmt.Sprintf("20000%d", result[i].Dd),
			Saldo:         fmt.Sprint(result[i].Saldo),
			Cantidad:      "NULO",
			CBU:           request.CbuCuentaTelco,
			Orden:         i,
		}

		ri = append(ri, ri20000)

		ri30000 := ribcradtos.RiCuentaCliente{
			CodigoPartida: fmt.Sprintf("30000%d", result[i].Dd),
			Saldo:         "0",
			Cantidad:      "0",
			CBU:           "NULO",
			Orden:         i,
		}
		ri = append(ri, ri30000)

		ri40000 := ribcradtos.RiCuentaCliente{
			CodigoPartida: fmt.Sprintf("40000%d", result[i].Dd),
			Saldo:         "0",
			Cantidad:      "NULO",
			CBU:           "0",
			Orden:         i,
		}
		ri = append(ri, ri40000)

	}

	return
}

func (r *repository) BuildRIDatosFondo(request ribcradtos.RiDatosFondosRequest) (ri []ribcradtos.RiDatosFondos, erro error) {

	// FIXME Estos datos hay que borrar se uso solamente para cargar algo en el front.
	ri = []ribcradtos.RiDatosFondos{
		{Numero: "15000", Denominacion: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Agente: ribcradtos.AgenteCustodia, DenominacionAgente: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", CuitAgente: "22222222222"},
		{Numero: "15001", Denominacion: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Agente: ribcradtos.AgenteAdministracion, DenominacionAgente: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", CuitAgente: "22222222223"},
		{Numero: "15002", Denominacion: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Agente: ribcradtos.AgenteColocacion, DenominacionAgente: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", CuitAgente: "22222222224"},
	}

	return
}

func (r *repository) BuilRIInfestaditica(request ribcradtos.RiInfestadisticaRequest) (ri []ribcradtos.RiInfestadistica, erro error) {

	fechaInicio := request.FechaInicio.Format("2006-01-02")
	fechaFin := request.FechaFin.Format("2006-01-02")

	//ApWebBotonDePago
	var ri1010400 []ribcradtos.RiInfestadistica
	response := r.SQLClient.Raw("select '1010400' as CodigoPartida, mp.codigo_bcra as MedioPago, c.codigo_bcra as EsquemaPago, count(c.id) as CantOperaciones, sum(m.monto) as MontoTotal from movimientos as m inner join pagointentos as pi on pi.id = m.pagointentos_id inner join mediopagos as mp on mp.id = pi.mediopagos_id inner join channels as c on mp.channels_id = c.id where tipo = 'D' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by mp.id", fechaInicio, fechaFin).Scan(&ri1010400)

	if response.Error != nil {
		erro = response.Error
		return
	}

	for i := range ri1010400 {
		ri = append(ri, ri1010400[i])
	}

	//Total de Clientes PJ involucrados en transacciones
	var ri6011000 int64
	response = r.SQLClient.Raw("select count(c.clientes_id) as ri6011000 from (select cue.clientes_id From cuentas as cue inner join clientes as cli on cli.id = cue. clientes_id inner join (select m.cuentas_id from movimientos as m where cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m on m.cuentas_id = cue.id where cli.personeria = 'J' group by cue.clientes_id) as c", fechaInicio, fechaFin).Scan(&ri6011000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6011000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6011000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6011000, MontoTotal: "NULO"})
	}

	//Total de Clientes PF involucrados en transacciones
	var ri6012000 int64
	response = r.SQLClient.Raw("select count(c.clientes_id) as ri6012000 from (select cue.clientes_id From cuentas as cue inner join clientes as cli on cli.id = cue. clientes_id inner join (select m.cuentas_id from movimientos as m where cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m on m.cuentas_id = cue.id where cli.personeria = 'F' group by cue.clientes_id) as c", fechaInicio, fechaFin).Scan(&ri6012000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6012000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6012000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6012000, MontoTotal: "NULO"})
	}

	//Total de Clientes PJ involucrados en transferencias
	var ri6013000 int64
	response = r.SQLClient.Raw("select count(c.clientes_id) as ri6013000 from (select cue.clientes_id From cuentas as cue inner join clientes as cli on cli.id = cue.clientes_id inner join (select m.cuentas_id from movimientos as m inner join pagointentos pi on pi.id=m.pagointentos_id inner join mediopagos as mp on mp.id=pi.mediopagos_id where mp.mediopago = 'Debin' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m on m.cuentas_id = cue.id where cli.personeria = 'J' group by cue.clientes_id) as c", fechaInicio, fechaFin).Scan(&ri6013000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6013000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6013000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6013000, MontoTotal: "NULO"})
	}

	//Total de Clientes PF involucrados en transferencias
	var ri6014000 int64
	response = r.SQLClient.Raw("select count(c.clientes_id) as ri6014000 from (select cue.clientes_id From cuentas as cue inner join clientes as cli on cli.id = cue.clientes_id inner join (select m.cuentas_id from movimientos as m inner join pagointentos pi on pi.id=m.pagointentos_id inner join mediopagos as mp on mp.id=pi.mediopagos_id where mp.mediopago = 'Debin' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m on m.cuentas_id = cue.id where cli.personeria = 'F' group by cue.clientes_id) as c", fechaInicio, fechaFin).Scan(&ri6014000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6014000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6014000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6014000, MontoTotal: "NULO"})
	}

	//Total de cuentas de pago PJ involucradas en transaciones
	var ri6021000 int64
	response = r.SQLClient.Raw("select count(m.cuentas_id) as ri6021000 from (select m.cuentas_id from movimientos as m inner join cuentas as c on c.id = m.cuentas_id inner join clientes as cli on cli.id=c.clientes_id where cli.personeria = 'J' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m", fechaInicio, fechaFin).Scan(&ri6021000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6021000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6021000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6021000, MontoTotal: "NULO"})
	}

	//Total de cuentas de pago PF involucradas en transaciones
	var ri6022000 int64
	response = r.SQLClient.Raw("select count(m.cuentas_id) as ri6022000 from (select m.cuentas_id from movimientos as m inner join cuentas as c on c.id = m.cuentas_id inner join clientes as cli on cli.id=c.clientes_id where cli.personeria = 'F' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m", fechaInicio, fechaFin).Scan(&ri6021000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6022000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6022000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6022000, MontoTotal: "NULO"})
	}

	//Total de cuentas de pago PJ involucradas en transferencias
	var ri6023000 int64
	response = r.SQLClient.Raw("select count(m.cuentas_id) as ri6023000 from (select m.cuentas_id from movimientos as m inner join cuentas as c on c.id = m.cuentas_id inner join clientes as cli on cli.id=c.clientes_id inner join pagointentos pi on pi.id=m.pagointentos_id inner join mediopagos as mp on mp.id=pi.mediopagos_id where mp.mediopago = 'Debin' and cli.personeria = 'J' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m", fechaInicio, fechaFin).Scan(&ri6023000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6023000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6023000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6023000, MontoTotal: "NULO"})
	}

	//Total de cuentas de pago PF involucradas en transferencias
	var ri6024000 int64
	response = r.SQLClient.Raw("select count(m.cuentas_id) as ri6024000 from (select m.cuentas_id from movimientos as m inner join cuentas as c on c.id = m.cuentas_id inner join clientes as cli on cli.id=c.clientes_id inner join pagointentos pi on pi.id=m.pagointentos_id inner join mediopagos as mp on mp.id=pi.mediopagos_id where mp.mediopago = 'Debin' and cli.personeria = 'F' and cast(m.created_at as date) between cast(? as date) AND cast(? as date) group by m.cuentas_id) as m", fechaInicio, fechaFin).Scan(&ri6024000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6024000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6024000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6024000, MontoTotal: "NULO"})
	}

	//Total de clientes PJ
	var ri6031000 int64
	response = r.SQLClient.Raw("SELECT count(id) as ri6031000 FROM pasarela.clientes where personeria = 'J' and deleted_at is null").Scan(&ri6031000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6031000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6031000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6031000, MontoTotal: "NULO"})
	}

	//Total Clientes PF

	var ri6032000 int64
	response = r.SQLClient.Raw("SELECT count(id) as ri6032000 FROM pasarela.clientes where personeria = 'F' and deleted_at is null").Scan(&ri6032000)

	if response.Error != nil {
		erro = response.Error
		return
	}
	if ri6032000 > 0 {
		ri = append(ri, ribcradtos.RiInfestadistica{CodigoPartida: "6032000", MedioPago: "NULO", EsquemaPago: "NULO", CantOperaciones: ri6032000, MontoTotal: "NULO"})
	}

	return
}
