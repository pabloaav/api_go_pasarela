package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	main "github.com/Corrientes-Telecomunicaciones/api_go_pasarela/api"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/pkg/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newDatabaseTest() *database.MySQLClient {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.USER_TEST, config.PASSW_TEST, config.HOST_TEST, config.PORT_TEST, config.DB_NAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Error("cannot create mysql client")
		panic(err)
	}
	// cada test va a crear una conexion asi que es importante que la libere cuando deje de usarla
	// mysqldb, err := db.DB()
	// if err != nil {
	// 	logs.Error("cannot get database specific driver")
	// 	panic(err)
	// }
	// defer mysqldb.Close()
	// nos aseguramos que las tablas esten creadas de acuerdo a las estructuras definidas en las entities
	// db.AutoMigrate(
	// 	entities.Adquiriente{},
	// 	entities.Apilinkcierrelote{},
	// 	entities.Channel{},
	// 	entities.Cliente{},
	// 	entities.Cuenta{},
	// 	entities.Cuentacomision{},
	// 	entities.Impuesto{},
	// 	entities.Installment{},
	// 	entities.Installmentdetail{},
	// 	entities.Log{},
	// 	entities.Mediopago{},
	// 	entities.Movimiento{},
	// 	entities.Movimientocomisiones{},
	// 	entities.Notificacione{},
	// 	entities.Pago{},
	// 	entities.Pagoestado{},
	// 	entities.Pagoestadoexterno{},
	// 	entities.Pagoestadologs{},
	// 	entities.Pagointento{},
	// 	entities.Pagoitems{},
	// 	entities.Pagotipo{},
	// 	entities.Prismacierrelote{},
	// 	entities.Transferencia{},
	// )

	return &database.MySQLClient{DB: db}
}

var app *fiber.App

func Inicializar() *fiber.App {
	httpClient := http.DefaultClient
	sqlClient := newDatabaseTest()
	osFile := os.File{}

	app = main.InicializarApp(httpClient, sqlClient, &osFile)
	return app
}

// NOTE test para generar nueva solicituds de pago
// NOTE agregar opcion para enviar iterar y realizar varias pruebas sobre distintas request
func TestNewCheckout(t *testing.T) {
	// httpClient := http.DefaultClient
	// sqlClient := newDatabaseTest()
	// osFile := os.File{}

	app = Inicializar()

	// limpiamos la tabla
	// insertamos datos necesarios para la prueba
	var pagoitems []entities.Pagoitems
	request := dtos.PagoRequest{
		PayerName:         "Jose",
		Description:       "Test de integracion",
		FirstTotal:        1000,
		FirstDueDate:      "10-07-2021",
		ExternalReference: "15685",
		SecondDueDate:     "10-08-2021",
		SecondTotal:       1010,
		PayerEmail:        "jose.alarcon@telco.com.ar",
		PaymentType:       "sellos",
		Items: append(pagoitems, entities.Pagoitems{
			Quantity:    1,
			Description: "pago prueba",
			Amount:      1000,
		}),
	}

	requestJson, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/checkout", bytes.NewBuffer(requestJson))
	req.Header.Add("apikey", "123123123123123")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	resp, err := app.Test(req, -1) //para simular la solicitud en tu aplicación y obtén la respuesta
	if err != nil {
		logs.Error("error al ejecutar el pago: " + err.Error())
	}
	defer resp.Body.Close()
	// if resp.StatusCode != 201 {
	// 	bytresp, _ := io.ReadAll(resp.Body)
	// 	fmt.Println(string(bytresp))
	// }
	// var response *dtos.CheckoutResponse
	// json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
