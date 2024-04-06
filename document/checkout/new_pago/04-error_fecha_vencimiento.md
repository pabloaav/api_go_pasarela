> # New Pago 

## Error en fechas de vencimientos new pago
1. El proceso se inicia al llamar al servicio NewPago
2. Se validan los datos enviados por parámetro en dtos.PageRequest.Validar()
3. Se busca el registro de cliente en el repositorio repository.GetClienteByApiKey(apiKey)
4. Busco el id del tipo de pago entre los Pagotipo del objeto Cliente
5. Convierto las Fechas de vencimiento de string a formato time
6. Devuelvo el Error en fecha de vencimiento

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
    S ->> S: return Error en fecha de vencimiento
    
```    