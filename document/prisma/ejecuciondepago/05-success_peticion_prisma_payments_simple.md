# ejecucion de pago simple y offline
***
## Success al srealizar peticion de ejecucion de pago simple prisma (prisma responde como respuesta exitosa un json con la informacion del pago realizado)
1. solicita ejecucion de pago Payments(request prismadtos.StructPayments)
2. verifica el tipo de pago "si es tipo pago simple"
3. valida datos recibidos objetoRequest.ValidarProcesoPagoRequest()
4. relaiza llamda al repositorio remoto PostEjecutarPago(&objetoRequest)
5. serializa la estructura de dato recibida
6. parsea y valida la ruta de prisma , respuesta exitosa 
7. concatena la URL con la URI para solicitar ejecucion de pago simple
8. construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
9. construye un header y se le agrega a la request buildHeaderDefault(reuqest, PRIVATE_APIKEY_PRISMA)
10. realizamos peticion a prisma HTTPClient.Do(req): retorna un json con la informacion del pago realizado
11. verifica si retorna un ERROR: No existe Error
12. cerrar la conexion, defer resp.Body.Close()
13. valida que el codigo de estado de la peticion sea distinto de 201: codigo de estado igual 201
14. deserializa la respuesta y captura informacion del pago realizado
15. al servicio, retorna UnPago
14. al checkout, retorna UnPago 

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
             sa->>rrp:  relaiza llamda al repositorio remoto PostEjecutarPago(&objetoRequest)
             activate rrp
                rrp-->>rrp:serialza la estructura de dato recibida
                rrp-->>rrp: parsea y valida la ruta de prisma , respuesta exitosa 
                rrp-->>rrp: concatena la URL con la URI para solicitar ejecucion de pago simple
                rrp-->>rrp: construye una request http.NewRequest(metodo, url, bytes.NewBuffer(estructura serealizada))
                rrp-->>rrp: construye un header y se le agrega a la request buildHeaderDefault(reuqest, PRIVATE_APIKEY_PRISMA)
                rrp->>p: realizamos peticion a prisma HTTPClient.Do(req)
                activate p
                     p-->>rrp: retorna json_con_la_informacion_del_pago_realizado
                deactivate p
                alt verifica si retorna un ERROR
                    Note over rrp: No existe Error
                end 
                alt valida que el codigo de estado de la peticion sea distinto de 201
                    Note over rrp:codigo de estado igual 201
                end
                rrp-->>rrp: deserializa la respuesta y captura informacion de pago
                rrp-->>sa:retorna UnPago
             deactivate rrp
             sa-->>co: retorna UnPago
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