# solicitud token pago simple y offline

***
## Error al solicitar permiso de pago(el tipo de pago enviado no es valodo)
1. solicita un permiso de pago SolicitarToken(request prismadtos.StructToken)
2. verifica el tipo de pago 
3. retorna ERRR_TIPO_PAGO
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
    else
    Note over sa: tipo pago offline
    else
    Note over sa: tipo pago no valido
    sa-->>co: retorna: ERRR_TIPO_PAGO
    end
    deactivate sa
```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/00-solicitud_permiso_de_pago.md