# solicitud token pago simple y offline

## estructura de datos para realizar una solicitud de token
- el servicio realiza una solisitud de token al prisma.
    * si la solicitudi token es por un pago con tarjeta debe enviar enviar la siguiente informacion
        + Card
        + TypePay = "simple"
    * si la solicitudi token es por un pago offline debe enviar enviar la siguiente informacion
        + DataOffline
        + TypePay = "offline"
- todo esto se realiza llamando al mismos servicio "SolicitarToken", este devuelve una interface o el error producido 
***
	Card        Card                
	DataOffline OfflineTokenRequest 
	TypePay     ["simple", "offline"] 
***
### Card
    CardNumber               string                   
	CardExpirationMonth      string                   
	CardExpirationYear       string                   
	SecurityCode             string                  
	CardHolderName           string                  
	CardHolderIdentification CardHolderIdentification {
                                        TypeDni   ["DNI", "CI", "LE", "LC"] 
	                                    NumberDni string            
                             }
***
### DataOffline
    Customer DataCustomer{
        Identification IdentificationCustomer {
                                Type   ["DNI", "CI", "LE", "LC"] 
	                            Number string            
                        }
	    Name           string                 
    }

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
- [Error al solicitar permiso de pago(el tipo de pago enviado no es valodo)][URL-Error1S]
## Solicitud de Token Pago Simple
- [Error al validar la estructura de dato recibida (para un pago simble)][URL-Error2S]
- [Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)][URL-Error3S]
- [Error al solicitar una peticion de token a prisma (prisma devuelve un error en la respuesta)][URL-Error4S]
- [Success al solicitar una peticion de token a prisma (prisma devuelve respuesta exitosa y tonken)][URL-SuccessS]
## Solicitud de Token Pago Off-Line

- [Error al validar la estructura de dato recibida (para un pago offline)][URL-Error1O] 
- [Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)][URL-Error2O]
- [Error al solicitar una peticion de token a prisma (prisma devuelve un error en la respuesta)][URL-Error3O]
- [Success al solicitar una peticion de token a prisma (prisma devuelve respuesta exitosa y tonken)][URL-SuccessO]







[URL-Error1S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/01-error_tipo_de_pago.md 
[URL-Error2S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/02-error_datos_recibidos_pago_simple.md
[URL-Error3S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/03-error_repository_remoto_armado_url_pago_simple.md
[URL-Error4S]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/04-error_peticion_prisma_pago_simple.md
[URL-SuccessS]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/05-success_peticion_prisma_pago_simple.md


[URL-Error1O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/06-error_datos_recibidos_pago_offline.md
[URL-Error2O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/07-error_repository_remoto_armado_url_pago_offline.md
[URL-Error3O]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/08-error_peticion_prisma_pago_offline.md
[URL-SuccessO]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/09-success_peticion_prisma_pago_offline.md
