import Vue from 'vue'
import Router from 'vue-router'
import vueFilterPrettyBytes from 'vue-filter-pretty-bytes'
import FileList from "../../../components/FileList";


Vue.use(vueFilterPrettyBytes);
Vue.use(Router);


const routes = [
  {
    path: "/",
    redirect: "/view"
  },
  {
    path: "/view",
    name: "View",
    component: FileList
  },
];

const router = new Router({
  routes: routes,
});


window.router = router;

export default router
