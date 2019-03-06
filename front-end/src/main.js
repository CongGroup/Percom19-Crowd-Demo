import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import VueCookies from 'vue-cookies';
import axios from 'axios'
import VueAxios from 'vue-axios'
import PacMan from 'vue-spinner/src/PacmanLoader.vue';

Vue.use(VueCookies);
Vue.config.productionTip = false;
Vue.use(VueAxios,axios);
Vue.component('pacman',PacMan);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
