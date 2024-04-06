# Pagos pendientes expirados

## Error al cargar el estado expirado
1. Busca la configuración TIEMPO_EXPIRACION_PAGOS
2. Busca el estado de pago pendiente
3. Busca los pagos expirados
4. Busca el estado Expirado
5. ERROR_PAGO_ESTADO
***


```mermaid
sequenceDiagram;
    participant BS as BackgroudServices
    participant ME as ModificarEstadoPagosExpirados
    participant GC as GetConfiguracion
    participant GPE as GetPagoEstado
    participant GP as GetPago
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
      Note over GPE: Error al Buscar
    GPE -->>ME: ERROR_PAGO_ESTADO
    ME -->>BS: ERROR_PAGO_ESTADO
```
