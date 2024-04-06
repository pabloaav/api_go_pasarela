> # New Pago 

## Error validación new pago
1. El proceso se inicia al llamar al servicio NewPago
2. Se validan los datos enviados por parámetro en dtos.PageRequest.Validar()
3. Devuelvo el Error de la validación

***


```mermaid
sequenceDiagram;
    participant S as Servicio
    participant PR as dtos.PagoRequest
    participant R as Repository
    participant UID as Commons
    S ->> PR: request
    PR ->> PR: validar
    PR -->> S: Error
    S ->> S: return Error en la validación
    
```    