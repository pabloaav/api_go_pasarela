> # Cierre de Lote ApiLink 

## Error al cargar los pagos estados externos
1. El proceso se inicia automaticamente en el horario definido
2. Busca el pago estado processing (estado inicial para los debines) filtroPagosEstado()
3. Busca el canal debin filtroChannel
4. Busco los pagos que pertenecen al pago estado punto 1 
5. Filtro los pagos que pertenecen al canl del punto 2
6. Busco en apilink todos los pagos que están en estado pendientes y que sean debin
7. Busco los pagos estados externos para poder comparar debines con pagos
8. ERROR_PAGO_ESTADO_EXTERNO

***


```mermaid
sequenceDiagram;
    participant B as BuildMovimiento
    participant BC as BuildCierreLoteApiLinkService
    participant GPE as GetPagoEstado
    participant GC as GetChannel
    participant GP as GetPagos
    participant GD as GetDebinesApiLinkService
    participant GPEX as GetPagosEstadosExternos
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
    GD -->> BC: debines
    BC ->> GPEX: filtroPagoEstadoExternos
    GPEX ->> GPEX: Busca Pagos Estados Externos
    alt Error BD
        GPEX -->> BC: ERROR_PAGO_ESTADO_EXTERNO
        BC -->> B: ERROR_PAGO_ESTADO_EXTERNO
    else No encontrado
        GPEX-->>BC: nil
        BC ->> BC : ¿pagos estados externos < 1?
        opt Pagos estados externos = 0
            GPEX ->> L: crea un log de error
            GPEX -->>BC: ERROR_PAGO_ESTADO_EXTERNO_LISTA
            BC -->>B: ERROR_PAGO_ESTADO_EXTERNO_LISTA
        end
    end
    
    
```


