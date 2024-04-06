Vue.component("pick-method", {
  template: `<v-container class="text-left my-3">
                            <v-btn large elevation="2" block v-on:click="pickDebit" class="my-3 pickmethod font-weight-regular text-none text-subtitle-1" v-show="pay.included_channels.includes('DEBIT')">
                              Tarjeta de débito
                              <v-icon right>
                                mdi-chevron-right
                              </v-icon>
                            </v-btn>
                            <v-btn large elevation="2" block v-on:click="pickCard" class="my-3 pickmethod text-subtitle-1" v-show="pay.included_channels.includes('CREDIT')">
                            <v-row class="align-center">
                            <v-col>
                              <span class="font-weight-regular text-none">Tarjeta de crédito</span>
                              <v-icon right>
                                mdi-chevron-right
                              </v-icon>
                            </v-col>
                            <v-col class="text-right">
                              <span class="font-weight-regular text-none float-right">
                                <v-dialog
                                transition="dialog-top-transition"
                                max-width="600"
                                >
                                <template v-slot:activator="{ on, attrs }">
                                  <v-btn
                                    text
                                    :color="getCssMainColorPref"
                                    class="font-weight-regular text-none pa-0"
                                    v-bind="attrs"
                                    v-on="on"
                                  >ver costos</v-btn>
                                </template>
                                <template v-slot:default="dialog">
                                  <v-card>
                                    <v-toolbar v-bind:style="{ background: personalizado }" dark>
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
                                        v-for="c in cuotas"
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
                              </span>
                            </v-col>
                            </v-row>
                            </v-btn>
                            <v-btn large elevation="2" block v-on:click="pickOffline" class="my-3 pickmethod font-weight-regular text-none text-subtitle-1" v-show="pay.included_channels.includes('OFFLINE')">
                            Efectivo
                            <v-icon right>
                              mdi-chevron-right
                            </v-icon>
                            </v-btn>
                            <v-btn large elevation="2" block v-on:click="pickDebin" class="my-3 pickmethod font-weight-regular text-none text-subtitle-1" v-show="pay.included_channels.includes('DEBIN')">
                            Debin
                            <v-icon right>
                              mdi-chevron-right
                            </v-icon>
                            </v-btn>
                           
                </v-container>`,
  data() {
    return {
      panel: 0,
      // personalizado: 'linear-gradient(to right, #8c03e5,#00dbe9)'
    };
  },
  computed: {
    ...mapState(["pay", "cuotas"]),
    ...mapGetters(["getCssMainColorPref", "getCssSecondaryColorPref"]),
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
    pickCard() {
      this.pay.paymentMethod = 1;
      this.pay.currentStep = 1;
    },
    pickDebit() {
      this.pay.paymentMethod = 2;
      this.pay.currentStep = 1;
    },
    pickOffline() {
      this.pay.paymentMethod = 3;
      this.pay.currentStep = 1;
    },
    pickDebin() {
      this.pay.paymentMethod = 4;
      this.pay.currentStep = 1;
    },
  },
});
