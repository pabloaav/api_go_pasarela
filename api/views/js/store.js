var mapState = Vuex.mapState;
var mapGetters = Vuex.mapGetters;

const store = new Vuex.Store({
  state: model,
  mutations: {
    seleccionaTarjeta(state, card) {
      state.cards.seleccionada = card;
      state.form.card_brand = card.mediopago;
    },
    setCargando(state, value) {
      state.cargando.isLoading = value;
    },
    setImportePagar(state, value) {
      state.importePagar = value;
    },
    setCostoFinanciero(state, value) {
      state.costoFinanciero = value;
    },
  },
  actions: {
    // Funcion Principal del Pago.
    async processPayment({ state }) {
      this.state.resultado.isLoading = true;
      this.state.pay.currentStep = 2;
      this.state.form.holder_email = this.state.form.holder_email.trim();

      const settings = {
        method: "POST",
        body: JSON.stringify(state.form),
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };
      fetch("/checkout/pagar", settings)
        .then(async (response) => {
          const data = await response.json();

          if (!response.ok) {
            const error = (data && data.message) || response.status;
            this.state.resultado.message = error;
            return Promise.reject(error);
          }

          // La response es OK.
          // El error puede venir dentro de la response OK desde las servicios de pagos externos.
          // Se pregunta por la data en los atributos que contienen el exito o no del proceso
          const estado = data.data.estado; // Un mensaje de error o una respuesta de exito
          const exito = data.status; // Un valor booleano que indica si el proceso fue exitoso o no

          if (!exito) {
            // Si en el mensaje de error aparece una referencia a el error de bin (4 ultimos digitos de la tarjeta):
            const error402 = estado.includes("402");

            if (!error402) {
              const binError = estado.includes("bin"); //validamos si el error es por el bin de la tarjeta
              const invalidCardError = estado.includes("invalid_card");

              if (binError) {
                const error =
                  "Error en el número de tarjeta. Por favor, revise los valores ingresados.";
                this.state.resultado.message = error;
                return Promise.reject(error);
              } else if (invalidCardError) {
                const error =
                  "Tarjeta no válida. Por favor, revise el valor ingresado.";
                this.state.resultado.message = error;
                return Promise.reject(error);
              } else {
                const error = estado.split("-")[0];
                this.state.resultado.message = error;
                return Promise.reject(error);
              }
            } else {
              //se ejecuta cuando existe un error con codigo de estado 402

              const error = estado.split("Descripcion: ")[1];
              this.state.resultado.message = error;
              return Promise.reject(error);
            }
          }

          this.state.resultado = data.data;

          this.state.resultado.message = data.message;
          this.state.resultado.status = data.status;

          this.state.resultado.isLoading = false;
        })
        .catch((error) => {
          this.state.resultado.status = false;
          if (this.state.resultado.message == "") {
            this.state.resultado.message = "Un error inesperado ha ocurrido";
          }
          this.state.resultado.isLoading = false;

          console.error("There was an error!", error);
        });
    },

    // mount index.html()
    async checkPrisma({ commit, state }) {
      const settings = {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };

      fetch("/checkout/prisma", settings)
        .then(async (response) => {
          const data = await response.json();

          if (!response.ok) {
            const error = (data && data.message) || response.status;
            this.state.prisma.message = error;
            return Promise.reject(error);
          }
          this.state.prisma.message = data.message;
          this.state.prisma.status = data.status;
        })
        .catch((error) => {
          this.state.prisma.status = false;
          if (this.state.prisma.message == "") {
            this.state.prisma.message =
              "El servicio de prisma no se encuentra disponible.";
          }
          console.error("There was an error!", error);
        });
    },
    // mount index.html
    async getPlandeCuotas({ commit, state }) {
      const settings = {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };

      await fetch("/administracion/plan-cuotas", settings)
        .then(async (response) => {
          const data = await response.json();

          if (!response.ok) {
            const error = (data && data.message) || response.status;
            return Promise.reject(error);
          }
          this.state.cuotas = data.reverse();
        })
        .catch((error) => {
          console.error("Error all consultar planes de cuotas!", error);
        });
    },
    // mount index.html
    async getTarjetas({ commit, state }) {
      const settings = {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };

      fetch("/checkout/tarjetas/all", settings)
        .then(async (response) => {
          const data = await response.json();

          if (!response.ok) {
            const error = (data && data.message) || response.status;
            return Promise.reject(error);
          }
          state.cards.lista = data.data;
        })
        .catch((error) => {
          console.error("There was an error!", error);
        });
    },
    // Obtiene el costo financiero relacionado con el IVA sobre el importe de recargo en las cuotas
    // el resultado se guarda en propiedad costoFinanciero del Store
    async getCostoFinanciero({ commit, state }) {
      const settings = {
        method: "GET",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      };

      fetch("/administracion/obtener-impuesto", settings)
        .then(async (response) => {
          const data = await response.json();
          if (!response.ok) {
            const error = (data && data.message) || response.status;
            this.state.resultado.message = error;
            return Promise.reject(error);
          }
          commit("setCostoFinanciero", data.porcentaje * 100);
          this.state.resultado.isLoading = false;
        })
        .catch((error) => {
          console.error("There was an error!", error);
        });
    },

    // Se encarga de comprobar si el pago que se intenta realizar ya fue pagado en otra instancia
    async verificarCheckout() {
      const uuid = this.state.form.uuid;

      const response = await fetch(`/checkout/verificar/pago/estado/${uuid}`);
      const { status } = await response.json();

      return status;
    },

    // Verifica si los primeros 8 digitos de la tarjeta ingresada es valida
    async verificarBinTarjeta({}, params) {
      const { NroBin, TipoTarjeta } = params;

      const ruta = "https://api-bin.wee.ar/bin/v1/validar-bin?";

      const response = await fetch(
        ruta + new URLSearchParams({ NroBin, TipoTarjeta })
      );
      const { data } = await response.json();

      return data;
    },
  },
  getters: {
    getCssAmountPref(state, getters) {
      return {
        "--background":
          "linear-gradient(to right," +
          state.preferences.mainColor +
          "," +
          getters.getCssSecondaryColorPref +
          ")",
        "--border-radius": "10px",
        "--font-weight": "normal",
      };
    },
    getCssBackgroundPref(state, getters) {
      if (state.preferences.client == "wee") {
        return {
          "--background-image":
            "linear-gradient(to right top, #d16ba5, #c777b9, #ba83ca, #aa8fd8, #9a9ae1, #8aa7ec, #79b3f4, #69bff8, #52cffe, #41dfff, #46eefa, #5ffbf1)",
        };
      }
      return {
        "--background-image":
          "linear-gradient(to right top," +
          state.preferences.mainColor +
          "," +
          getters.getCssSecondaryColorPref +
          ")",
      };
    },
    getCssMainColorPref(state) {
      return state.preferences.mainColor;
    },
    getCssSecondaryColorPref(state) {
      var color =
        state.preferences.secondaryColor == ""
          ? "#ffffff"
          : state.preferences.secondaryColor;
      return color;
    },
    getCssTablePref(state) {
      return {
        "--border-top": "2px solid " + state.preferences.mainColor,
        "--background": state.preferences.secondaryColor,
      };
    },
    getCssTabSliderPref(state, getters) {
      return {
        "--background":
          "linear-gradient(to right," +
          state.preferences.mainColor +
          "," +
          getters.getCssSecondaryColorPref +
          ")",
      };
    },
    getCssButtonPref(state) {
      if (state.preferences.client == "wee") {
        return "teal";
      }
      return state.preferences.mainColor;
    },
    getLogoClientPref(state) {
      if (state.preferences.client == "wee") {
        return "/imgs/wee_por_telco.png";
      }
      return "data:image/png;base64," + state.preferences.logo;
      // return "/imgs/logos/" + state.preferences.logo;
    },
    getIfClientCheckoutWee(state) {
      return state.preferences.client == "wee" ? true : false;
    },
  },
});
