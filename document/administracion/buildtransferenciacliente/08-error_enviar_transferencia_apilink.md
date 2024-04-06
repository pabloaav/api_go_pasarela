> # Build Transferencia Cliente

## Error al enviar las transferencia a apilink
1. Busca una lista de movimientos con los ids seleccionados por el usuario.
2. Valida si la cantidad de elementos en la lista obtenida es igual a la solicitada por el usuario
3. Busca el estado Accredited para validar si todos los pagos están en estado acreditado
4. Valida si todos los pagos están en estado acreditado
5. Valida si el total solicitado corresponde al todal de los pagos
6. Buscar el saldo de la cuenta del cliente
7. Valida si el saldo de la cuenta es suficiente para realizar la transferencia
8. Crea los movimientos de salida
9. Envía la transferencia para apilink
10. En caso de error baja los movimientos creados anteriormente
11. En caso de error al bajar los movimientos crea una notificacion
12. Error enviar transferencias
***


```mermaid
sequenceDiagram;
    participant B as BuildTransferenciaCliente
    participant GM as GetMovimientos
    participant GPE as GetPagoEstado
    participant GS as GetSaldoCuenta
    participant MT as C_M_Transferencia
    participant CT as C_T_ApiLinkService
    participant BM as BajaMovimiento
    participant CN as C_Notificacion_S
    B ->> GM : filtroMovimiento (ids movimiento)
    GM -->> B: movimientos
    B->>B: Valida cantidad elementos
    B->>GPE: FiltroPagoEstado
    GPE-->>B: estadoAcreditado
    B->>B: Valida Pagos Acreditados
    B->>B: Valida Importe
    B->>GS: cuentaId
    GS-->>B: saldoCuenta
    B->>B: Valida si saldo es suficiente 
    B->>MT: crea BD listaMovimientos
    B->>CT: envía Apilink requerimientoId, request.Transferencia
    CT-->>B: Error Cargar Transferencia
    B->> BM: Baja BD listaMovimientos, err.Error()
    Note over BM: No se pudo Bajar
    BM-->>B: ERROR_BAJAR_MOVIMIENTOS
    B->> CN: crea BD notificacion
    B-->>B: Error enviar transferencias
```


