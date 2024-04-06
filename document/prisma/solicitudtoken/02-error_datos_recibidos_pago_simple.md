# solicitud token pago simple y offline

***
## Error al validar la estructura de dato recibida (para un pago simble)
1. solicita un permiso de pago SolicitarToken(request prismadtos.StructToken)
2. verifica el tipo de pago simple
3. valida datos recibidos
4. retorna ERRROR
***
## posible constantes de ERRROR que se puede recibir al validar los datos para una solicitud de pago simple
    - ERROR_ESTRUCTURA_INCORRECTA
    - ERROR_NUMBER_CARD
    - ERROR_NUMBER_CARD
    - ERROR_DATE_CARD
    - ERROR_DATE_CARD
    - ERROR_HOLDER_NAME
    - ERROR_TIPO_DOCUMENTO
***

```mermaid
sequenceDiagram;
    participant co as CheckOut
    participant sa as ServiceTelcoPrisma
    participant rrp as RepositoryRemotePrisma
    co ->> sa: SolicitarToken(request prismadtos.StructToken)
    activate sa
    alt verifica el tipo pago
    Note over sa: tipo pago simple
    sa-->>sa: valida estructura de dato recibida Card.Validar()
    sa-->>co: retorna: ERROR
    else
    Note over sa: tipo pago offline
    else
    Note over sa: tipo pago no valido
    
    end
    deactivate sa
```
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/00-solicitud_permiso_de_pago.md