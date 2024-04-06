# ejecucion de pago simple y offline

***
## Error al validar la estructura de dato recibida (para ejecucion de pago simble)
1. solicita ejecucion de pago Payments(request prismadtos.StructPayments)
2. verifica el tipo de pago "si es tipo pago simple"
3. valida datos recibidos objetoRequest.ValidarProcesoPagoRequest()
4. retorna ERRROR

## posible constantes de ERRROR que se puede recibir al validar los datos para una ejecucion de pago simple
    - ERROR_SITE_TRANSACTION_ID
    - ERROR_TOKEN_PAGO
    - ERROR_BIN
    - ERROR_AMOUNT
    - ERROR_INSTALLMENTS
    - ERROR_PAYMENT_TYPE
    - ERROR_NOMBRE_ESTABLECIMIENTO
    - ERROR_EMAIL
***
```mermaid
sequenceDiagram;
    participant co as CheckOut
    participant sa as ServiceTelcoPrisma
    participant rrp as RepositoryRemotePrisma
    co ->> sa: Payments(request prismadtos.StructPayments)
    activate sa
    alt verifica el tipo pago
    Note over sa: tipo pago simple
    sa-->>sa: valida datos recibidos objetoRequest.ValidarProcesoPagoRequest()
    sa-->>co: retorna: ERROR
    else
    Note over sa: tipo pago offline
    else
    Note over sa: tipo pago no valido
    end
    deactivate sa
```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/00-ejecucion_de_pago.md