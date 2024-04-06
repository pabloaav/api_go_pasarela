> # Get PagoResultado

## error al obtener el medio de pago
1. El proceso se inicia al llamar al servicio GetPagoResultado
2. Se valida que se envíen correctamente los parámetros
3. Consulto datos del Channel en el repositorio 
4. Consulto datos del medio de pago al repositorio
5. Devuelvo el error si no se encuentra

***

```mermaid
sequenceDiagram;
    participant S as Servicio
    participant PR as dtos.ResultadoRequest
    participant R as Repository
    S ->> PR: request
    PR ->> PR: validar
    PR -->> S: ok
    S ->> R: channel
    R ->> R: GetChannelByName
    R -->> S: Channel
    S ->> R: ChannelID, CardBrand
    R ->> R: GetMedioPago
    R -->> S: Error al obtener medio de pago
    S ->> S: return Error
        
```