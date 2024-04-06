# Servicio Leer Cierre Lote

## Error al obtener lista de estados desde la base  de datos 

### (se consulta la base de datos para obtener una lista de estados)
- BC: BackGround
- SCL: Servico Cierre de Lote
- SAdmin: Servicio Administración
- SCom: Servicio Commons
1. diariamente despues de la media noche BC llama al servicio cierre lote LeerCierreLote()
2. se llama al servicio SAdmin  s.adminService.GetPagosEstadosService(true, true) para obtener una lista de estados
3. verifica si devuelve error o la lista de estados, retorna error 
4. retorna constante de error: ERROR_READ_ARCHIVO

***
```mermaid
sequenceDiagram;
    participant bc as BackGround
    participant scl as Servicio Cierre Lote
    participant sa as Servicio Administración
    participant sc as Servicio Commons
    activate bc
        note over bc: diariamente despues de la media noche se llama al servicio cierre lote
        bc->>scl: LeerCierreLote()
        activate scl
            scl ->> sa: s.adminService.GetPagosEstadosService(true, true)
            activate sa
                sa-->>scl: retorna
            deactivate sa
            alt verifica si devuelve error o la lista de estados
                Note over scl: ocurrer error
                scl-->bc: retorna constante de error: ERROR_READ_ARCHIVO
            else
                Note over scl: encuentra archivos
            end
        deactivate scl
    deactivate bc
```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/01-servicio_leer_cierre_lote_prisma.md
