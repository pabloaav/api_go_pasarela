> # Cierre de Lote ApiLink 

## Error al buscar pagos por canal
1. Busca el pago estado processing (estado inicial para los debines) filtroPagosEstado()
2. Busca el canal debin filtroChannel
3. Busco los pagos que pertenecen al pago estado punto 1
4. Filtro los pagos que pertenecen al canl del punto 2
5. ERROR_PAGO_PENDIENTE

***


```mermaid
sequenceDiagram;
    participant B as BuildMovimiento
    participant BC as BuildCierreLoteApiLinkService
    participant GPE as GetPagoEstado
    participant GC as GetChannel
    participant GP as GetPagos
    B ->> BC : Inicio proceso automático
    BC ->> GPE: filtroPagosEstado
    GPE -->> BC: pagoEstado
    BC ->> GC: filtroChannel
    GC -->> BC: Canal
    BC ->> GP: filtroPagos
    GP -->>BC: pagos
    BC ->> BC : ¿Pago Channel Processing?
    opt No existen pagos (nil)
        BC -->>B: ERROR_PAGO_PENDIENTE
    end
    
```


