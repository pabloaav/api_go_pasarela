package util

import (
	"fmt"

	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/filtros/administracion"
)

type UtilRepository interface {
	CreateNotificacion(notificacion entities.Notificacione) error
	CreateLog(log entities.Log) (erro error)
	GetConfiguracion(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, erro error)
	GetConfiguracionesRepository(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, erro error)
	CreateConfiguracion(config entities.Configuracione) (id uint, erro error)
	CrearPeticionesRepository(peticionWeb entities.Webservicespeticione) (erro error)
	GetImpuestoByIdRepository(id int64) (impuesto entities.Impuesto, erro error)
	CreateHistorialOperacionesRepository(historial entities.HistorialOperaciones) (erro error)
}

func NewUtilRepository(conn *database.MySQLClient) UtilRepository {
	return &utilRepository{
		SqlClient: conn,
	}
}

type utilRepository struct {
	SqlClient *database.MySQLClient
}

func (r *utilRepository) CreateNotificacion(notificacion entities.Notificacione) error {

	resp := r.SqlClient.Omit("id").Create(&notificacion)

	if resp.Error != nil {

		erro := fmt.Errorf(ERROR_CREAR_NOTIFICACION)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "CreateNotificacion",
		}

		err := r.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
		return erro

	}
	return nil
}

func (r *utilRepository) CreateLog(log entities.Log) (erro error) {

	resp := r.SqlClient.Omit("id").Create(&log)

	if resp.Error != nil {
		return fmt.Errorf("error al crear log %s", resp.Error.Error())
	}

	return nil
}

func (r *utilRepository) GetConfiguracion(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, erro error) {

	resp := r.SqlClient.Model(entities.Configuracione{})

	if len(filtro.Nombre) > 0 {

		resp.Where("nombre", filtro.Nombre)
	}

	resp.Find(&configuracion)

	if resp.Error != nil {

		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguracion",
		}

		err := r.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *utilRepository) GetConfiguracionesRepository(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, erro error) {
	resp := r.SqlClient.Model(entities.Configuracione{})
	if filtro.Buscar {
		resp.Where("nombre like ?", "%"+filtro.Nombrelike+"%")
	}
	resp.Find(&configuraciones)
	if resp.Error != nil {
		erro = fmt.Errorf(ERROR_CONFIGURACIONES)

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguracionesRepository",
		}
		err := r.CreateLog(log)
		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}
	return
}

func (r *utilRepository) CreateConfiguracion(config entities.Configuracione) (id uint, erro error) {

	result := r.SqlClient.Omit("id").Create(&config)

	if result.Error != nil {
		erro = fmt.Errorf(ERROR_CREAR_CONFIGURACIONES)
		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       result.Error.Error(),
			Funcionalidad: "CreateConfiguracion",
		}

		err := r.CreateLog(log)

		if err != nil {
			mensaje := fmt.Sprintf("%s, %s", err.Error(), result.Error.Error())
			logs.Error(mensaje)
		}

		return
	}

	id = config.ID

	return
}

func (r *utilRepository) CrearPeticionesRepository(peticionWeb entities.Webservicespeticione) (erro error) {
	resp := r.SqlClient.Create(&peticionWeb)
	if resp.Error != nil {
		return fmt.Errorf("error al registrar peticiones %s", resp.Error.Error())
	}
	return nil
}

func (r *utilRepository) GetImpuestoByIdRepository(id int64) (impuesto entities.Impuesto, erro error) {
	resp := r.SqlClient.Where("id = ?", id).Find(&impuesto)
	if resp.Error != nil {
		return entities.Impuesto{}, fmt.Errorf("error al obtener impuesto %s", resp.Error.Error())
	}
	return
}

/*
	permite leer el contenido de un archivo y retorna
	elcontenido del archivo en byte, el nombre del archivo y tipo del archivo, o error
*/
func LeerDatosArchivo(rutadestino string, ruta string, nombreArchivo string) (data []byte, archivonombre string, archivotipo string, erro error) {
	data, erro = ioutil.ReadFile(fmt.Sprintf("%s/%s", ruta, nombreArchivo))

	if erro != nil {
		msj := "error a leer datos del archivo:" + nombreArchivo
		logs.Error(msj)
		erro = errors.New(msj)
		return
	}
	milissconds := time.Now().Unix()
	// fecha := fmt.Sprintf("%v-%v-%v_%v:%v:%v_", time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	fecha := fmt.Sprintf("%v_", milissconds)
	archivo_extension := strings.Split(nombreArchivo, ".")
	archivonombre = fmt.Sprintf("%s/%v%s.%s", rutadestino, fecha, archivo_extension[0], archivo_extension[1])
	archivotipo = archivo_extension[len(archivo_extension)-1]
	return
}

func LeerArchivo(rutadestino string, ruta string, nombreArchivo string) (data []byte, archivonombre string, archivotipo string, erro error) {
	data, erro = ioutil.ReadFile(fmt.Sprintf("%s/%s", ruta, nombreArchivo))

	if erro != nil {
		msj := "error a leer datos del archivo:" + nombreArchivo
		logs.Error(msj)
		erro = errors.New(msj)
		return
	}

	archivo_extension := strings.Split(nombreArchivo, ".")
	archivonombre = fmt.Sprintf("%s/%v", rutadestino, archivo_extension[0])
	archivotipo = archivo_extension[len(archivo_extension)-1]
	return
}

func (r *utilRepository) CreateHistorialOperacionesRepository(historial entities.HistorialOperaciones) error {
	resp := r.SqlClient.Create(&historial)

	if resp.Error != nil {
		return fmt.Errorf("error al registrar el historial %s", resp.Error.Error())
	}

	return nil
}
