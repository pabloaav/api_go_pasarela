> # Get Billing

## caso exito get billing
1. El proceso se inicia al llamar al servicio GetBilling
2. Se consulta al repositorio con el uuid
3. Consulta el pago intento asociado al pago correcto
4. Agrega el formato con la información al generador de PDF
5. El generador de pdf devuelve el archivo
6. devuelve el archivo al frontend

***

```mermaid
sequenceDiagram;
    participant S as Servicio
    participant R as Repository
    participant M as Generador PDF
    S ->> R: uuid
    R ->> R: GetPagoByUuid
    R -->> S: Pago
    S ->> R: pagoId
    R ->> R: GetValidPagointento
    R -->> S: Pagointento
    S ->> M: formato e info
    M ->> M: GenerarPDF
    M ->> S: PDF
    S ->> S: return PDF
        
```   