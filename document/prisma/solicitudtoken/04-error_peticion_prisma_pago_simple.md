# solicitud token pago simple y offline

***
## Error al solicitar una peticion de token a prisma (prisma devuelve un error en la respuesta)
1. solicita un permiso de pago SolicitarToken(request prismadtos.StructToken)
2. verifica el tipo de pago simple
3. valida datos recibidos
5. relaiza llamda al repositorio remoto PostSolicitudTokenPago(&objetoRequest)
6. serializa la estructura de dato recibida
7. parsea y valida la ruta de prisma, respuesta exitosa 
8. concatena a la URL con la URI para solicitar token
9. construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
10. construye un header y se le agrega a la request buildHeaderDefault(reuqest, PUBLIC_APIKEY_PRISMA)
11. realizamos peticion a prisma HTTPClient.Do(req): retorna ERROR
12. verifica si retorna un ERROR
    - 12.1. genera log de error, logs.Error("error al solicitar token: " + err.Error())
13. cerrar la conexion, defer resp.Body.Close()
14. se valida el codigo sea distinto a 201, si es distinto
    - 14.1. deserializa la respuesta 
    - 14.2. genera log de error, logs.Error("Error en la peticion " + err.Error())
15. al servicio, retorna ERROR 
16. al checkout, retorna ERROR 
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
    sa-->>sa: valida estructura de dato recibida Card.Validar():retorna nil
    sa->>rrp:  relaiza llamda al repositorio remoto PostSolicitudTokenPago(&objetoRequest)
    activate rrp
    rrp-->>rrp:serialza la estructura de dato recibida
    rrp-->>rrp:parsea y valida la ruta de prisma: respuesta exitosa
    rrp-->>rrp:concatena a la URL con la URI para solicitar token
    rrp-->>rrp:construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
    rrp-->>rrp:construye un header y se le agrega a la request buildHeaderDefault(reuqest, PUBLIC_APIKEY_PRISMA)
    rrp->>p:realizamos peticion a prisma HTTPClient.Do(req)
    activate p
    p-->>rrp:retorna "error al solicitar token:" + el codigo de errror
    deactivate p
    alt
    Note over rrp: responde error en la repuesta
    rrp-->>rrp:genera log de error, logs.Error("error al solicitar token: " + err.Error())
    else
    Note over rrp: responde un objeto response con el token
    end 
    rrp-->>rrp:cerrar la conexion, defer resp.Body.Close()
    alt
    Note over rrp: codigo status es distinto de 201
    rrp-->>rrp:deserializa la respuesta 
    rrp-->>rrp:genera log de error, logs.Error("Error en la peticion " + err.Error())
    end
    rrp-->>sa:retorna ERROR 
    deactivate rrp
    sa-->>co:retorna ERROR 
    else
    Note over sa: tipo pago offline
    else
    Note over sa: tipo pago no valido
    end
    deactivate sa
```
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/solicitudtoken/00-solicitud_permiso_de_pago.md