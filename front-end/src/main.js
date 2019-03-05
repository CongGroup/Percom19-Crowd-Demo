import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import VueCookies from 'vue-cookies';
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.use(VueCookies);
Vue.config.productionTip = false;
Vue.use(VueAxios,axios);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
