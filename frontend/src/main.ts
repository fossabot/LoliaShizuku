import { createApp } from "vue";
import { createPinia } from "pinia";
import "vuetify/styles";
import { createVuetify } from "vuetify";
import { aliases, fa } from "vuetify/iconsets/fa";
import { md3 } from "vuetify/blueprints";
import { lightTheme, darkTheme } from "./plugins/theme";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";

import App from "./App.vue";
import router from "./router";
import "./assets/css/style.css";
import "./assets/css/comfortaa-fonts.css";
import "@fortawesome/fontawesome-free/css/all.css";
import "unfonts.css";

const vuetify = createVuetify({
  components,
  directives,
  blueprint: md3,
  theme: {
    defaultTheme: "lightTheme",
    themes: {
      lightTheme,
      darkTheme,
    },
  },
  icons: {
    defaultSet: "fa",
    aliases,
    sets: {
      fa,
    },
  },
});

const pinia = createPinia();

createApp(App).use(pinia).use(router).use(vuetify).mount("#app");
