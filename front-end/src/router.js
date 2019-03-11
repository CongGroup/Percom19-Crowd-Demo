import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

const router = new Router({
  routes: [
      {
          path: "/admin",
          name: "admin",
          // route level code-splitting
          // this generates a separate chunk (about.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () =>
              import("./views/admin.vue"),
          meta: {
              title: 'crowd-demo',
              requireAuth: true
          }
      },

    {
        path: "",
        name: "home",
        component: ()=> import("./views/Home.vue"),
        meta: {
            title: 'crowd-demo',
            requireAuth: false
        }
    },
      {
          path: '/login',
          name: "login",
          component: ()=> import("./views/login.vue"),
          meta: {
              title: 'crowd-demo',
              requireAuth: false,
          }
      }
  ]
});

router.beforeEach((to, from, next) => {
    document.title="crowd demo";
    next();
    // if (to.name !== 'admin') {
    //     next('/')
    //     if (store.state.ks != undefined) {
    //         next('/');
    //     } else {
    //         next();
    //     }
    // } else {
    //     if (to.meta.requireAuth && store.state.ks == undefined) {
    //         next('/login')
    //
    //     } else {
    //         next();
    //     }
    // }
})


export default router;
