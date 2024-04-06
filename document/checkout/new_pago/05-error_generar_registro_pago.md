> # New Pago 

## Error al generar registro de pago new pago
1. El proceso se inicia al llamar al servicio NewPago
2. Se validan los datos enviados por parámetro en dtos.PageRequest.Validar()
3. Se busca el registro de cliente en el repositorio repository.GetClienteByApiKey(apiKey)
4. Busco el id del tipo de pago entre los Pagotipo del objeto Cliente
5. Convierto las Fechas de vencimiento de string a formato time
6. Genero un id unico con el servicio commons.NewUuid
7. Genero el nuevo objeto Pago
8. Envío el objeto pago al repositorio para que lo almacene en la base de datos
9. Devuelvo el Error

***


```mermaid
sequenceDiagram;
    participant S as Servicio
    participant PR as dtos.PagoRequest
    participant R as Repository
    participant UID as Commons
    S ->> PR: request
    PR ->> PR: validar
    PR -->> S: 
    S ->> R: apiKey
    R ->> R: GetClienteByApikey
    R -->> S: Cliente
    S ->> S: busco tipoPagoID
    S ->> S: convierto Fechas de Vencimiento
    S ->> UID: 
    UID ->> UID: NewUUID
    UID -->> S: UUID
    S ->> S: Genero Objeto Pago
    S ->> R: Pago
    R ->> R: CreatePago
    R -->> S: Error no se pudo generar registro pago
    S ->> S: return Error
    
```    