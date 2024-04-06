Vue.component("card-form", {
  template: `<v-form 
    @submit.prevent="postCardForm"
    ref="cardform"
    v-model="valid"
    >
    <v-container class="py=0" fluid>
      <v-row justify="center">
        <v-col cols="12" md="10">
          <h4 align="center" class="my-2">PAGO CON TARJETA DE CREDITO</h4>
          <v-alert text type="error" v-if="disabled">
            <p>Este Método de pago no está disponible</p>
            <p>El servicio de <strong>Prisma</strong> no se encuentra disponible</p>
            <p>Por favos comuniquese con soporte o utilice un método diferente</p>
          </v-alert>
          <v-row>
            <v-col>
            <div class="card">
            <div class="card__inner">
            <v-sheet
            class="mx-auto pa-5 rounded-lg card__face card__face--front" 
            :color="this.cardColor"
            elevation="5" 
            max-width="400"
            min-height="120">
                <v-row>
                  <v-col cols="3">
                  
                    <v-img src="/imgs/chip.png" class="pa-3 mx-auto" height="50" contain></v-img>
                
                  </v-col>
                  <v-col cols="9">
                    <div class="card-type float-right">
                      <v-scroll-x-transition>
                        <v-img v-show="cardType !== ''"
                          :src="'/imgs/' + cardType + '.png'"
                          height="50"
                          width="65"
                          contain
                        ></v-img>
                      </v-scroll-x-transition>
                    </div>
                  </v-col>
                </v-row>
                <div class="card-number text-center">
                  <div v-if="card_display !== ''">
                    <v-scroll-y-transition class="py-0" group tag="span">
                      <span class="text-h5 primary--text" v-for="(n, key) in card_display"
                        :key="key + 1"
                      >{{ n }}</span>
                    </v-scroll-y-transition>
                  </div>
                  <span class="text-h5 primary--text" v-else> #### #### #### #### </span>
                </div>
                <div class="card-name text-center">
                  <div v-if="nombre_titular !== ''">
                  <v-scroll-x-transition class="py-0" group tag="span">
                    <span class="text-h5 primary--text text-uppercase" v-for="(h, key) in nombre_titular"
                    :key="key + 1">{{ h }}</span>
                  </v-scroll-x-transition>
                  </div>
                  <span class="text-h5 primary--text" v-else> NOMBRE DEL TITULAR </span>
                </div>
                <v-row no-gutters class="my-5">
                  <v-col float-left>
                  <v-row no-gutters>
                    <div v-if="form.card_month !== ''">
                      <v-scroll-y-transition>
                        <span class="text-subtitle-2 primary--text">{{formatCardMonth(form.card_month)}}</span>
                      </v-scroll-y-transition>
                    </div>
                    <span v-else class="text-subtitle-2 primary--text">MM</span>
                    <span class="text-subtitle-2 primary--text">&nbsp; / &nbsp;</span>
                    <div v-if="form.card_year !== ''">
                      <v-scroll-y-transition>
                        <span class="text-subtitle-2 primary--text">{{String(form.card_year).slice(2, 4)}}</span>
                      </v-scroll-y-transition>
                    </div>
                    <span v-else class="text-subtitle-2 primary--text">AA</span>
                  </v-row>
                  </v-col>
                </v-row>
              </v-sheet>
              <v-sheet
                class="mx-auto rounded-lg card__face card__face--back" 
                color="blue lighten-5"
                elevation="5" 
                width="400"
                min-height="216">
                <v-col px-0>
                <div class="band__space"></div>
                <div class="grey darken-4 magnetic__band" px-0></div>
                <div class="band__space"></div>
                <v-row justify="center">
                <v-card color="white" min-height="40"></v-card>
                </v-row>
                <v-row>
                <div class="card__code text-right">
                  <v-scroll-y-transition class="py-0" group tag="span">
                    <span class="text-h5 primary--text text-right" v-for="(h, key) in form.card_code"
                      :key="key + 1"
                    > * </span>
                  </v-scroll-y-transition>
                </div>
                </v-row>
                </v-col>
              </v-sheet>
              </div>
              </div>
            </v-col>
          </v-row>

          <v-row class="mt-0">
            <v-col>
              <v-text-field
              class="mayus"
                v-model.lazy="nombre_titular"
                outlined
                dense
                label="Titular (como aparece en la Tarjeta)"              
                hint="Tal como aparece en la tarjeta"
                :disabled="disabled"
                :rules="holderNameRules"
                required
                name="nombre"
              />
            </v-col>
          </v-row>

          <v-row class="mt-0">
            <v-col>
              <v-text-field
                type="number"
                min="0"
                outlined
                dense
                v-model="form.card_number"
                label="Número de la Tarjeta"  
                autocomplete="off"
                name="CardNumber"            
                :disabled="disabled"
                :rules="holderCardNumRules"
                required
                counter
                :maxlength="cardNumberMaxLength"
                @keypress="restrictChars($event)"
                @input="recortarCVV"
                @blur="num => validarTarjeta(num.target.value)"
                ref="cardNumber"
              ></v-text-field>
              <v-row>
              <v-col>
              </v-col>
              </v-row>
            </v-col>
            <v-col cols="4">
              <v-text-field
                outlined
                dense
                v-model="form.card_code"
                label="Codigo de seguridad"
                :rules="holderCardCodeRules"
                autocomplete="off"
                :type="showCode ? 'text' : 'password'"
                required
                :append-icon="showCode ? 'mdi-eye' : 'mdi-eye-off'"
                @focus="flip()"
                @blur="flip()"
                @click:append="showCode = !showCode"
                @keypress="restrictChars($event)"
                :maxlength="longitud_codigo"
                ref="myInput"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row class="mt-0">
            <v-col>
              <v-select
                outlined
                dense
                v-model="form.card_month"
                label="Mes de Vencimiento"
                :items="meses"
                :rules="holderCardMonthRules"
                maxlength="2"
                required
              ></v-select>
            </v-col>
            <v-col>
              <v-select
                outlined
                dense
                v-model="form.card_year"
                label="Año de Vencimiento"
                :items="anios"
                :rules="holderCardYearRules"
                maxlength="4"
                required
              ></v-select>
            </v-col>          
            <v-col>
              <v-select
                outlined
                dense
                v-model="form.installments"
                label="Cuotas"
                :items="itemsAndTotal"
                @change="(valor)=>calcularMontoCuota(valor,true)"
                required
              ></v-select>
            </v-col>
          
          </v-row>

          <v-dialog
          transition="dialog-top-transition"
          max-width="600"
          >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
            small
              outlined
              :color="getCssMainColorPref"
              v-bind="attrs"
              v-on="on"
              @click="panel=0"
            >VER COSTOS</v-btn>
          </template>
          <template v-slot:default="dialog">
            <v-card>
                <v-toolbar
                  v-bind:style="{background: personalizado}"
                  dark
                >
                <v-toolbar-title>COSTOS POR PAGO EN CUOTAS</v-toolbar-title>
                <v-spacer></v-spacer>
                  <v-btn
                    icon
                    dark
                    @click="dialog.value = false"
                  >
                    <v-icon>mdi-close</v-icon>
                  </v-btn>
                </v-toolbar>
              <v-card-text class="pt-2">
              <v-expansion-panels accordion v-model="panel">
                <v-expansion-panel
                  v-for="c in cuotas.filter(cuota => cuota.mediopagoinstallments_id === mediopagoinstallmentsId)"
                  :key=c.id
                >
                  <v-expansion-panel-header>
                  <!-- <img v-if="c.descripcion == 'TELCO'" src="/imgs/login_logo.png"  style="width: 120px; flex:none !important "> -->
                  <p v-if="c.descripcion == 'TELCO'"> GENERAL </p>
                  <p v-else> {{c.descripcion}} </p>
                  </v-expansion-panel-header>
                  <v-expansion-panel-content>
                    <cuotas-detail v-bind="c"></cuotas-detail>
                    <p v-if="c.descripcion == 'TELCO'">* Los importes incluyen costo financiero de IVA</p>
                  </v-expansion-panel-content>
                </v-expansion-panel>
              </v-expansion-panels>
              </v-card-text>
              <v-card-actions class="justify-center">
                <v-btn
                  text
                  @click="dialog.value = false"
                >Cerrar</v-btn>
              </v-card-actions>
            </v-card>
          </template>
          </v-dialog>

          <v-row class="mt-0">
            <v-col cols="4">
              <v-select
                v-model="form.holder_docType"
                :items="doctype_items"
                label="Tipo"
                return-object
                single-line
                outlined
                dense
                required
              ></v-select>
            </v-col>
            <v-col>
              <v-text-field
                type="number"
                outlined
                dense
                v-model="form.holder_docNum"
                label="Número de Identificación"
                :disabled="disabled"
                :rules="holderDocNumRules"
                required
                name="documento"
                min="0"
                autocomplete="off"
              />
            </v-col>
          </v-row>

          <v-row class="mt-0">
            <v-col>
              <v-text-field
                outlined
                dense
                v-model="form.holder_email"
                label="Correo electrónico"              
                :disabled="disabled"
                :rules="holderEmailRules"
                required
                name="email"
                validate-on-blur
                autocomplete="off"
              ></v-text-field>
            </v-col>
          </v-row>

          <v-row class="mt-0">
            <v-col class="text-center pt-2">
              <v-btn v-on:click.prevent="volverAIntentar" class="rounded-xl" width="80%">
              <v-icon left>
                mdi-chevron-left
              </v-icon>
              VOLVER
              </v-btn>
            </v-col>
            <v-col class="text-center pt-2">
              <v-btn dark :color="getCssButtonPref" type="submit" class="rounded-xl" width="80%" :loading="loading">
              Pagar
                <v-icon right>
                  mdi-chevron-right
                </v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-container>


    <v-dialog v-model="alertaCredito" width="600" persistent v-if="pay.paymentMethod===1 && pay.currentStep === 1">
      <template>
        <v-alert
          prominent
          type="info"
          class="mb-0"
        >
          <v-row align="center" no-gutters>
            <v-col class="grow ml-3" cols="12">
              <p class="font-weight-regular">Ud. seleccionó medio de pago <span class="font-weight-bold text-subtitle-1">Tarjeta de Crédito</span></p>
              <p class="font-weight-regular">Confirmar pagar con TARJETA DE CRÉDITO</p>
              <v-btn @click="volverAIntentar" color="blue lighten-5" outlined small>
              <v-icon>
                mdi-chevron-left
              </v-icon>
                Volver
              </v-btn>
              <v-btn @click="alertaCredito=false" small>Aceptar</v-btn>
            </v-col>
          </v-row>
        </v-alert>
      </template>
    </v-dialog>

    </v-form>
    `,
  data() {
    return {
      loading: false,

      prepaga: false, // validacion de bin: tarjeta
      debito: false, // validacion de bin: tipo de tarjeta

      alertaCredito: true,
      nombre_titular: "",
      meses: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"],
      anios: [
        "2023",
        "2024",
        "2025",
        "2026",
        "2027",
        "2028",
        "2029",
        "2030",
        "2031",
        "2032",
        "2033",
        "2034",
        "2035",
      ],
      mediopagoinstallmentsId: 1, // valor que corresponde a Telco por default
      myInput: "",
      longitud_cvv: 3,
      valid: true,
      cardNumberMaxLength: 19,
      card_display: "",
      showCode: false,
      minCardYear: new Date().getFullYear(),
      doctype_items: ["DNI", "CUIL", "LE", "LC"],
      holderEmailRules: [
        (v) => !!v || "el E-mail es obligatorio.",
        (v) => /.+@.+\..+/.test(v) || "el E-mail debe ser válido.",
      ],
      holderNameRules: [(v) => !!v || "El nombre del titular es obligatorio"],
      holderDocNumRules: [
        (v) => !!v || "El número de identificación es un campo obligatorio.",
      ],
      holderCardNumRules: [
        (v) => !!v || "El número de tarjeta es un campo obligatorio.",
        () => this.cardType.length !== 0 || "Número de tarjeta no válido",
        () => this.prepaga === false || "No se permiten tarjetas prepagas",
        () => this.debito === false || "Número de tarjeta no válido",
        (v) =>
          v.length === this.cardNumberMaxLength ||
          `El número de tarjeta debe tener ${this.cardNumberMaxLength} dígitos`,
      ],
      holderCardCodeRules: [
        (v) => !!v || "El código de seguridad es un campo obligatorio.",
      ],
      holderCardMonthRules: [
        (v) => !!v || "La fecha de expiración es un campo obligatorio.",
      ],
      holderCardYearRules: [
        (v) => !!v || "La fecha de expiración es un campo obligatorio.",
      ],
      panel: 0,
      seleccionado: "",
    };
  },
  watch: {
    async "form.card_number"(newCard) {
      const medioPago = this.pay.paymentMethod; //medio de pago seleccionado

      if (newCard.length === 8 && medioPago === 1) this.validarTarjeta(newCard);

      if (newCard.length < 8) {
        this.prepaga = false;
        this.debito = false;
      }
    },
  },
  computed: {
    ...mapState(["pay", "cuotas"]),
    longitud_codigo() {
      return this.longitud_cvv;
    },
    disabled() {
      return !this.prisma.status;
    },
    minCardMonth() {
      if (this.form.card_year === this.minCardYear)
        return new Date().getMonth() + 1;
      return 1;
    },
    // cardType() se ejecuta cada vez que se ingresa un numero en el input de numero de tarjeta
    // retorna un string vacio si no matchea con la tarjeta, o un string con el nombre de la tarjeta seleccionada
    cardType() {
      // actualizo el numero en la tarjeta dibujada
      var v = this.form.card_number; // el valor ingresado en el campo de numero de tarjeta
      var reg = new RegExp(".{4}", "g");
      this.card_display = v.replace(reg, function (a, b, c) {
        return a + " ";
      });

      // se toman los 6 primeros digitos del valor ingresado
      let number = String(this.form.card_number).slice(0, 6);
      // blanqueo el tipo de tarjeta
      let nombreMedioPago = "";

      // recorro la lista de tarjetas
      // por cada tarjeta de credito (credit o channel 1) se evalua si el numero ingresado se corresponde (match) con la expresion regular
      // de cada tarjeta obtenida de la tabla correspondiente de la BD

      //filtrando tarjetas de debito
      const creditCards = this.cards.lista.filter(
        (card) => card.channel.ID === 1
      );

      creditCards.forEach((card) => {
        // tomo la expresión regular
        re = new RegExp(card.regexp);
        // si coincide la expresión y es channel 1 selecciono
        // channel 1 es credit. Correspnde a tarjetas de credito
        if (
          number.match(re) != null &&
          card.channels_id === 1 &&
          this.form.card_number.length <= card.longitud_pan
        ) {
          this.longitud_cvv = card.longitud_cvv;
          nombreMedioPago = card.mediopago; // El nombre de la tarjeta, en este caso (medio de pago en general)
          this.cardNumberMaxLength = card.longitud_pan; // Primary Account Number (PAN) es el numero de la tarjeta. puede haber diferentes longitudes

          // Si el mediopagoinstallments no es Telco se debe recalcular los montos de las cuotas
          if (nombreMedioPago) {
            // recalcular con el nuevo Installments
            // filtrar el array de cuotas (Installments)
            var distintoInstallment = this.cuotas.filter(
              (cuota) =>
                cuota.mediopagoinstallments_id == card.mediopagoinstallments_id
            );
            // La funcion filter devuelve un array aunque contenga cero o un objetos
            this.mediopagoinstallmentsId =
              distintoInstallment[0].mediopagoinstallments_id;
            // Recalcular los valores de cuotas
            this.itemsAndTotal;
          }
        }
      });

      return nombreMedioPago;
    },
    cardColor() {
      let color = "primary";
      switch (this.cardType) {
        case "VISA":
          color = "indigo lighten-3";
          break;
        case "VISA DÉBITO":
          color = "indigo lighten-3";
          break;
        case "MASTERCARD Prisma":
          color = "light-blue lighten-3";
          break;
        case "MASTERCARD DEBIT PRISMA":
          color = "light-blue lighten-3";
          break;
        case "DINERS CLUB":
          color = "blue-grey lighten-3";
          break;
        case "TARJETA NARANJA":
          color = "orange lighten-3";
          break;
        case "AMERICAN EXPRESS":
          color = "teal lighten-3";
          break;
        case "CABAL PRISMA":
          color = "deep-purple lighten-3";
          break;
        case "TARJETA NEVADA":
          color = "red lighten-3";
          break;
        case "NATIVA":
          color = "deep-purple lighten-3";
          break;
        case "TARJETA CLUB DÍA":
          color = "red lighten-3";
          break;

        default:
          color = "blue lighten-5";
          break;
      }
      return color;
    },
    /**
     * Retorna un array con los valores de las cuotas mas el importe a pagar
     * @returns {Array}
     */
    itemsAndTotal() {
      // Las cuotas (installments) que puede elegir el usuario. Cada numero de cuota es un string
      var arrayCuotas = this.pay.included_installments.split(",");

      var itemsTotal = arrayCuotas.map((item) => {
        return (
          item + " - " + this.toPesos(this.calcularMontoCuota(item, false))
        );
      });
      return itemsTotal;
    },
    ...mapState([
      "pay",
      "form",
      "product",
      "formErrors",
      "prisma",
      "cards",
      "cuotas",
      "costoFinanciero",
      "importePagar",
    ]),
    ...mapGetters([
      "getCssMainColorPref",
      "getCssSecondaryColorPref",
      "getCssButtonPref",
    ]),
    personalizado() {
      return (
        "linear-gradient(to right," +
        this.getCssMainColorPref +
        "," +
        this.getCssSecondaryColorPref +
        ")"
      );
    },
  },
  methods: {
    // Para el caso de que se pone un numero de American Express el codigo cambia de 4 a 3 digitos
    recortarCVV() {
      this.form.card_code = this.$refs.myInput.value.slice(
        0,
        this.longitud_cvv
      );
    },
    toPesos(num) {
      return Number(num).toLocaleString("es-ar", {
        style: "currency",
        currency: "ARS",
        minimumFractionDigits: 2,
      });
    },
    // Volver mayusculas las letras del input del holder name
    uppercase() {
      this.nombre_titular = this.nombre_titular.toUpperCase();
    },
    // Validacion en input de CBU para restringir la entrada solo a numeros
    restrictChars: function ($event) {
      if (
        $event.charCode === 0 ||
        /\d/.test(String.fromCharCode($event.charCode))
      ) {
        return true;
      } else {
        $event.preventDefault();
      }
    },
    flip() {
      const card = document.querySelector(".card__inner");
      card.classList.toggle("is-flipped");
    },
    volverAIntentar() {
      this.nombre_titular = "";
      this.alertaCredito = true;
      this.pay.currentStep = 0;
      (this.form.channel = ""),
        (this.form.holder_name = ""),
        //this.form.holder_email = "",
        (this.form.holder_docType = "DNI"),
        (this.form.holder_docNum = ""),
        (this.form.holder_cuit = ""),
        (this.form.card_number = ""),
        (this.form.card_brand = ""),
        (this.form.card_expiration = ""),
        (this.form.card_month = ""),
        (this.form.card_year = ""),
        (this.form.card_code = ""),
        (this.form.cbu = ""),
        (this.form.alias = ""),
        (this.form.installments = ""),
        (this.form.uuid = ""),
        (this.form.id = ""),
        (this.form.es_cuenta_propia = true),
        (this.form.concepto_abreviado = "VAR"),
        (this.form.tiempo_expiracion = 0),
        (this.form.importe = 0),
        (this.form.moneda = "ARS"),
        (this.form.recurrente = false),
        (this.form.descripcion_prestacion = "");
      this.card_display;
      // this.$refs.cardform.reset()
      this.$emit("emit-clear");
      // cuando se presiona el boton de cancelar se debe setear el importe inicial como importe a pagar
      store.commit("setImportePagar", this.product.total);
    },
    formatCardMonth(currentCardMonth) {
      if (currentCardMonth.length == 1) {
        return "0" + currentCardMonth;
      }
      return currentCardMonth;
    },

    // Este metodo es la accion que se ejecuta en el boton submit del form de card_form. Enviar los datos del formulario al backend
    // La funcion del Store que se despacha es processPayment.
    // Los datos del formulario se toman en dicha funcion desde el state del store con state.form
    async postCardForm() {
      if (this.$refs.cardform.validate()) {
        this.form.holder_name = this.nombre_titular; // necesario porque se cambio por performance el v-model directo al form.holder_name
        this.form.card_month = this.formatCardMonth(this.form.card_month);
        this.form.card_brand = this.cardType;
        this.form.channel = "credit"; // Indica en canal de pago utilizado
        this.form.uuid = this.pay.uuid;

        // Enviar el importe de Pago bruto que se recibe del backend tal cual viene, pero de tipo entero sin decimales
        this.form.importe = parseInt(parseFloat(this.product.total) * 100);

        // Enviar el importe a pagar final que es la suma del importe bruto, mas cuotas mas costo financiero. Se toma del Store
        this.form.valorcupon = parseInt(
          parseFloat(this.importePagar).toFixed(2) * 100
        );

        // Como existe un v-model en el input de cuotas: v-model="form.installments" los installments quedan con un string compuesto.
        // se debe extraer el numero de la cuota
        this.form.installments = this.form.installments.split("-")[0].trim();

        this.loading = true;

        const status = await this.$store.dispatch("verificarCheckout");

        if (!status) window.location.reload();

        if (this.valid && status) {
          // llama al endpoint pagar
          this.$store.dispatch("processPayment");
          this.loading = false;
        }
      }
    },
    soloNumeros() {
      var key = window.event.keyCode;
      if (key < 48 || key > 57) {
        window.event.returnValue = false;
      }
    },
    // Este metodo calcula el monto a pagar segun  el usuario selecciona un numero de cuotas de tarjeta de credito
    calcularMontoCuota(cuota, setStoreImporte = false) {
      var cuotasValores = [];

      // el importe total a pagar
      let importeTotal = parseFloat(this.product.total);

      var cuotaSeleccionada = "";
      if (setStoreImporte) {
        // el usuario selecciona una opcion del select. EL value viene concatenado con cuota + importe total. se debe separar la cuota del string
        var valorSelect = cuota.split("-");
        cuotaSeleccionada = valorSelect[0].trim();
      } else {
        cuotaSeleccionada = cuota;
      }

      // transformo included_installments string con las cuotas aceptadas
      // en un array de numeros de cuotas. Array de numeros enteros
      var ArrayCuotas = this.pay.included_installments
        .split(",")
        .map(function (value) {
          return parseInt(value, 10); // El segundo argunmento de parseInt es la base numerica
        });

      // Determinar el installment_id segun el mediopagoinstallmentsId. Filter devuelve array
      var inst = this.cuotas.filter(
        (c) => c.mediopagoinstallments_id == this.mediopagoinstallmentsId
      );

      // filtro los planes de cuotas de acuerdo a las cuotas permitidas para este pago.
      // en la funcion de filtro, cada objeto detallecuota es un objeto installmentdetail
      let planesFiltrados = inst[0].installmentdetail.filter((detallecuota) =>
        ArrayCuotas.includes(detallecuota.cuota)
      );

      // El array de planesFiltrados tiene los objetos cuyos atributos corresponden a la tabla detalles de cuotas de la base de datos
      planesFiltrados.forEach((element) => {
        let valor = new Object();
        if (element.cuota == cuotaSeleccionada) {
          // divido el importe total por la cantidad de cuotas del plan
          let importeCuota = importeTotal / element.cuota;

          valor.cuota = element.cuota;

          valor.tem = element.tem;

          valor.importeFinal = element.coeficiente * importeTotal;

          valor.importe = valor.importeFinal / element.cuota;

          cuotasValores.push(valor);
        }
      });

      if (setStoreImporte) {
        store.commit("setImportePagar", cuotasValores[0].importeFinal);
      } else {
        return cuotasValores[0].importeFinal;
      }
    },
    async validarTarjeta(NroBin) {
      const params = { NroBin: NroBin.substr(0, 8), TipoTarjeta: "C" };
      const data = await this.$store.dispatch("verificarBinTarjeta", params);

      // validamos si la tarjeta no es valida
      if (data && !data.tarjeta_valida) {
        // comprobamos si el bin ingresado pertenece a una tarjeta prepaga
        if (data.pre_pago) this.prepaga = true;

        // comprobamos si el bin ingresado pertenece a una tarjeta de debito
        if (data.type.includes("debit")) this.debito = true;

        this.$refs.cardNumber.validate(); // mostrar mensaje de error
      } else {
        this.prepaga = false;
        this.debito = false;
        this.$refs.cardNumber.validate();
      } // limpiar mensaje de error
    },
  },
  mounted() {
    this.$store.dispatch("getPlandeCuotas");
  },
});
