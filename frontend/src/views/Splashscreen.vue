<template>
  <div class="background">
    <div class="surface-card p-4 shadow-2 border-round w-full lg:w-2 splash">
      <div class="text-center mb-5">
        <img
          src="~@/assets/images/logos/cb_large.png"
          alt="Logo Image"
          height="130"
          class="mb-3 logo">
        <div class="text-900 text-3xl font-medium mb-3">Welcome to the cookbook</div>
      </div>
      <div>
        <Button
          label="Sign In"
          icon="pi pi-user"
          class="w-full"
          @click="$router.push('App')"
        >
        </Button>
      </div>
    </div>
    <div class="attribute">
      Background by <a href="https://unsplash.com/@toddquackenbush?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Todd Quackenbush</a> on <a href="https://unsplash.com/@ih16/likes?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Unsplash</a>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue } from 'vue-class-component';

export default class Splashscreen extends Vue {}
</script>

<style lang="scss" scoped>
  .splash {
    position: fixed;
    top: 45%;
    left: 50%;
    transform: translate(-50%, -50%);
  }

  .loginbutton {
    position: fixed;
    left: 50%;
    transform: translate(-50%, -50%);
  }

  .attribute {
    text-align: center;
    color: gray;
    position: fixed;
    top: 95%;
    left: 50%;
    transform: translate(-50%, -50%);
  }

  a:link {
    color: gray;
    text-decoration: none;
  }

  a:visited {
    color: gray;
    text-decoration: none;
  }

  a:hover {
    text-decoration: underline;
  }

  a:active {
    text-decoration: underline;
  }

  .background {
    width: 100vw;
    height: 100vh;
    background: url('~@/assets/images/splash/one.jpg') no-repeat center center fixed;
    -webkit-background-size: cover;
    -moz-background-size: cover;
    -o-background-size: cover;
    background-size: cover;
  }

</style>
