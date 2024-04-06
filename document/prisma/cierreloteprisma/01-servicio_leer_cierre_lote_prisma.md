# Servicio Leer Cierre Lote
este servicio se encarga de leer cada uno de los archivos txt de cierre de lote enviado por prisma.
el procedimiento es el siguiente, por cada archivo se recorre el contenido del mismo y se va interpretando su contenido, de acuerdo a la estructura definida en el archivo principal.
luego de interpretado todo el contenido del archivo el mismo es guardado en la base de datos.


## Casos

- [Error al obtener lista de estados desde la base  de datos ][CASO1-ERROR]
- [Error al directorios de cierre de lotes][CASO2-ERROR]
- [Error al abrir archivo de cierre de lote][CASO3-ERROR]
- [Error al recorrer el contenido del archivo][CASO4-ERROR]
- [Exito al Recorrer los Archivo Cierre de Lote][CASOEXITO-ERROR]



<!-- rutas -->
[CASO1-ERROR]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/02-servicio_leer_cierre_lote_error_obtener_lista_estados.md

[CASO2-ERROR]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/03-servicio_leer_cierre_lote_error_al_leer_directorio.md

[CASO3-ERROR]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/04-servicio_leer_cierre_lote_error_al_abrir_archivo.md

[CASO4-ERROR]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/05-servicio_leer_cierre_lote_error_al_recorrer_contenido_archivo.md

[CASOEXITO-ERROR]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/06-servicio_leer_cierre_lote_exito.md