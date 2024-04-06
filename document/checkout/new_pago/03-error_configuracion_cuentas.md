> # New Pago 

## error en configuracion de cuentas new pago
1. El proceso se inicia al llamar al servicio NewPago
2. Se validan los datos enviados por parámetro en dtos.PageRequest.Validar()
3. Se busca el registro de cliente en el repositorio repository.GetClienteByApiKey(apiKey)
4. Busco el id del tipo de pago entre los Pagotipo del objeto Cliente
5. Devuelvo el Error en la configuración de cuentas

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
    S ->> S: return Error en la configuración de cuentas
    
```    