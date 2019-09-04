// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

import 'element-ui/lib/theme-chalk/index.css'

import {
  Input,
  InputNumber,
  Button,

  Form,
  FormItem,
    
  Alert,
  
  Icon,
  Row,
  Col,
   
  Container,
  Header,
  Main,

  Backtop,
  
  Notification
} from 'element-ui';

Vue.use(Input);
Vue.use(InputNumber);
Vue.use(Form);
Vue.use(FormItem);
Vue.use(Button);
Vue.use(Backtop);

Vue.use(Alert);

Vue.use(Icon);
Vue.use(Row);
Vue.use(Col);

Vue.use(Container);
Vue.use(Header);
Vue.use(Main);

Vue.prototype.$notify = Notification;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: {
    App
  },
  template: '<App/>'
});