Vue.component('cuotas-detail', {
  props: ['id', 'descripcion', 'installmentdetail'],
  template: `
    <v-card flat tile>
        <v-simple-table class="table">
            <template v-slot:default>
                      <thead>
                        <th>Número de cuotas</th>
                        <th>Porcentaje</th>
                        <th>Importe cuota</th>
                        <th>Importe Total*</th>
                      </thead>
                      <tbody style="text-align: center;">
                        <tr v-for="plan in valorCuotas">
                          <td>{{plan.cuota}}</td>
                          <td>{{plan.tem}}%</td>
                          <td>{{toPesos(plan.importe)}}</td>
                          <td>{{toPesos(plan.importeFinal)}}</td>
                      </tr>
                      </tbody>
            </template>
        </v-simple-table>
    </v-card>
    `,
  computed: {
    /**
     * Retorna cada uno de los item de la lista de cuotas que se muestra en el checkout "ver costos"
     * Cada item contiene: cuota, tem, importe, importeFinal
     * @return {Array}
     */
    valorCuotas() {
      // tomo el valor del importe total a pagar
      let importeTotal = parseFloat(this.product.total)

      // transformo included_installments string con las cuotas aceptadas
      // en un array de numeros de cuotas. Array de numeros enteros
      var ArrayCuotas = this.pay.included_installments.split(',').map(function (value) {
        return parseInt(value, 10); // El segundo argunmento de parseInt es la base numerica
      })

      // filtro los planes de cuotas de acuerdo a las cuotas permitidas para este pago
      let planesFiltrados = this.installmentdetail.filter(plan => ArrayCuotas.includes(plan.cuota))
      // el valor a pagar por cada cuota
      cuotasValores = []
      // Intl.NumberFormat Constructor para objetos que permiten el formato de números con sensibilidad al lenguaje.
      var formatter = new Intl.NumberFormat('es-419', {
        // These options are needed to round to whole numbers if that's what you want.
        minimumFractionDigits: 2, // (this suffices for whole numbers, but will print 2500.10 as $2,500.1)
        maximumFractionDigits: 2, // (causes 2500.99 to be printed as $2,501)
      });
      // El array de planesFiltrados tiene los objetos cuyos atributos corresponden a la tabla detalles de cuotas de la base de datos
      planesFiltrados.forEach(element => {
        let valor = new Object

        // let diferenciaFinanciacion = 0
        // let valorCostoFinancieroIva = 0
        // divido el importe total por la cantidad de cuotas del plan
        let importeCuota = importeTotal / element.cuota
        // se crea un objeto vacio llamado valor
        valor.cuota = element.cuota
        valor.tem = element.tem
        // a cada cuota se multiplica por el coeficiente. Da como resultado el importe total a pagar con recargo por cuota
        valor.importe = element.coeficiente * importeCuota
        valor.importeFinal = element.coeficiente * importeTotal
        // se agrega el calculo del importe de iva sobre la diferencia entre importe en una cuota y en mas de una cuota
        // diferenciaFinanciacion = valor.importeFinal - importeTotal
        // valorCostoFinancieroIva = diferenciaFinanciacion * (this.costoFinanciero / 100)
        // valor.importeFinal = valor.importeFinal + valorCostoFinancieroIva
        valor.importe = valor.importeFinal / element.cuota // el importe de cada cuota despues del recargo
        cuotasValores.push(valor)
      })

      return cuotasValores
    },
    ...mapState(['product', 'pay', 'costoFinanciero']),
  },
  methods: {
    toPesos(num) {
      return Number(num).toLocaleString('es-ar', { style: 'currency', currency: 'ARS', minimumFractionDigits: 2 })
    }
  }
})