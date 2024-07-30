<template>
  <div id="app">
    <MainTemplate />
  </div>
</template>

<script>
import MainTemplate from './components/MainTemplate.vue'

export default {
  name: 'App',
  components: {
    MainTemplate
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
