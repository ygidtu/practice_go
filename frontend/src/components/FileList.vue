<template>
  <el-row>
    <el-row>
      <el-col :span="20">
        <el-breadcrumb separator="/" id="main">
          <el-breadcrumb-item v-for="(b, index) in bread" v-bind:key="index">
            <a v-if="b.href.length > 0" :data-href="b.href" @click="openDir(b.href)"><b>{{ b.name }}</b></a>
            <span v-else>{{ b.name }}</span>
          </el-breadcrumb-item>
        </el-breadcrumb>
      </el-col>
      <el-col :span="4">
        <el-link type="primary" @click="openDir(ajaxData.current)">刷新</el-link>
      </el-col>
    </el-row>

    <el-row>
      <el-tabs type="border-card" closable @tab-remove="removeTab"  v-model="editableTabsValue">
        <el-tab-pane label="文件" name="0">
          <el-row ref="mainContent">
            <el-scrollbar style="height: 100%;" wrapStyle="overflow-x: hidden;" viewStyle="overflow-x: hidden">
              <el-row v-for="(row, index) in ajaxData.data" :key="index" :gutter="20">
                <el-col :span="span" :key="e.path" v-for="e in row">
                  <el-card shadow="hover" class="box-card">
                    <div class="text item">
                  <span>
                    <i :class="e.is_dir ? 'el-icon-folder-opened' : 'el-icon-files'" />
                    <p v-line-clamp:15>{{ e.name }}</p>
                  </span>
                      <el-divider>
                        Size: {{ e.size | prettyBytes }}
                      </el-divider>

                      <el-row v-if="e.type === fileType.Image">
                        <el-image style="width: 100px; height: 100px" :src="urls.preview + '/' + e.path" :preview-src-list="imgUrls"></el-image>
                      </el-row>
                      <el-row :gutter="10">
                        <el-button-group>
                          <el-button v-if="e.is_dir" icon="el-icon-folder-opened" @click="openDir(e.path)" size="small">打开</el-button>
                          <el-button v-if="!e.is_dir" icon="el-icon-paperclip" v-on:click="downloadFile(e.path)" size="small">下载</el-button>
                          <el-button v-if="fileType.Not !== e.type" icon="el-icon-view" @click="previewFile(e.name, e.path, e.type)" size="small">预览</el-button>
                          <el-button icon="el-icon-s-fold" @click="compressFile(e.path)" size="small">压缩</el-button>
                          <el-button icon="el-icon-delete" @click="deleteFile(e.path)" size="small" v-if="!e.disable_delete">删除</el-button>
                        </el-button-group>
                      </el-row>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
              <el-backtop target="#main"></el-backtop>
            </el-scrollbar>
          </el-row>
        </el-tab-pane>

        <el-tab-pane
                v-for="item in editableTabs"
                :key="item.name"
                :label="item.name"
                :name="item.name"
        >
          <embed v-if="item.type === fileType.Html || item.type === fileType.Pdf" :src="item.href" style="width: 100%;" :height="clientHeight - 250 + 'px'" />
          <Media v-if="item.type === fileType.Video" kind="video" :controls="true" :src="[item.href]" style="width: 100%" :height="clientHeight - 250 + 'px'"></Media>
          <el-image v-if="item.type === fileType.Image" fit="contain" :src="item.href" :height="clientHeight - 250 + 'px'"></el-image>
        </el-tab-pane>
      </el-tabs>
    </el-row>
  </el-row>

</template>
  
<script>
  import Media from '@dongido/vue-viaudio'

  export default {
    components: {
        Media: Media,
    },
    data() {
 
      return {
        span: 6,
        title: 'Welcome to Your Vue.js App',
        ajaxData: {
          data: [],
          current: "/"
        },
        clientHeight: "600px",
        bread: [],
        drawer: false,
        urls: {
          login: new URL("/login", document.baseURI).href,
          logout: new URL("/logout", document.baseURI).href,
          list: new URL("/api/list", document.baseURI).href,
          preview: new URL("/api/preview", document.baseURI).href,
          download: new URL("/api/download", document.baseURI).href,
          compress: new URL("/api/compress", document.baseURI).href,
          delete: new URL("/api/delete", document.baseURI).href,
        },
        fileType: Object.freeze({
          "Not": 0,
          "Image":1, 
          "Html":2, 
          "Pdf":3, 
          "Video": 4
        }),
        imgUrls: [],
        editableTabsValue: '0',
        editableTabs: []
      }
    },
    methods: {

      getData: function (path) {
        const self = this;

        self.ajaxData = {
          data: [],
          current: ""
        };

        this.axios.get(self.urls.list, {
            params: {
              path: path
            }
          })
          .then(response => {
            if (response.status !== 200) {
              this.$notify({
                showClose: true,
                type: 'error',
                message: response.data["Message"]
              });
            } else {
              let temp_data = response.data;

              temp_data.data.sort(function(a, b){
                if (b.is_dir && !a.is_dir) {
                  return 1
                } else if (!b.is_dir && a.is_dir) {
                  return -1
                } else if (a.name > b.name) {
                  return 1
                } else {
                  return -1
                }
              });

              let single_row = [];
              let image_urls = [];
              temp_data.data.forEach(function (element, index) {
                if (element.type.startsWith("text/html")) {
                  element.type = self.fileType.Html
                } else if ( element.type.startsWith("image") ) {
                  element.type = self.fileType.Image;
                  image_urls.push(`${self.urls.preview}/${element.path}`)
                } else if (
                        element.type.startsWith("video") ||
                        element.type.startsWith("application/octet-stream")
                ) {
                  element.type = self.fileType.Video;
                } else if (element.type.startsWith("application/pdf")) {
                  element.type = self.fileType.Pdf;
                } else {
                  element.type = self.fileType.Not;
                }

                if (index > 0 && index % (24 / self.span) === 0) {
                  self.ajaxData.data.push(single_row);
                  single_row = []
                }
                single_row.push(element);
              });

              if (single_row.length > 0) {
                self.ajaxData.data.push(single_row);
                single_row = []
              }

              self.imgUrls = image_urls;
              this.createBread(temp_data.current)
            }
          });
      },

      createBread: function (current) {

        this.ajaxData.current = current;

        let temp_bread = [];
        temp_bread.push({
          "href": "/",
          "name": "Home"
        });

        let breads = current.split('/');
        let path = "/";
        breads.forEach(function (element, index) {
          if (element.length > 0) {
            path = `${path}/${element}`;
            if (index < breads.length - 1) {
              temp_bread.push({
                "href": `${path}`,
                "name": element
              })
            } else {
              temp_bread.push({
                "name": element,
                "href": ''
              })
            }
          }
        });

        this.bread = temp_bread
      },

      openDir: function (href) {
        this.getData(href);
      },

      deleteFile: function (href) {
        this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.axios.get(this.urls.delete, {
            params: {
              path: href
            }
          }).then((response) => {
            if (response.status !== 200) {
              this.$notify({
                showClose: true,
                type: 'error',
                message: response.data["Message"]
              });
            } else {
              this.$notify({
                showClose: true,
                type: 'success',
                message: '删除成功!'
              });

              this.openDir(this.ajaxData.current)
            }
          }).catch((error) => {
            this.$notify({
              showClose: true,
              type: 'error',
              message: error.toString()
            });
          });

        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消删除'
          });
        });
      },

      compressFile: function (href) {
        this.axios.get(this.urls.compress, {
          params: {
            path: href
          }
        }).then((response) => {
          if (response.status !== 200) {
            this.$notify({
              showClose: true,
              type: 'error',
              message: response.data["Message"]
            });
          } else {
            this.$notify({
              showClose: true,
              type: 'success',
              message: '后台压缩中!'
            });
          }
        }).catch((error) => {
          this.$notify({
            showClose: true,
            type: 'error',
            message: error.toString()
          });
        })
      },

      downloadFile: function (href) {
        window.open(`${this.urls.download}/?path=${href}`)
      },

      previewFile: function (name, href, type) {

        let res = {
          "name": name,
          "href": `${this.urls.preview}/${href}`,
          "type": type
        };

        let exists = false;

        for (let element of this.editableTabs) {
          if (element.name === name) {
            exists = true;
            break
          }
        }

        if (!exists) {
          this.editableTabs.push(res);
        }

        this.editableTabsValue = name
      },

      changeFixed(clientHeight){ //动态修改样式
        this.$refs.mainContent.$el.style.height = clientHeight - 250 + 'px';
      },

      removeTab(targetName) {
        let tabs = this.editableTabs;
        let activeName = this.editableTabsValue;
        if (activeName === targetName) {
          tabs.forEach((tab, index) => {
            if (tab.name === targetName) {
              let nextTab = tabs[index + 1] || tabs[index - 1];
              if (nextTab) {
                activeName = nextTab.name;
              }
            }
          });
        }

        this.editableTabs = tabs.filter(tab => tab.name !== targetName);

        if (this.editableTabs.length === 0) {
          this.editableTabsValue = '0'
        } else {
          this.editableTabsValue = this.editableTabs.last().name;
        }
      }
    },

    mounted() {
      // 获取浏览器可视区域高度
      this.clientHeight =   `${document.documentElement.clientHeight}`;

      window.onresize = function temp() {
        this.clientHeight = `${document.documentElement.clientHeight}`;
      };

      this.openDir("")
    },
    watch: {
      // 如果 `clientHeight` 发生改变，这个函数就会运行
      clientHeight: function () {
        console.log(this.clientHeight);
        this.changeFixed(this.clientHeight)
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1,
  h2 {
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }

  .el-dropdown+.el-dropdown {
    margin-left: 15px;
  }


</style>