# ejecucion de pago simple y offline

***
## Error al parsear y validar URL (para relaizar una llamada al servicio de prisma)
1. solicita ejecucion de pago Payments(request prismadtos.StructPayments)
2. verifica el tipo de pago "si es tipo pago offline"
3. valida datos recibidos objetoRequest.Validar()
4. relaiza llamda al repositorio remoto PostEjecutarPagoOffLine(&objetoRequest)
5. serializa la estructura de dato recibida
6. parsea y valida la ruta de prisma 
7. retorna "Error al crear base url" + el codigo de errror
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
    else
    Note over sa: tipo pago offline
    sa-->>sa: valida datos recibidos objetoRequest.Validar()
    sa->>rrp: relaiza llamda al repositorio remoto PostEjecutarPagoOffLine(&objetoRequest)
    activate rrp
    rrp-->>rrp:serialza la estructura de dato recibida
    rrp-->>rrp:parsea y valida la ruta de prisma
    rrp-->>sa:retorna "Error al crear base url" + el codigo de errror
    deactivate rrp
    sa-->>co: retorna: ERROR
    else
    Note over sa: tipo pago no valido
    end
    deactivate sa
```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/ejecuciondepago/00-ejecucion_de_pago.md