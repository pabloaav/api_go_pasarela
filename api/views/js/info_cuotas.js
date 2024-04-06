Vue.component('info-cuotas',{
    template: `
    <div class="block">
        <p>Cuotas permitidas</p>
        <table class="table">
                      <thead>
                        <th>cuota</th>
                        <th>tem</th>
                        <th>importe</th>
                      </thead>
                      <tbody>
                        <tr v-for="plan in valorCuotas">
                          <td>{{plan.cuota}}</td>
                          <td>{{plan.tem}}%</td>
                          <td>$ {{plan.importe}}</td>
                      </tr>
                      </tbody>
        </table>
    </div>
    `,
    computed: {
        valorCuotas() {
            // tomo el valor del importe total a pagar
            let importe = this.product.total
            // transformo included_installments string con las cuotas aceptadas
            // en un array de numeros de cuotas
            var ArrayCuotas = this.pay.included_installments.split(',').map(function(value){
                return parseInt(value, 10);
            })
            // filtro los planes de cuotas de acuerdo a las cuotas permitidas para este pago
            planesFiltrados = this.cuotas.planes.filter(plan => ArrayCuotas.includes(plan.cuota))
            // voy a crear el valor a pagar por cada cuota
            cuotasValores = []
            var formatter = new Intl.NumberFormat('es-419', {
                // These options are needed to round to whole numbers if that's what you want.
                minimumFractionDigits: 2, // (this suffices for whole numbers, but will print 2500.10 as $2,500.1)
                maximumFractionDigits: 2, // (causes 2500.99 to be printed as $2,501)
              });
            planesFiltrados.forEach(element => {
                // divido el importe por la cantidad de cuotas del plan
                importeCuota = importe / element.cuota
                valor = new Object
                valor.cuota = element.cuota
                valor.tem = element.tem
                // a cada cuota la multiplico por el coeficiente
                valor.importe = formatter.format(element.coeficiente * importeCuota)
                cuotasValores.push(valor)
            })
            return cuotasValores
        },
        ...mapState(['cuotas', 'product', 'pay']),
    }
})