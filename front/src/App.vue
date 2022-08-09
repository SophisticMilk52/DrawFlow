<template>
<FileServer/>
</template>

<script>
import FileServer from './components/FileServer.vue'
export default {
  name: 'App',
  components: {
    FileServer
  }
}
</script>

<style>
*,
:after,
:before {
  box-sizing: border-box;
  margin: 0;
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
 
}
html {
  padding: 0;
  margin: 0;
  overflow: hidden;
  background: transparent;
}
 body {
  max-width: calc(100vw - 40px);
  max-height:  calc(100vh - 40px);
  margin: 20px;
  padding: 0;
}
#app {
  margin: 0px;
  border-radius: 8px; 
  border: 1px solid #494949;
  
}
.el-main {
  padding: 0;
}
</style>
