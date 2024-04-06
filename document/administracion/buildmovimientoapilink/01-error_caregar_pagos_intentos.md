> # Build Movimiento ApiLink

## Error al cargar pagos intentos  
1. Valida si la lista de cierre de lote inicial tiene elementos
2. Busca los pagos intentos que pertenecen a la lista de cierre de lote
3. ERROR_PAGO_INTENTO
***


```mermaid
sequenceDiagram;
    participant B as BuildMovimientoApiLink
    participant GPI as GetPagosIntentos
    B ->> B : validar lista cierre lotes
    B ->> GPI: filtroPagoIntento
    Note over GPI: No se pudo cargar
    GPI -->> B: ERROR_PAGO_INTENTO
```


