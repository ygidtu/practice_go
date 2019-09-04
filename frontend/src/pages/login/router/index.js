import Vue from 'vue'
import Router from 'vue-router'
import vueFilterPrettyBytes from 'vue-filter-pretty-bytes'
import Login from "../Login";

Vue.use(vueFilterPrettyBytes)
Vue.use(Router)


const routes = [
    {
        path: "/",
        name: "Login",
        component: Login
    }
];

const router = new Router({
    routes: routes,
})


window.router = router

export default router
