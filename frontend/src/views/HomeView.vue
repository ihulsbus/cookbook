<template>
  <div class="suface-ground empty">
    It's very empty here...
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';

@Options({
  components: {
  },
})
export default class HomeView extends Vue {}
</script>

<style lang="scss">
  .surface-ground {
    height: 100%
  }

  .empty {
    text-align: center;
    color: gray;
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
  }
</style>
