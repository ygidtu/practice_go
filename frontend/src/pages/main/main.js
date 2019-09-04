// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import 'element-ui/lib/theme-chalk/index.css'
import lineClamp from 'vue-line-clamp'

import {
  Menu,
  MenuItem,
  Input,
  InputNumber,
  Button,
  ButtonGroup,
  Table,
  TableColumn,
  Breadcrumb,
  BreadcrumbItem,
  Form,
  FormItem,

  Tabs,
  TabPane,
  Tag,
  
  Alert,
  
  Icon,
  Row,
  Col,
  
  Card,
 
  Container,
  Header,
  Aside,
  Main,

  Link,
  Divider,
  Image,

  MessageBox,
  Message,
  Notification,
  Scrollbar
} from 'element-ui';


Vue.config.productionTip = false;
Vue.use(VueAxios, axios);

Vue.use(lineClamp);
Vue.use(Menu);
Vue.use(MenuItem);
Vue.use(Input);
Vue.use(InputNumber);
Vue.use(Button);
Vue.use(ButtonGroup);
Vue.use(Tabs);
Vue.use(TabPane);
Vue.use(Breadcrumb);
Vue.use(BreadcrumbItem);
Vue.use(Tag);
Vue.use(Alert);
Vue.use(Icon);
Vue.use(Row);
Vue.use(Col);
Vue.use(Card);
Vue.use(Container);
Vue.use(Header);
Vue.use(Aside);
Vue.use(Main);
Vue.use(Link);
Vue.use(Divider);
Vue.use(Image);
Vue.use(Scrollbar);


Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$alert = MessageBox.alert;
Vue.prototype.$confirm = MessageBox.confirm;
Vue.prototype.$prompt = MessageBox.prompt;
Vue.prototype.$notify = Notification;
Vue.prototype.$message = Message;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: {
    App
  },
  template: '<App/>'
});
