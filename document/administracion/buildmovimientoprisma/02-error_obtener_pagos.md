> # Build Prisma Movimiento

## Error al obtener pagos relacionados  (Validar si existe error al intentar obtener los pagos relacionados con los Cierre de lotes desde la base de datos)
1. Intenta obtener los cierre de lotes en DB por medio de la función GetPrismaCierreLotes()
2. Regresa una lista CierreLote
3. Verifica si la lista esta vacía (listaCL distinto a 0), devuelve lista con información
4. Recorre listaCL y obtiene como resultado una lista de Uuid
5. Genera tipo dato PagoFiltro
6. Intenta obtener una lista de Pagos
    - 6.1. Se produce error, genera estructura de tipo log donde se guarda el tipo de error, el mensaje y la función donde se realiza el error
    - 6.2. Inserta en la tabla log información del error, se le pasa como parametro el objeto log que contiene "tipo de error, el mensaje y la función donde se realiza el error"
7. retorna al sistema el error existente, ERROR_PAGO   
***
- nota:
    * DB = Data Base
    * listaCL = lista Cierre de Lote
    * PagoFiltro = es una estructura de dato definda que contiene 
                    (	Ids []uint64, PagoEstadosId uint64, PagoEstadosIds []uint64, CargaPagoIntentos bool, CargaMedioPagos   bool,	CargarPagoTipos bool,	CargarCuenta bool,	CargarPagoEstado  bool,	Uuids []string )
    * log = es una entidad
***
```mermaid
sequenceDiagram;
    participant SA as ServiceAdministración
    participant RA as RepositoryAdministración
    SA ->> RA: GetPrismaCierreLotes()
    RA -->>SA: Retorna una lista CierreLote
    SA -->>SA: Verifica si la lista esta vacía (lista con información)
    activate SA
    SA-->>SA: recorre listaCL:Genera una lista de UuIdd 
    deactivate SA
    activate SA
    SA-->>SA: Genera dato tipo PagoFiltro con los parametros "listaUuId, CargaPagoIntentos,	CargarPagoTipos, CargarPagoEstado" 
    deactivate SA
    SA ->> RA: GetPagos(filtroPago)
    activate RA
    Note over RA: No se pudo cargar
    
    RA->>RA: error al consultar los pagos
    activate RA
    RA-->>RA:genera estructura de tipo log (tipo de error, el mensaje y la función donde se realiza el error)     
    deactivate RA
    activate RA
     RA-->>RA:CreateLog(log)
    deactivate RA
    RA -->> SA: retorna: ERROR_PAGO 
    deactivate RA
```