# ejecucion de pago simple y offline

## estructura de datos para realizar una ejecucion de pago
- el servicio realiza una ejecucion de pago en prisma.
    * si la ejecucion de pago, es por un pago con tarjeta debe enviar enviar la siguiente informacion
        + PagoSimple
        + TypePay = "simple"
    * si la ejecucion de pago, es por un pago offline debe enviar enviar la siguiente informacion
        + PagoOffline
        + TypePay = "offline"
- todo esto se realiza llamando al mismos servicio "Payments", este devuelve una interface o el error producido 
***
	PagoSimple  PaymentsSimpleRequest                
	PagoOffline PaymentsOfflineRequest 
	TypePay     ["simple", "offline"] 
***

### PagoSimple
    Customerid        Customerid{
                            ID string
                    }
    SiteTransactionID string
    SiteID            int64
    Token             string
    PaymentMethodID   int64
    Bin               string
    Amount            int64
    Currency          EnumTipoMoneda["ARS", "USD"]
    Installments      int64
    Description       string
    PaymentType       EnumPaymentType["single", "distributed"]
    EstablishmentName strin
    Customeremail     Customeremail{
                            Email string
                        }
    SubPayments       []interface{}
***
### PagoOffline
    Customer   DataCustomer{
                     Name           string
                     Identification IdentificationCustomer{
                                Type   EnumTipoDocumento["DNI", "CI", "LE", "LC"]
                                Number string
                     }
               }
    SiteTransactionID string
    Token             string
    PaymentMethodID   int64
    Currency          EnumTipoMoneda["ARS", "USD"]
    PaymentType       string
    Email             string
    InvoiceExpiration string
    CodP3             string
    CodP4             string
    Client            string
    Surcharge         int64
    PaymentMode       string
***
### Estructura del Error
ErrorEstructura{
    	ErrorType        string            
	    ValidationErrors []ValidationError{
	                                Code   string
	                                Param  string
	                                Status string
                            } 
	    Message          string            
	    Code             string            
}
***
- nota: los elemento que interactuan son:
    * checkout
    * Servicio telco 
    * Repository Remoto Prisma
    * Prisma
***
# Casos 

- [Error al solicitar la ejecucion de pago(el tipo de pago enviado no es valido)][URL-Error1S]
## Ejecución de Pago Simple
- [Error al validar la estructura de dato recibida (para ejecucion de pago simble)][URL-Error2S]
- [Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)][URL-Error3S]
- [Error en peticion ejecucion de pago simple (para relaizar una llamada al servicio de prisma)][URL-Error4S]
- [Success al srealizar peticion de ejecucion de pago simple prisma (prisma responde como respuesta exitosa un json con la informacion del pago realizado)][URL-SuccessS]
## Ejecución de Pago Off-Line

- [Error al validar la estructura de dato recibida (para ejecucion de pago offline)][URL-Error1O] 
- [Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)][URL-Error2O]
- [Error en peticion ejecucion de pago offline (para relaizar una llamada al servicio de prisma)][URL-Error3O]
- [Success al srealizar peticion de ejecucion de pago offline prisma (prisma responde como respuesta exitosa un json con la informacion del pago realizado)][URL-SuccessO]

[URL-Error1S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/01-error_tipo_de_pago.md
[URL-Error2S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/02-error_datos_recibidos_payments_simple.md
[URL-Error3S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/03-error_repository_remoto_armado_url_payments_simple.md
[URL-Error4S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/03-error_repository_remoto_armado_url_payments_simple.md
[URL-SuccessS]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/05-success_peticion_prisma_payments_simple.md


[URL-Error1O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/02-error_datos_recibidos_payments_offline.md
[URL-Error2O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/03-error_repository_remoto_armado_url_payments_offline.md
[URL-Error3O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/04-error_peticion_prisma_payments_offline.md
[URL-SuccessO]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/05-success_peticion_prisma_payments_offline.md
