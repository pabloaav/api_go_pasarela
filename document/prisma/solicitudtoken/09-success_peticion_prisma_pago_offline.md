# solicitud token pago simple y offline

***
## Success al solicitar una peticion de token a prisma (prisma devuelve respuesta exitosa y tonken)
1. solicita un permiso de pago SolicitarToken(request prismadtos.StructToken)
2. verifica el tipo de pago simple
3. valida datos recibidos
5. relaiza llamda al repositorio remoto PostSolicitudTokenPago(&objetoRequest)
6. serializa la estructura de dato recibida
7. parsea y valida la ruta de prisma, respuesta exitosa 
8. concatena a la URL con la URI para solicitar token
9. construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
10. construye un header y se le agrega a la request buildHeaderDefault(reuqest, PUBLIC_APIKEY_PRISMA)
11. realizamos peticion a prisma HTTPClient.Do(req): retorna repuesta exito
12. verifica si retorna un ERROR
    - 12.1. no retorna error
13. cerrar la conexion, defer resp.Body.Close()
14. se valida el codigo status sea distinto a 201, retorna si es igual 201
15. deserializa la respuesta 
16. retorna la respuesta 
17. al servicio, retorna respuesta 
18. al checkout, retorna respuesta 
***
```mermaid
sequenceDiagram;
    participant co as CheckOut
    participant sa as ServiceTelcoPrisma
    participant rrp as RepositoryRemotePrisma
    participant p as Prisma
    co ->> sa: SolicitarToken(request prismadtos.StructToken)
    activate sa
    alt verifica el tipo pago
    Note over sa: tipo pago simple
    else
    Note over sa: tipo pago offline
     sa-->>sa: valida estructura de dato recibida objetoRequest.ValidarSolicitudTokenOfflineRequest():retorna nil
    sa->>rrp:  relaiza llamda al repositorio remoto PostSolicitarTokenOffLine(&objetoRequest)
    activate rrp
    rrp-->>rrp:serialza la estructura de dato recibida
    rrp-->>rrp:parsea y valida la ruta de prisma: respuesta exitosa
    rrp-->>rrp:concatena a la URL con la URI para solicitar token
    rrp-->>rrp:construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
    rrp-->>rrp:construye un header y se le agrega a la request buildHeaderDefault(reuqest, PUBLIC_OFFLINE_APIKEY_PRISMA)
    rrp->>p:realizamos peticion a prisma HTTPClient.Do(req)
    activate p
    p-->>rrp:retorna objeto response + token
    deactivate p
    alt
    Note over rrp: no responde error en la repuesta
    end 
    Note over rrp: responde un objeto response con el token
    rrp-->>rrp:cerrar la conexion, defer resp.Body.Close()
    alt
    Note over rrp: codigo status es distinto de 201: retorna codigo 201
    end
    rrp-->>rrp:cerrar la conexion, defer resp.Body.Close()
    rrp-->>rrp:deserializa la respuesta 
    rrp-->>sa:retorna la respuesta
    deactivate rrp
    sa-->>co:retorna la respuesta 
    else
    Note over sa: tipo pago no valido
    end
    deactivate sa
```
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/00-solicitud_permiso_de_pago.md