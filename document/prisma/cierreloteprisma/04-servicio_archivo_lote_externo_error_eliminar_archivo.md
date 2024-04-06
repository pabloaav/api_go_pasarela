# Servicio Archivo Lote Externo

## Error al intentar remover un archivo
### (luego de mover un archivo desde una ubicación origen a una ubicacion destino con exito, el archivo de la ubicacion origen es borrado)
- SP: Servico Prisma
- SCom: Servicio Commons
- SAdmin: Servicio Administración
1. diariamente despues de la media noche BC llama al servicio SP funcion ArchivoLoteExterno() (err error)
1. obtiene de las variable de entorno la ruta donde se localizan los archivos de cierre de lote
2. se llama al servicio SCom y se le pasa el valor de la ruta s.commonsService.LeerDirectorio(config.RUTA_LOTE_FTP) 
3. verifica si devuelve error o los archivos encontrados, retorna lista_de_archivo 
4. retorna lista_de_archivo
5. define contador con total de archivos entcontrados, totalArchivos = len(archivos) 
6. recorre lista de archivos, por cada archivo 
7. llama servico SCom, s.commonsService.MoverArchivos(config.RUTA_LOTE_FTP, config.RUTA_LOTES_SIN_VERIFICAR, archivo.Name()) retorna valor nil o error
8. verifica si mover archivo produce error, retornta nil
9. la funcion retorna un valor nil al realizar el movimiento de archivo con exito
10. llama servico SCom s.commonsService.BorrarArchivo(config.RUTA_LOTE_FTP, archivo.Name()) retorna valor nil o error
11. verifica si al borrar archivo se produce error, retornta error
12. arma el mensaje de error, err = errors.New(ERROR_REMOVER_ARCHIVO + err.Error()) 
13. genera log de error, logs.Error(err.Error())
14. construye msaje error para el logs de notificaciones
15. arma la notificación, ArmarNotificacionCierreLote(msjError)
16. llama al servico SAdmin para guarda la notificación s.adminService.CreateNotificacionService(notificacion) retorna valor nil o error
17. verifica si al guardar devuelve error, retorno error
18. genera log de error, logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
19. retorna contador de archivos movidos y error

***
## posible constantes de ERRROR que se puede recibir mover un archivo
    - ERROR_REMOVER_ARCHIVO
    - ERROR_CREAR_NOTIFICACION
***
```mermaid
sequenceDiagram;
    participant bc as  BackGround
    participant sap as ServicePrisma
    participant sac as ServiceCommonds
    participant sa as ServiceAdministración
    activate bc
        note over bc: diariamente despues de la media noche se llama al servicio Archivo lote Externo
        bc->>sap: ArchivoLoteExterno() (err error)
        activate sap
            sap -->> sap: obtiene ruta variable entorno
        deactivate sap
            sap ->> sac: s.commonsService.LeerDirectorio(config.RUTA_LOTE_FTP) 
            activate sac
                sac-->>sac: Lee directorio
                alt verifica si devuelve error o los archivos encontrados
                    Note over sac: ocurrer error 
                else
                    Note over sac: encuentra archivos 
                    sac-->> sap: retorna lista_de_archivo
                end
                activate sap
                    sap-->>sap: totalArchivos = len(archivos) 
                deactivate sap
                    loop recorre lista de archivo
                        Note  over sap, sac: por cada archivo  
                        sap->>sac: s.commonsService.MoverArchivos(config.RUTA_LOTE_FTP, config.RUTA_LOTES_SIN_VERIFICAR, archivo.Name())
                        activate sac
                            sac-->>sac: abre archivo en el origen
                            alt verifica si ocurrer error 
                                note over sac: abrir archivo en destino:OK
                            end                        
                            sac-->>sac: crea archivo en el destino
                            alt verifica si ocurrer error 
                                note over sac: crear archivo en origen:ok
                            end
                            sac-->>sac: copia el archivo del origen al destino
                            alt verifica si ocurrer error 
                                note over sac: copiar del origen al destino:ok
                            end
                            sac-->>sap: retorna valor nil
                        deactivate sac
                        activate sap
                            alt verifica si mover archivo produce error
                                note over sap: retornta error                         
                            else
                                note over sap: exito al mover archivo
                                sap->>sac: s.commonsService.BorrarArchivo(config.RUTA_LOTE_FTP, archivo.Name()) retorna error
                                sac-->>sac: removerArchivo(archivo) retorna error
                                alt verifica si ocurrer error 
                                    note over sac: error al intentar Remover archivo:ok
                                end
                                    sac-->>sap: retornta error
                            end
                            activate sap
                                note over sap:arma el mensaje de error
                                sap-->>sap:err = errors.New(ERROR_REMOVER_ARCHIVO + err.Error())
                            deactivate sap
                            activate sap
                                note over sap: genera log de error
                                sap-->>sap: logs.Error(err.Error())
                            deactivate sap
                            activate sap
                                sap-->>sap:construye msaje error para el logs de notificaciones
                            deactivate sap
                            activate sap
                                note over sap:arma la notificación
                                sap-->>sap: ArmarNotificacionCierreLote(msjError)
                            deactivate sap
                            note over sap,sa:guarda la notificación
                            sap->>sa:s.adminService.CreateNotificacionService(notificacion) retorna valor nil o error
                            activate sa
                                sa-->>sap: retorna ERROR_CREAR_NOTIFICACION
                                alt verifica si al guardar devuelve error
                                    note over sap: retorno error, genera log de error
                                    activate sap
                                        sap-->>sap:logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
                                    deactivate sap
                                end
                            deactivate sa
                        deactivate sap
                    end
                sap-->>bc: retorna contador de archivos movidos y error
            deactivate sac
    deactivate bc 

```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/01-servicio_archivo_Lote_externo.md





