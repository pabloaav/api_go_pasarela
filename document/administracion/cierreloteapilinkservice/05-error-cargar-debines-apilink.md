> # Cierre de Lote ApiLink 

## Error al cargar los debines en apilink
1. El proceso se inicia automaticamente en el horario definido
2. Busca el pago estado processing (estado inicial para los debines) filtroPagosEstado()
3. Busca el canal debin filtroChannel
4. Busco los pagos que pertenecen al pago estado punto 1
5. Filtro los pagos que pertenecen al canl del punto 2
6. Busco en apilink todos los pagos que están en estado pendientes y que sean debin
7. ERROR_GET_DEBINES

***


```mermaid
sequenceDiagram;
    participant B as BuildMovimiento
    participant BC as BuildCierreLoteApiLinkService
    participant GPE as GetPagoEstado
    participant GC as GetChannel
    participant GP as GetPagos
    participant GD as GetDebinesApiLinkService
    B ->> BC : Inicio proceso automático
    BC ->> GPE: filtroPagosEstado
    GPE ->> GPE: Busca pago estado
    GPE -->> BC: pagoEstado
    BC ->> GC: filtroChannel
    GC ->> GC: Busca Canal
    GC -->> BC: Canal
    BC ->> GP: filtroPagos
    GP ->> GP: Busca Pagos
    GP -->>BC: pagos
    BC ->> BC: Filtra Pago por Canal
    BC ->> GD: uuid, request
    GD ->> GD:  Busca Debines Pendientes
    alt Error BD
        GD-->>BC: ERROR_GET_DEBINES
        BC-->>B: ERROR_GET_DEBINES
    else No encontrado
        GD-->>BC: nil
        BC ->> BC : ¿existe debines?
        opt debines = 0
            BC -->>B: ERROR_DEBINES
        end
    end
    
```


