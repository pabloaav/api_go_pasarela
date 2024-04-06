Vue.component("offline-form", {
  template: `<v-form 
    @submit.prevent="postOfflineForm"
    ref="offlineform"
    v-model="valid"
    >
    <v-container class="py=0" fluid>
    <v-row justify="center">
        <v-col 
        cols="12"
        md="10"
        >
        <h4 align="center" class="my-2">PAGO CON EFECTIVO</h4>
        <v-row>
          <v-col cols="4">
            <v-select
              v-model="form.holder_docType"
              :items="doctype_items"
              item-text="abbr"
              item-value="abbr"
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
              autocomplete="off"
            />
          </v-col>
        </v-row>

        <v-row class="mt-0">
          <v-col>
            <v-text-field
              class="mayus"
              outlined
              dense
              v-model="form.holder_name"
              label="Nombre Completo"
              :disabled="disabled"
              :rules="holderNameRules"
              required
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
    </v-form>
    `,
  data() {
    return {
      valid: true,
      // metodo: 'efectivo',
      doctype_items: ["DNI", "CI", "LE", "LC"],
      holderDocNumRules: [
        (v) => !!v || "El número de identificación es un campo obligatorio.",
      ],
      holderNameRules: [(v) => !!v || "El nombre del titular es obligatorio"],
      holderEmailRules: [
        (v) => !!v || "el E-mail es obligatorio.",
        (v) => /.+@.+\..+/.test(v) || "el E-mail debe ser válido.",
      ],
      loading: false,
    };
  },
  computed: {
    disabled() {
      return !this.prisma.status;
    },
    ...mapState(["pay", "form", "formErrors", "prisma"]),
    ...mapGetters(["getCssMainColorPref", "getCssButtonPref"]),
  },
  methods: {
    volverAIntentar() {
      this.nombre_titular = "";
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
      // this.$refs.offlineform.reset()
      this.$emit("emit-clear");
    },
    async postOfflineForm() {
      if (this.$refs.offlineform.validate()) {
        this.loading = true;

        this.form.channel = "offline";
        this.form.card_brand = "rapipago";
        this.form.uuid = this.pay.uuid;
        this.form.installments = "1"; // Es el numero de cuotas

        const status = await this.$store.dispatch("verificarCheckout");

        if (!status) window.location.reload(); // si el pago ya fue realizado recargar pagina

        // llama al endpoint pagar
        if (this.valid && status) this.$store.dispatch("processPayment");

        this.loading = false;
      }
    },
  },
});
