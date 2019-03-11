import Vue from "vue";
import App from "./App.vue";
import VueCookies from 'vue-cookies';
Vue.use(VueCookies);
import router from "./router";
import store from "./store";
import axios from 'axios';
import VueAxios from 'vue-axios';
import VeeValidate from 'vee-validate';
import PacMan from 'vue-spinner/src/PacmanLoader.vue';

const dictionary = {
    en: {
        messages: {
            max: (_,ref)=> `input length should be less than ${ref[0]}`,
            required: (_)=> 'Input field must not be empty',
        }
    },
};


// console.log(VeeValidate.Validator.dictionary);
VeeValidate.Validator.localize(dictionary);
// console.log(VeeValidate.Validator.dictionary);
console.log(typeof process.env.SERVER_PATH);

Vue.use(VeeValidate);
Vue.config.productionTip = false;
Vue.use(VueAxios,axios);
Vue.component('pacman',PacMan);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
