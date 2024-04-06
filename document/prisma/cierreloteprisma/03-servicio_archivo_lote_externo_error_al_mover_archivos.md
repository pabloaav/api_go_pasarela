# Servicio Archivo Lote Externo

## Error al mover archivos 
### (luego de obtener una lista de archivo, estos son movidos de un directorio origen a un directorio destino)

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
8. verifica si se produce error, retornta error
9. disminuye contador en 1, totalArchivos -= 1    
10. arma el mensaje de error, err = errors.New(ERROR_MOVER_ARCHIVO + err.Error()) 
11. genera log de error, logs.Error(err.Error())
12. construye msaje error para el logs de notificaciones
13. arma la notificación, ArmarNotificacionCierreLote(msjError)
14. llama al servico SAdmin para guarda la notificación s.adminService.CreateNotificacionService(notificacion) retorna valor nil o error
15. verifica si al guardar devuelve error, retorno error
16. genera log de error, logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
17. retorna contador de archivos movidos y error

***
```mermaid
sequenceDiagram;
    participant bc as  BackGround
    participant sap as ServicePrisma
    participant sac as ServiceCommonds
    participant sa as ServiceAdministración
    activate bc
        note over bc: diariamente despues de la media noche se llama al servicio Archivo lote Externo
        bc->>sap: ArchivoLoteExterno()(err error)
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
                sac-->> sap: retorna listadearchivo
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
                            note over sac: error al abrir archivo en destino:OK
                        end                        
                        sac-->>sac: crea archivo en el destino
                        alt verifica si ocurrer error 
                            note over sac: error al crear archivo en origen:ok
                        end
                        sac-->>sac: copia el archivo del origen al destino
                        alt verifica si ocurrer error 
                            note over sac: error al copiar del origen al destino:ok
                        end
                        sac-->>sap: retorna Constante ERROR
                deactivate sac
                activate sap
                    alt verifica si mover archivo produce error
                        note over sap: retornta error 
                            activate sap
                                note over sap: disminuye contador en 1
                                sap-->>sap: totalArchivos -= 1 
                            deactivate sap
                            activate sap
                                note over sap: arma el mensaje de error
                                sap-->>sap: err = errors.New(ERROR_MOVER_ARCHIVO + err.Error())
                            deactivate sap
                            activate sap
                                note over sap: genera log de error
                                sap-->>sap: logs.Error(err.Error())
                            deactivate sap
                            activate sap
                                sap-->>sap: construye mensaje error para el logs de notificaciones
                            deactivate sap
                            activate sap
                                note over sap: arma la notificación
                                sap-->>sap: ArmarNotificacionCierreLote(msjError)
                            deactivate sap
                            note over sap,sa: guarda la notificación
                            sap->> sa: s.adminService.CreateNotificacionService(notificacion) retorna valor nil o error
                        activate sa
                            sa-->> sap: retorna ERROR_CREAR_NOTIFICACION
                            alt verifica si al guardar devuelve error
                                note over sap: retorno error, genera log de error
                                activate sap
                                    sap-->>sap:logs.Error(ERROR_AL_CREAR_NOTIFICACION + errNotificacion.Error())
                                deactivate sap
                            end
                        deactivate sa
                    else
                        note over sap: exito al mover archivo
                    end
                deactivate sap
            end
            sap-->> bc: retorna contador de archivos movidos y error
        deactivate sac
    deactivate bc
```
***
[Volver][URL-Volver]

[URL-Volver]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/01-servicio_archivo_Lote_externo.md



