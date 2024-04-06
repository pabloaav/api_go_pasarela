# Pagos pendientes expirados

## Caso de succeso no se encontraron pagos expirados
1. Busca la configuración TIEMPO_EXPIRACION_PAGOS
2. Busca el estado de pago pendiente
3. Busca los pagos expirados
4. No se encontraron pagos expirados
5. Finaliza el proceso
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
    GP -->> ME: Lista Vacia
    ME -->>BS: Finaliza proceso
```
