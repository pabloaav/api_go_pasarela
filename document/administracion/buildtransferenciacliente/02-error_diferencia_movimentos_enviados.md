> # Build Transferencia Cliente

## Error en la cantidad de movimentos encontrados  
1. Busca una lista de movimientos con los ids seleccionados por el usuario.
2. Valida si la cantidad de elementos en la lista obtenida es igual a la solicitada por el usuario
3. ERROR_MOVIMIENTO_LISTA_DIFERENCIA
***


```mermaid
sequenceDiagram;
    participant B as BuildTransferenciaCliente
    participant GM as GetMovimientos
    participant L as CreateLog
    B ->> GM : filtroMovimiento (ids movimiento)
    GM -->> B: movimientos
    B->>B: Valida cantidad elementos
    Note over B: Se encontro diferencia
    B->>L: log
    B-->>B: ERROR_MOVIMIENTO_LISTA_DIFERENCIA
```


