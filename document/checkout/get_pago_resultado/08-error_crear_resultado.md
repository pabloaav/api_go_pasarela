> # Get PagoResultado

## error al procesar pago con metodo externo
1. El proceso se inicia al llamar al servicio GetPagoResultado
2. Se valida que se envíen correctamente los parámetros
3. Consulta datos del Channel en el repositorio 
4. Consulta datos del medio de pago al repositorio
5. Consulta datos del pago mediante el uuid
6. Consulta datos del tipo de pago
7. Consulta datos de la cuenta 
8. Calcula importe a pagar según fechas de vencimiento
9. Convierte a entero el importe 
10. Obtiene el metodo de pago segun channelID
11. Ejecuta el crear resultado del metodo obtenido
12. Devuelvo el error si falla

***

```mermaid
sequenceDiagram;
    participant S as Servicio
    participant PR as dtos.ResultadoRequest
    participant R as Repository
    participant PF as PaymentFactory
    participant PM as PaymentMethod
    S ->> PR: request
    PR ->> PR: validar
    PR -->> S: ok
    S ->> R: channel
    R ->> R: GetChannelByName
    R -->> S: Channel
    S ->> R: ChannelID, CardBrand
    R ->> R: GetMedioPago
    R -->> S: Mediopago
    S ->> R: uuid
    R ->> R: GetPagoByUuid
    R -->> S: Pago
    S ->> R: PagostipoID
    R ->> R: GetPagotipoById
    R -->> S: Pagotipo
    S ->> R: CuentasID
    R ->> R: GetCuentaByID
    R -->> S: Cuenta
    S ->> S: Compara fechas
    S ->> S: Convierte importe a entero
    S ->> PF: ChannelID
    PF -->> S: PaymentMethod
    S ->> PM: request, pago, cuenta
    PM ->> PM: CreateResultado
    PM -->> S: Error al procesar pago
    S ->> S: return Error
        
```