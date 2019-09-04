<template>
  <div id="app">

    <el-container style="border: 1px solid #eee">
        <el-header style="text-align: left; font-size: 12px">
          <span style="font-size: 20px; color: white">Enjoy!!</span>

          <el-link @click="jumpTo('/login')" icon="el-icon-s-platform" disabled>
              登录
          </el-link>
          <el-link @click="jumpTo('/logout')" type="warning" icon="el-icon-warning">
              登出
          </el-link>
      </el-header>

      <el-container>
        <el-main @scroll.native="infiniteScroll">
            <router-view></router-view>
        </el-main>
      </el-container> 
    </el-container>

  </div>
</template>


<script>
  export default {
    name: 'App',
      data() {
        return {
            asideWidth: "200px",
            mainWidth: (window.innerWidth - this.width) + "px",
            isCollapse: false
        }
      },
      methods: {
          jumpTo: function(href) {
              window.location.replace(href)
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
    height: 100%
  }
  .el-header {
    background-color: #409EFF;
    color: #333;
    line-height: 60px;
    padding-top: 0px;
  }

  .el-aside {
    color: #333;
  }
  html, body {
    overflow: hidden;
  }
</style>
