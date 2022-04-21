<template>
<div class="container">
  <div class="header">
    <HeaderBar/>
  </div>
  <div class="content">
    <router-view />
  </div>
</div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import HeaderBar from '@/components/HeaderBar.vue';

@Options({
  components: {
    HeaderBar,
  },
})
export default class cookbook extends Vue {}
</script>

<style lang="scss" scoped>
  #app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
  }

html, body {
  margin: 0;
  height: 100%;
}

.container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  flex-shrink: 0;
  height: 55px;
  z-index: 10;
}

.content {
  flex-grow: 1;
  overflow-y: scroll;
}
</style>
