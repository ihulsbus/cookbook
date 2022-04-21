<template>
  <div class="flex header">
    <div class="show-mobile">
      <Button
        type="button"
        icon="pi pi-bars"
        @click="toggle"
        aria-haspopup="true"
        aria-controls="overlay_menu"/>
      <Menu
        id="overlay_menu"
        ref="menu"
        :model="menuOptions"
        :popup="true" />
    </div>
    <div
    class="flex-initial flex align-items-center justify-content-center title"
    v-on:click="goHome()">
      <img
        src="~@/assets/images/logos/logo-placeholder.png"
        alt="Logo Image"
        height="50"
        class="mb-3 logo">
    </div>
    <div class="hide-mobile">
      <TabMenu :model="menuOptions" />
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import TabMenu from 'primevue/tabmenu';
import Menu from 'primevue/menu';

@Options({
  components: {
    TabMenu,
    Menu,
  },
  data() {
    return {
      menuOptions: [
        {
          label: 'Recipes',
          to: '/app/recipes',
        },
        {
          label: 'Ingredients',
          to: '/app/ingredients',
        },
        {
          label: 'Shopping Lists',
          to: '/app/shopping',
        },
        {
          label: 'Mealplanner',
          to: '/app/planner',
        },
      ],
    };
  },
  methods: {
    goHome() {
      this.$router.push({ name: 'HomeView' });
    },
    toggle(event:Event) {
      this.$refs.menu.toggle(event);
    },
  },
})
export default class HeaderBar extends Vue {}
</script>

<style lang="scss" scoped>
  @media only screen and (max-width: 620px) {
    .show-mobile {
      display: block;
    }
    .hide-mobile {
      display: none;
    }
  }

  @media only screen and (min-width: 621px) {
    .show-mobile {
      display: none;
    }
    .hide-mobile {
      display: block;
    }
  }

  .header {
    height: 55px;
    box-shadow: 0px 1px 5px 0px;
    width: 100%;
    z-index: 10;
    padding-top: 5px;
    padding-left: 5px;
    background-color: white;
  }

  .title {
    margin-right: 5px;
    margin-bottom: 5px;
  }

  .logo {
    position: relative;
    top: 6px;
  }

  .nav {
    padding-right: 5px;
    margin-bottom: 5px;

    a {
      margin-left: 5px;
      font-weight: bold;
      color: #2c3e50;

      &.router-link-exact-active {
        color: #42b983;
      }
    }
  }

</style>
