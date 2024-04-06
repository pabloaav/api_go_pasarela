> # Check Prisma

## Error al consultar api de prisma
1. El proceso se inicia al llamar al servicio CheckPrisma
2. Se envía el request a la api de prisma
3. Devuelvo el Error de llamado a la api

***


```mermaid
sequenceDiagram;
    participant S as Servicio
    participant P as Prisma
    S ->> P: request
    P ->> P: CheckService
    P -->> S: Error en llamado a la api
    S ->> S: return Error 
    
```    