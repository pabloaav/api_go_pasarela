> # Build Movimiento ApiLink

## Error diferencia entre la cantidad de elementos de la lista de cierre de lote y pagos intentos  
1. Valida si la lista de cierre de lote inicial tiene elementos
2. Busca los pagos intentos que pertenecen a la lista de cierre de lote
3. Busca el pago estado Accredited porque para este se irá crear un movimiento
4. Valida si se encontro un pago intento para cada debin en la lista de cierre lote
5. Crea un log para informar la diferencia encontrada
6. ERROR_CIERRE_PAGO_INTENTO
***


```mermaid
sequenceDiagram;
    participant B as BuildMovimientoApiLink
    participant GPI as GetPagosIntentos
    participant GPE as GetPagoEstado
    participant L as CreateLog
    B ->> B : validar lista cierre lotes
    B ->> GPI: filtroPagoIntento (Pagos Intentos por debines)
    GPI -->> B: listaPagosIntentos 
    B ->> GPE: filtroPagoEstado (Pago Estado Accredited)
    GPE -->> B: pagoEstadoAcreditado
    B ->> B: Validar Cantidad Pagos Intentos
      Note over B: Cantidad Invalida
    B ->> L: log
    B -->> B: ERROR_CIERRE_PAGO_INTENTO
```


