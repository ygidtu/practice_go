<template>
  <div id="app">

    <el-container style="border: 1px solid #eee">
      <el-header style="text-align: right; font-size: 12px">
        <span style="font-size: 20px">Enjoy!!</span>
      </el-header>
      <el-container :width="main_width_">

        <el-main>
          <el-row :gutter="24">
            <el-col :sm="8" :md="10" :lg="12" :offset="6">
              <div class="grid-content bg-purple">
                <el-form :model="loginForm" status-icon :rules="rules" ref="loginForm" label-width="100px"
                  method="post" action="/login">
                  <el-form-item label="用户名" prop="username">
                    <el-input id="username" name="username" placeholder="请输入用户名" type="text" autocomplete="off"
                      v-model="loginForm.username"></el-input>
                  </el-form-item>
                  <el-form-item label="密码" prop="passwd">
                    <el-input id="password" name="password" placeholder="请输入密码" type="password" autocomplete="off"
                      v-model="loginForm.passwd"></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" native-type="submit">提交</el-button>
                    <el-button v-on:click="resetForm()">重置</el-button>
                  </el-form-item>
                </el-form>
              </div>
            </el-col>
          </el-row>
        </el-main>

      </el-container>
    </el-container>

    <el-backtop target="#app"></el-backtop>
  </div>
</template>


<script>
  export default {
    name: 'Login',

    data() {
      let validateUsername = (rule, value, callback) => {

        if (value === '') {
          callback(new Error('请输入用户名'))
        } else {
          callback()
        }
      };

      let validatePass = (rule, value, callback) => {

        if (value === '') {
          callback(new Error('请输入密码'));
        } else if (value < 5) {
          callback(new Error('密码长度过短'));
        } else {
          callback();
        }
      };


      return {
        defaultActive: "login",
        loginForm: {
          username: 'admin',
          passwd: '',
        },
        rules: {
          passwd: [{
            required: true,
            validator: validatePass,
            trigger: 'blur'
          }],
          username: [{
            required: true,
            validator: validateUsername,
            trigger: 'blur',
          }]
        },
      }
    },
    methods: {
      resetForm() {
        this.loginForm.username = "admin";
        this.loginForm.passwd = "";
      }
    },
    mounted() {
      if (document.getElementById("error")) {
        this.$notify.error({
          title: '错误',
          message: '用户名/密码错误'
        });
      }
    }
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    /*margin-top: 60px;*/
  }

  .el-header {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
  }

  .el-aside {
    color: #333;
  }
</style>