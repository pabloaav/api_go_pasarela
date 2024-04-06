Vue.component("result-card", {
  template: `<v-container class="py=0" fluid>
    <v-row justify="center">
      <v-col 
      cols="12"
      md="10"
      >

      <v-card class="pa-5 rounded-xl" elevation="3">
      <div v-show="!this.resultado.isLoading">
      <v-container v-show="this.resultado.status">

        <div v-if="this.form.channel == 'debin' || this.form.channel == 'offline'">
          <v-alert text color='amber' icon="mdi-check-circle">
              <h3>Su pago está siendo procesado</h3>
          </v-alert>
          <p class="font-weight-bold amber--text">
            {{this.form.channel == 'debin'
              ? "Para finalizar la operación acceda a su home banking y autorice su debin."  
              : "Para finalizar la operación diríjase a una sucursal de Rapipago."
            }}
          </p>
            <div class="notification is-primary is-light">
              <p>Estado: En proceso</p>
              <p>Concepto: {{ this.resultado.description }}</p>
              <p>Monto: {{ toPesos(this.resultado.importe_pagado) }}</p>
            </div>
        </div>

        <div v-else>
          <v-alert text type="success">
              <h3>Su pago ha sido procesado exitosamente</h3>
          </v-alert>
            <div class="notification is-primary is-light">
              <p>Estado: Aprobado</p>
              <p>Concepto: {{ this.resultado.description }}</p>
              <p>Monto: {{ toPesos(this.resultado.importe_pagado) }}</p>
            </div>
        </div>

          <v-row class="field is-grouped">
            <v-col class="text-center">
              <v-btn dark :color="getCssButtonPref" v-on:click="abrirPdfLink" width="100%">
                    <v-icon dense left>
                      fas fa-print
                    </v-icon>
                    {{this.form.channel == 'offline' 
                      ? "Imprimir cupón" 
                      : "Imprimir Recibo"
                    }}
              </v-btn>
            </v-col>
            <v-col class="text-center">
              <v-btn dark :color="getCssButtonPref" v-on:click="redirigirExitoso" width="100%">
                    <v-icon dense left>
                      fas fa-undo
                    </v-icon>
                    Aceptar
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        <v-container v-show="!this.resultado.status">
        <v-alert text type="error" >
          <h3 class="subtitle has-text-danger">{{ this.resultado.message }}</h3>
        </v-alert>
          <div class="notification is-light">
            <p>Vuelva a intentarlo haciendo click <a v-on:click="volverAIntentar">aquí!</a></p>
            <p>O pongase en contacto con soporte <strong>TelCo</strong>: reclamoswee@telco.com.ar</p>
          </div>
          <v-row class="field is-grouped">
            <v-col class="text-center">
              <v-btn dark color="teal" v-on:click="volverAIntentar" width="100%">
                <v-icon dense left>
                  fas fa-undo
                </v-icon>
                Reintentar
              </v-btn>
            </v-col>  
            <v-col class="text-center">
              <v-btn v-on:click="redirigirRechazado" width="100%">
                Cancelar
                <v-icon dense right>
                  fas fa-times
                </v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        </div>
        <div class="text-center" v-show="this.resultado.isLoading">
          <v-progress-circular
          class="ma-2"
          :rotate="360"
          :size="100"
          :width="15"
          :value="value"
          color="teal"
          >
          {{ value }}
          </v-progress-circular>
        </div>  
      </v-card>

      </v-col>
    </v-row>
    </v-container>
    `,
  data() {
    return {
      interval: {},
      value: 0,
    };
  },
  beforeDestroy() {
    clearInterval(this.interval);
  },
  mounted() {
    this.interval = setInterval(() => {
      if (this.value === 100) {
        return (this.value = 0);
      }
      this.value += 10;
    }, 1000);
  },
  computed: {
    ...mapState(["resultado", "form", "pay", "product"]),
    ...mapGetters(["getCssMainColorPref", "getCssButtonPref"]),
  },
  methods: {
    volverAIntentar() {
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
      store.commit("setImportePagar", this.product.total); // En caso de reinternar el pago, se resetea el importe total inicial
      this.$emit("emit-clear");
    },
    redirigirRechazado() {
      window.location = this.pay.back_url_rejected;
    },
    redirigirExitoso() {
      // controlar que venga algo en la metadata
      let location = this.resultado.metadata
        ? this.pay.back_url_success + "/" + this.resultado.metadata
        : this.pay.back_url_success;
      window.location = location;
    },
    abrirPdfLink() {
      window.open(this.resultado.pdf_url, "_blank");
    },
    toPesos(num) {
      return Number(num).toLocaleString("es-ar", {
        style: "currency",
        currency: "ARS",
        minimumFractionDigits: 2,
      });
    },
  },
});
