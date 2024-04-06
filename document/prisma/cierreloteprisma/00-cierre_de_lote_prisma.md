# cierre de lote prisma - [Documento API Rest][URL-Decidir1]


El cierre de lote le permite realizar la presentación ante cada Marca de tarjeta de las operaciones de Compras, Anulaciones y Devoluciones realizadas, y de esta manera puedan ser liquidadas por cada medio de pago. Los cierres de lotes de cada medio de pago se pueden realizarse de 2 maneras:

1.Manual: esta modalidad es “on demand”. Para ello, un usuario del comercio debe ingresar a la consola de Decidir y seleccionar el medio de pago a cerrar lote. Para más detalle por favor consultar el Manual de Administración de Decidir. [Ver manual][URL-Decidir]


2.Automática: Los procesos se ejecutan diariamente luego de la medianoche, y al finalizar, se envían al comercio cada uno de los archivos del cierre de lote de cada medio de pago habilitado. Los resúmenes correspondientes a los cierres de lotes automáticos efectuados pueden ser enviados por:
    E-MAIL
    FTP/SFTP
En caso de que el comercio opte por recibir los resúmenes vía e-mail, debe indicarnos a qué dirección o direcciones de correo electrónico desea recibir tales archivos. En caso de que el comercio opte por recibir los resúmenes vía FTP o SFTP, debe indicarnos los siguientes datos: URL del servidor, usuario y clave.

Los cierre de lotes automatico tiene el sigunete formato:

Al finalizar el cierre de lote automático, se genera un archivo por medio de pago con el detalle de las transacciones incluidas en dicho cierre.

Las especificaciones del archivo son:

Nomenclatura del nombre de cada archivo generado: lote########_ddmmyy.MP.txt
########: ID Site Decidir (8 dígitos) ddmmyy:
Fecha en la que se realizó el cierre (6 dígitos).
MP: ID de medio de pago (1-3 dígitos).

Ejemplo del contenido de un archivo de cierre de lote:

D99cb9313bcdddss02400005895627249884000C210920180000018990069787200037370000000000909153570122109201800000000000000000000055506441100000000000000000000000000000000000909150000000000000000000 D99cb9313qwertyu02400005895623572629122C210920180000003699007727580037460000000000909153570022109201800000000000000000000055506441100000000000000000000000000000000000909150000000000000000000 D99cb9313cvbgnmm02400005895625363122010C210920180000006999004159200037250000000000909153570062109201800000000000000000000055506441100000000000000000000000000000000000909150000000000000000000 D99cb9313fghjklu02400005895623862118038C210920180000004799046823920037400000000000909153570032109201800000000000000000000055506441100000000000000000000000000000000000909150000000000000000000 D99vbn112bbbnnmr02400005895622572621138C210920180000002398989212490037470000000000909153570012109201800000000000000000000055506441100000000000000000000000000000000000909150000000000000000000 T000000000502435700050000036886080000000000000000000000000000000000000000000000000000000000000000000

***
## El diseño del archivo enviado es el siguiente:
### Estructura Registros Detalle:
PrismaRegistroDetalle{
	
    TipoRegistro       string  -0 TIPOREGISTRO 	      [1-1]     Tipo de Registro, Char default "D".
    IdTransacciones    string  -1 IDTRANSACCIONSITE   [2-16]    Id de Transacción, Alfanumérico de 15 dígitos, completando con "0" a la izquierda.
	IdMedioPago        int64   -2 IDMEDIOPAGO         [17-19]   Medio de Pago, numérico, 3 dígitos completando con "0" a la izquierda. Por ejemplo: 001 
                                                                identifica a Visa.
	NroTarjetaCompleto int64   -3 NROTARJETACOMPLETO  [20-39]   Nro de Tarjeta, numérico de 20 dígitos. Se informan los seis primeros digitos (BIN), últimos 4 dígitos 
                                                                del número de tarjeta y se completa con ""0"" los digitos restantes.
	TipoOperacion      string  -4 TIPOOPERACION 	  [40-40]   Operación, Char valores posibles:“C”:Compra, “D”:Devolución, “A”:Anulación.
	Fecha              string  -5 FECHA 	          [41-48]   Fecha de Operación, numérico de 8 dígitos, formato ""DDMMYYYY"".
	Monto              float64 -6 MONTO 	          [49-60]   Monto de Operación, numérico de 12 dígitos, 10 enteros + 2 decimales (sin punto decimal).
	CodAut             int64   -7 CODAUT* 	          [61-66]   Código de Autorización, numérico de 6 dígitos, completando con ""0"" a la izquierda.
	NroTicket          int64   -8 NROTICKET 	      [67-72]   Número de Cupón, numérico de 6 dígitos.
	IdSite             int64   -9 IDSITE 	          [73-87]   Id Site Decidir, numérico de 15 dígitos, el Id Site siempre es de 8 dígitos y se completa 
                                                                con 7 ceros a la izquierda.
	IdLote             int64   -10 IDLOTE 	          [88-90]   Número de lote, numérico de 3 dígitos.
	Cuotas             int64   -11 CUOTAS 	          [91-93]   Cantidad de coutas, numérico de 3 dígitos.
	FechaCierre        string  -12 FECHACIERRE        [94-101]  Fecha de cierre,numérico de 8 dígitos.
	NroEstablecimiento int64   -13 NROESTABLECIMIENTO [102-131] Número de establecimiento, numérico de 30 dígitos.
	IdCliente          string  -14 IDCLIENTE 	      [132-171] IDCLIENTE 40 caracteres completados con "0".
	Filler             string  -15 FILLER 	          [172-190] Filler, 19 caracteres completados con "0".
}
### Estructura Registro Trailer:
type PrismaRegistroTrailer struct {

	TipoRegistro      string  // 0  TIPOREGISTRO 	   [1-1] 	Tipo de Registro, Char, default ""T"".
	CantidadRegistros int64   // 1  CANTIDADREGISTROS  [2-11] 	Cantidad Registros ""Detalle"", numérico de 10 dígitos, completando con ""0"" a la izquierda.
	IdMedioPago       int64   // 2  IDMEDIOPAGO 	   [12-14] 	Medio de Pago, numérico de 3 dígitos, completando con ""0"" a la izquierda.Por ejemplo: 001 identifica a Visa
	IdLote            int64   // 3  IDLOTE 	           [15-17]  Número de Lote, numérico de 3 dígitos (000...999).
	CantidadCompras   int64   // 4  CANTCOMPRAS 	   [18-21] 	Contador de Compras, numérico de 4 dígitos (0000...9999), cantidad de compras netas.
	MontoCompras      float64 // 5  MONTOCOMPRAS 	   [22-33] 	Monto de Compras, numérico de 12 dígitos, formato $$$$$$$$$$CC, monto total de compras netas.
	CantidadDevueltas int64   // 6  CANTDEVUELTAS 	   [34-37] 	Cantidad de Devoluciones, numérico de 4 dígitos (0000...9999), cantidad de devoluciones netas.
	MontoDevueltas    float64 // 7  MONTODEVUELTAS 	   [38-49] 	Contador de Anulaciones, numérico de 12, formato $$$$$$$$$$CC, cantidad de anulaciones.
	CantidadAnuladas  int64   // 8  CANTANULADAS 	   [50-53] 	Cantidad de Anulaciones, numérico de 4 dígitos (0000...9999), monto de anulaciones.
	MontoAnulacion    float64 // 9  MONTOANULADAS 	   [54-65] 	Monto de Anulaciones, numérico de 12 dígitos, formato, monto de anulaciones.
	Filler            string  // 10 FILLER 	           [66-100] Filler, 35 caracteres completados con ""0"".
}
***

## Servicio Cierre de Lote
El servicio cierre de lote tiene como objetivo automatizar la obtención de los archivos recibido por parte de prisma interpretar su contenido y guardar los datos interpretados en una tabla de la base de datos, para que luego sea procesados por otros servicios.

- [Servicio Mover Archivo Externos][URL-SMAE]

- [Servicio Leer CierreLote][URL-SMA]


<!-- rutas -->
[URL-Decidir1]: https://decidirv2.api-docs.io/1.0/introduccion/cierre-de-lote
[URL-Decidir]: https://developers.decidir.com/sac/estaticos/manuales/Instructivo_Administracion_SPS.pdf
[URL-SMAE]:  https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/01-servicio_archivo_Lote_externo.md
[URL-SMA]: https://github.com/Corrientes-Telecomunicaciones/api_go_pasarela/blob/development/document/prisma/cierreloteprisma/01-servicio_leer_cierre_lote_prisma.md





