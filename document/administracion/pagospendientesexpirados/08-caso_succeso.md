# Pagos pendientes expirados

## Caso de succeso
1. Busca la configuración TIEMPO_EXPIRACION_PAGOS
2. Busca el estado de pago pendiente
3. Busca los pagos expirados
4. Busca el estado Expirado
5. Modifica el estado de los pagos y crea los PagoEstadosLogs
6. Finaliza proceso
***


```mermaid
sequenceDiagram;
    participant BS as BackgroudServices
    participant ME as ModificarEstadoPagosExpirados
    participant GC as GetConfiguracion
    participant GPE as GetPagoEstado
    participant GP as GetPago
    participant UE as UpdateEstadoPagos
    BS ->> ME: Inicio proceso
    ME ->> GC: Filtro Configuración
    GC ->> GC: Cargar Configuración
    GC -->> ME: Configuracion
    ME ->> GPE: FiltroPagoEstado
    GPE ->> GPE: Cargar PagoEstado
    GPE -->>ME: Pago Estado
    ME ->> GP: FiltroPago
    GP ->> GP: Carga pagos expirados
    GP -->> ME: pagos expirados
    ME ->> GPE: FiltroPagoEstado
    GPE ->> GPE: Carga Estado Expirado
    GPE -->> ME: Pago Estado Expirado
    ME ->> UE: pagos expirados | estado expirado
    UE ->> UE: Modifica Pagos Expirados
    UE -->> ME: Finaliza Proceso
    ME -->>BS: Finaliza Proceso
```
