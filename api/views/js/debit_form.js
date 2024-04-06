Vue.component("debit-form", {
  template: `<v-form 
    @submit.prevent="postDebitForm"
    ref="debitform"
    v-model="valid"
    >
    <v-container class="py=0" fluid>
      <v-row justify="center">
        <v-col 
        cols="12"
        md="10"
        >
        <h4 align="center" class="my-2">PAGO CON TARJETA DE DEBITO</h4>
        <v-alert text type="error" v-if="disabled">
        <p>Este Método de pago no está disponible</p>
        <p>El servicio de <strong>Prisma</strong> no se encuentra disponible</p>
        <p>Por favor comuniquese con soporte o utilice un método diferente</p>
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
                      <v-img v-show="debitCardType !== ''"
                        :src="'/imgs/' + debitCardType + '.png'"
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
                    :key="key + 1"
                  >{{ h }}</span>
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
              outlined
              dense
              v-model.lazy="nombre_titular"
              label="Titular (como aparece en la Tarjeta)"              
              hint="Tal como aparece en la tarjeta"
              :disabled="disabled"
              :rules="holderNameRules"
              required
              name="nombre"
              autocomplete="off"
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
              ref="cardNumber"
              @blur="num => validarTarjeta(num.target.value)"
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
              maxlength="3"
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
          
        </v-row>

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
              autocomplete="off"
            ></v-select>
          </v-col>
          <v-col>
            <v-text-field
              type="number"
              min="0"
              outlined
              dense
              v-model="form.holder_docNum"
              label="Número de Identificación"
              :disabled="disabled"
              :rules="holderDocNumRules"
              required
              name="documento"
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
              label="E-mail"              
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


        <v-dialog v-model="alertaDebito" width="600" persistent v-if="pay.paymentMethod===2 && pay.currentStep === 1 ">
          <template>
            <v-alert
              prominent
              type="info"
              class="mb-0"
            >
              <v-row align="center" no-gutters>
                <v-col class="grow ml-3" cols="12">
                  <p class="font-weight-regular">Ud. seleccionó medio de pago <span class="font-weight-bold">Tarjeta de Débito</span></p>
                  <p class="font-weight-regular">Verifique que su tarjeta sea de DÉBITO</p>
                  <v-btn @click="volverAIntentar" color="blue lighten-5" outlined small>
                  <v-icon>
                    mdi-chevron-left
                  </v-icon>
                    Volver
                  </v-btn>
                  <v-btn @click="alertaDebito=false" small>Aceptar</v-btn>
                </v-col>
              </v-row>
            </v-alert>
          </template>
        </v-dialog>
    </v-form>
    `,
  data() {
    return {
      alertaDebito: true,

      prepaga: false, // validacion de bin: tarjeta
      credito: false, // validacion de bin: tipo de tarjeta

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
        () => this.debitCardType.length !== 0 || "Número de tarjeta no válido",
        () => this.prepaga === false || "No se permiten tarjetas prepagas",
        () => this.credito === false || "Número de tarjeta no válido",
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
      loading: false,
    };
  },

  computed: {
    disabled() {
      return !this.prisma.status;
    },
    minCardMonth() {
      if (this.form.card_year === this.minCardYear)
        return new Date().getMonth() + 1;
      return 1;
    },
    debitCardType() {
      // actualizo el numero en la tarjeta dibujada
      var v = this.form.card_number; //.replace(/[^\dA-Z]/g, '')
      var reg = new RegExp(".{4}", "g");
      this.card_display = v.replace(reg, function (a, b, c) {
        return a + " ";
      });
      // cada vez que el campo del formulario cambia
      let number = String(this.form.card_number).slice(0, 6);
      let seleccion = "";

      //filtrando tarjetas de debito
      const debitCards = this.cards.lista.filter(
        (card) => card.channel.ID === 2
      );

      // recorro la lista de tarjetas
      debitCards.forEach((card) => {
        // tomo la expresión regular
        re = new RegExp(card.regexp);
        // si coincide la expresión y es channel 2 selecciono
        if (number.match(re) != null && card.channels_id === 2) {
          seleccion = card.mediopago;
          this.cardNumberMaxLength = card.longitud_pan;

          //si la tarjeta es maestro validar que su longitud sea 18
          if (
            card.ID === 35 &&
            this.form.card_number.length !== card.longitud_pan
          ) {
            //se quita la seleccion si la longitud de la tarjeta es distinto de 18
            seleccion = "";
          }
        }
      });
      return seleccion;
    },
    cardColor() {
      let color = "primary";
      switch (this.debitCardType) {
        case "VISA":
          color = "indigo lighten-3";
          break;
        case "VISA DÉBITO":
          color = "indigo lighten-3";
          break;
        case "MASTERCARD PRISMA":
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
        case "TARJETA CLUB DíA":
          color = "red lighten-3";
          break;

        default:
          color = "blue lighten-5";
          break;
      }
      return color;
    },
    ...mapState(["pay", "form", "formErrors", "prisma", "cards", "product"]),
    ...mapGetters(["getCssMainColorPref", "getCssButtonPref"]),
  },
  watch: {
    cardYear() {
      if (this.form.card_month < this.minCardMonth) {
        this.form.card_month = "";
      }
    },
    async "form.card_number"(newCard) {
      const medioPago = this.pay.paymentMethod; //medio de pago seleccionado

      if (newCard.length === 8 && medioPago === 2) this.validarTarjeta(newCard);

      // reset validacion
      if (newCard.length < 8) {
        this.prepaga = false;
        this.credito = false;
      }
    },
  },
  methods: {
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
      const card = document.querySelectorAll(".card__inner");
      card.forEach((c) => c.classList.toggle("is-flipped"));
    },
    volverAIntentar() {
      this.nombre_titular = "";
      this.alertaDebito = true;
      this.pay.currentStep = 0;
      (this.form.channel = ""),
        (this.form.holder_name = ""),
        //this.form.holder_email = "",
        (this.form.holder_docType = "DNI"),
        (this.form.holder_docNum = ""),
        (this.form.holder_cuit = ""),
        (this.form.card_brand = ""),
        (this.form.card_number = ""),
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
      // this.$refs.debitform.reset()
      this.$emit("emit-clear");
    },
    formatCardMonth(currentCardMonth) {
      if (currentCardMonth.length == 1) {
        return "0" + currentCardMonth;
      }
      return currentCardMonth;
    },
    async postDebitForm() {
      if (this.$refs.debitform.validate()) {
        this.loading = true;
        this.form.holder_name = this.nombre_titular; // necesario porque se cambio por performance el v-model directo al form.holder_name
        this.form.card_month = this.formatCardMonth(this.form.card_month);
        this.form.card_brand = this.debitCardType;
        this.form.channel = "debit"; // Indica en canal de pago utilizado
        this.form.uuid = this.pay.uuid;
        this.form.installments = "1"; // Es el numero de cuotas
        this.form.importe = parseInt(
          parseFloat(this.product.total).toFixed(2) * 100
        );
        this.form.valorcupon = parseInt(
          parseFloat(this.product.total).toFixed(2) * 100
        );

        const status = await this.$store.dispatch("verificarCheckout");

        if (!status) window.location.reload(); // si el pago ya fue realizado recargar pagina

        if (this.valid && status) this.$store.dispatch("processPayment"); // realizar pago

        this.loading = false;
      }
    },
    soloNumeros() {
      var key = window.event.keyCode;
      if (key < 48 || key > 57) {
        window.event.returnValue = false;
      }
    },
    async validarTarjeta(nroBin) {
      const params = { NroBin: nroBin.slice(0, 8), TipoTarjeta: "D" };
      const data = await this.$store.dispatch("verificarBinTarjeta", params);

      // validamos si la tarjeta no es valida
      if (data && !data.tarjeta_valida) {
        // comprobamos si el bin ingresado pertenece a una tarjeta prepaga
        if (data.pre_pago) this.prepaga = true;

        // comprobamos si el bin ingresado pertenece a una tarjeta de credito
        if (data.type.includes("credit")) this.credito = true;

        this.$refs.cardNumber.validate();
      } else {
        this.prepaga = false;
        this.credito = false;

        this.$refs.cardNumber.validate();
      }
    },
  },
});
