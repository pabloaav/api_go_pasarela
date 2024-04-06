> # Build Transferencia Cliente

## Error con saldo de cuenta
1. Busca una lista de movimientos con los ids seleccionados por el usuario.
2. Valida si la cantidad de elementos en la lista obtenida es igual a la solicitada por el usuario
3. Busca el estado Accredited para validar si todos los pagos están en estado acreditado
4. Valida si todos los pagos están en estado acreditado
5. Valida si el total solicitado corresponde al todal de los pagos
6. Buscar el saldo de la cuenta del cliente
7. Valida si el saldo de la cuenta es suficiente para realizar la transferencia
8. ERROR_SALDO_CUENTA / ERROR_SALDO_CUENTA_INSUFICIENTE
***


```mermaid
sequenceDiagram;
    participant B as BuildTransferenciaCliente
    participant GM as GetMovimientos
    participant GPE as GetPagoEstado
    participant GS as GetSaldoCuenta
    B ->> GM : filtroMovimiento (ids movimiento)
    GM -->> B: movimientos
    B->>B: Valida cantidad elementos
    B->>GPE: FiltroPagoEstado
    GPE-->>B: estadoAcreditado
    B->>B: Valida Pagos Acreditados
    B->>B: Valida Importe
    B->>GS: cuentaId
    alt Error BD
        Note over GS: No se pudo cargar
        GS-->>B: ERROR_SALDO_CUENTA
    else Saldo Insuficiente
        GS-->>B: saldoCuenta
        B->>B: Valida si saldo es suficiente 
        B-->>B: ERROR_SALDO_CUENTA_INSUFICIENTE
    end
```


