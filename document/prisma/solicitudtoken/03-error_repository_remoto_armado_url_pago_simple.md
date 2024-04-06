# solicitud token pago simple y offline

***
## Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)
1. solicita un permiso de pago SolicitarToken(request prismadtos.StructToken)
2. verifica el tipo de pago simple
3. valida datos recibidos
4. relaiza llamda al repositorio remoto PostSolicitudTokenPago(&objetoRequest)
5. serializa la estructura de dato recibida
6. parsea y valida la ruta de prisma 
7. retorna "Error al crear base url" + el codigo de errror
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
    sa-->>sa: valida estructura de dato recibida Card.Validar():retorna nil
    sa->>rrp:  relaiza llamda al repositorio remoto PostSolicitudTokenPago(&objetoRequest)
    activate rrp
    rrp-->>rrp:serialza la estructura de dato recibida
    rrp-->>rrp:parsea y valida la ruta de prisma
    rrp-->>sa:retorna "Error al crear base url" + el codigo de errror
    sa-->>co:retorna ERROR 
    deactivate rrp
    else
    Note over sa: tipo pago offline
    else
    Note over sa: tipo pago no valido
    end
    deactivate sa
```
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/00-solicitud_permiso_de_pago.md