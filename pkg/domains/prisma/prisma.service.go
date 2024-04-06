package prisma

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"

	//prismaCierreLote "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaCierreLote"
	prismaOperaciones "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismaOperaciones"
	prismainforme "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismainformes"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos/prismadtos/prismatransacciones"
)

type Service interface {

	//////////////////trannsacciones simples y offline////////////////////////////////////////////////////
	/*
		CheckService Chequear el servicio de prisma
	*/
	CheckService() (estado bool, err error)

	/*
		SolicitarToken solicita un toquen de pago para una transaccion que puede ser simple o offline
		se bebe enviar el tipo de pago los valores pueden ser: simple o offline
	*/
	SolicitarToken(request prismadtos.StructToken) (response interface{}, erro error)

	/*
		payments  por medio de esta funcion se realiza la ejecucion del pago que puede ser simple o offline
		se bebe enviar el tipo de pago los valores pueden ser: simple o offline
	*/
	Payments(request prismadtos.StructPayments) (response interface{}, erro error)

	///////////////////////////////////////INFORMES DE PAGOS//////////////////////////////////////////////////////////////////
	/*
		GetInformePago permite obtener infromacion de un pago, recibe id de pago "payment_id"
	*/
	GetInformePago(pagoId string) (response *prismainforme.UnPagoResponse, erro error)

	/*
		ListarPagos permite obtener una lista de pago enviando una fecha desde y una fecha hasta
	*/
	ListarPagosPorFecha(requestPago prismainforme.ListaPagosRequest) (listaPago []prismainforme.Result, erro error)

	/*
		servicio para obtener todos los pagos cuyo estadoid sea igual al valor pasado por paramentro y los pagos intentos realocionado que
		coincida con el canal tambien pasado por parametros
	*/
	ListarPagosService(estadoPago int, channel string) (pagointentos []entities.Pagointento, erro error)

	///////////////////////////////////////OPERACIONES SOBRE TRANSACCIONES//////////////////////////////////////////////////////////////////
	/*
		PostAnulacionDevolucionTotalPago se realiza la solicitud de anulación o devolución total de pago
	*/
	PostAnulacionDevolucionTotalPago(ExternalId string) (response *prismaOperaciones.SolicitudAnulacionDevolucionResponse, erro error)

	/*
		PostSolicitudAnulacionDevolucionPagoParcial se realiza la solicitud de anulación o devolución parcial de pago a prisma
		se le debe enviar el payment_id y el monto
	*/
	PostSolicitudAnulacionDevolucionPagoParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.SolicitudAnulacionDevolucionPagoParcialResponse, erro error)

	/*
		DelAnulacionDevolucionPagoTotalParcial servicio para eliminar una solicitud de anulacion o devolucion de pago total o parcial
	*/
	DelAnulacionDevolucionPagoTotalParcial(params prismaOperaciones.ParamsPagoParcialTotalService) (response *prismaOperaciones.DeletSolicitudPagoResponse, erro error)
}

// prisma variable que va a manejar la instancia del servicio
var prisma *service

type service struct {
	remoteRepository RepositoryRemotePrisma
	repository       Repository
	commonsService   commons.Commons
	// adminService     administracion.Service
}

func NewService(rm RepositoryRemotePrisma, r Repository, c commons.Commons) Service { //, admin administracion.Service
	// al instanciar el servicio lo almaceno en la variable apilink
	prisma = &service{
		remoteRepository: rm,
		repository:       r,
		commonsService:   c,
		// adminService:     admin,
	}
	return prisma
}

// Resolve devuelve la instancia antes creada
func Resolve() *service {
	return prisma
}
