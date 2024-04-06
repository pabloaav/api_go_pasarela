> # Cierre de Lote ApiLink 

## Error al cargar los pagos
1. Busca el pago estado processing (estado inicial para los debines)
2. Busca el canal debin
3. Busco los pagos que pertenecen al pago estado punto 1
4. ERROR_PAGO | ERROR_PAGO_LISTA

***


```mermaid
sequenceDiagram;
    participant B as BuildMovimiento
    participant BC as BuildCierreLoteApiLinkService
    participant GPE as GetPagoEstado
    participant GC as GetChannel
    participant GP as GetPagos
    B ->> BC : Inicio proceso automático
    BC ->> GPE: filtroPagosEstado()
    GPE -->> BC: pagoEstado
    BC ->> GC: filtroChannel
    GC -->> BC: Canal
    BC ->> GP: filtroPagos
    alt Error BD
        GP-->>BC: ERROR_PAGO
        BC-->>B: ERROR_PAGO
    else No encontrado
        GP-->>BC: nil
        BC ->> BC : ¿pagos pendientes < 1?
        opt pagosPendientes = 0
            GP -->>BC: ERROR_PAGO_LISTA
            BC -->>B: ERROR_PAGO_LISTA
        end
    end
```


