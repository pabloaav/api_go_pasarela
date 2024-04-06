Vue.component("debin-form", {
  template: `<v-form 
    id="debin-form"
    @submit.prevent="postDebinForm"
    ref="debinform"
    v-model="valid"
    >
    <v-container class="py=0" fluid>
    <v-row justify="center">
        <v-col 
        cols="12"
        md="10"
        >
        <h4 align="center" class="my-2">PAGO CON DEBIN</h4>
        <v-row>
          <v-col>
            <v-text-field
              class="mayus"
              outlined
              dense
              v-model="form.holder_name"
              label="Titular"
              :rules="holderNameRules"
              required
              autocomplete="off"
            />
          </v-col>
        </v-row>
        <v-row class="mt-0">
          <v-col cols="3">
            <span>CUIL / CUIT</span>
          </v-col>
          <v-col>
            <v-text-field
              type="number"
              min="0"
              outlined
              dense
              v-model="form.holder_cuit"
              label="Número de Identificación / Sin guiones"
              hint="Sin guiones o espacios"
              :rules="holderCuitCuilRules"
              required
              autocomplete="off"
            />
          </v-col>
        </v-row>
        <v-row class="mt-0">
          <v-col>
            <v-radio-group
              v-model="debin.cbuoalias"
              row
            >
              <v-radio
                label="CBU"
                value="cbu"
                >
              </v-radio>
              <v-radio
                label="Alias"
                value="alias"
                >
              </v-radio>
            </v-radio-group>
            <v-text-field
              @keypress="restrictChars($event)"
              v-show="debin.cbuoalias == 'cbu'"
              min="0"
              outlined
              dense
              v-model="form.cbu"
              label="CBU"
              :rules="cbuRules"
              autocomplete="off"
            />
            <v-text-field
              v-show="debin.cbuoalias == 'alias'"
              outlined
              dense
              v-model="form.alias"
              label="Alias"
              :rules="aliasRules"
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
          >
            </v-text-field>
          </v-col>
        <!--    <v-col>
            <v-checkbox
              outlined
              dense
              v-model="form.es_cuenta_propia"
              label="Es cuenta propia"
            ></v-checkbox>
          </v-col>
          <v-col>
            <v-select
              outlined
              dense
              v-model="form.concepto_abreviado"
              label="Motivo"
              :items="['VAR','ALQ','CUO','EXP','FAC','PRE','SEG','HON','HAB','OIN']"
            ></v-select>
          </v-col> -->
        </v-row>

        <v-row class="mt-0">
        <v-col class="text-center pt-2"">
        <v-btn v-on:click.prevent="volverAIntentar" class="rounded-xl" width="80%">
        <v-icon left>
          mdi-chevron-left
        </v-icon>
        VOLVER
        </v-btn>
      </v-col>
      <v-col class="text-center pt-2"">
        <v-btn 
          dark 
          :color="getCssButtonPref" 
          type="submit" 
          class="rounded-xl" 
          width="80%"
          :loading="loading"
        >
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
    </v-form>`,
  data() {
    return {
      valid: true,
      holderNameRules: [(v) => !!v || "el Nombre es obligatorio."],
      holderCuitCuilRules: [
        (v) => !!v || "el CUIT / CUIL es obligatorio.",
        (v) => (v && v.length < 12) || "La longitud debe ser 11 caracteres",
      ],
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
    // Reglas de validacion para el input de CBU
    cbuRules() {
      if (this.debin.cbuoalias == "cbu") {
        // Los CBU son de 22 digitos
        return [
          (v) => !!v || "el cbu es obligatorio",
          (v) => (v && v.length <= 22) || "La longitud debe ser 22 caracteres",
        ];
      } else {
        return [];
      }
    },
    // Reglas de validacion para el input de Alias
    aliasRules() {
      if (this.debin.cbuoalias == "alias") {
        return [(v) => !!v || "el alias es obligatorio"];
      } else {
        return [];
      }
    },
    ...mapState(["pay", "debin", "form", "prisma"]),
    ...mapGetters(["getCssMainColorPref", "getCssButtonPref"]),
  },

  methods: {
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
    volverAIntentar() {
      this.nombre_titular = "";
      this.pay.currentStep = 0;
      (this.form.channel = ""),
        (this.form.holder_name = ""),
        //(this.form.holder_email = ""),
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
        (this.form.es_cuenta_propia = false),
        (this.form.concepto_abreviado = ""),
        (this.form.tiempo_expiracion = 0),
        (this.form.importe = 0),
        (this.form.moneda = "ARS"),
        (this.form.recurrente = false),
        (this.form.descripcion_prestacion = "");
      // this.$refs.debinform.reset()
      this.$emit("emit-clear");
    },
    async postDebinForm() {
      if (this.$refs.debinform.validate()) {
        this.loading = true;
        this.form.channel = "debin"; // Indica en canal de pago utilizado
        this.form.card_brand = "debin";
        this.form.installments = "1"; // Es el numero de cuotas
        this.form.uuid = this.pay.uuid;

        const status = await this.$store.dispatch("verificarCheckout");

        if (!status) window.location.reload(); // si el pago ya fue realizado recargar pagina

        if (this.valid && status) this.$store.dispatch("processPayment"); // realizar pago

        this.loading = false;
      }
    },
  },
});
