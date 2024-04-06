# Pagos pendientes expirados

## Error al crear configuración
1. Busca la configuración TIEMPO_EXPIRACION_PAGOS
2. En caso de que no encuentre la configuración la crea
3. ERROR_CREAR_CONFIGURACIONES
***


```mermaid
sequenceDiagram;
    participant BS as BackgroudServices
    participant ME as ModificarEstadoPagosExpirados
    participant GC as GetConfiguracion
    participant CC as CreateConfiguracion
    BS ->> ME: Inicio proceso
    ME ->> GC: Filtro Configuración
    GC ->> GC: Cargar Configuración
    GC -->> ME: Configuracion
    Note over ME: Configuracion.Id = 0
    ME ->> CC: Nueva Configuración
    CC ->> CC: Crear Configuración
    Note over CC: Error al Crear
    CC -->>ME: ERROR_CREAR_CONFIGURACIONES
    ME -->>BS: ERROR_CREAR_CONFIGURACIONES
```
