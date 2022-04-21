<template>
  <div class="surface-ground px-4 py-4 md:px-6 lg:px-8">
    <div class="text-3xl font-medium text-900 header">
      <span class="title">Recipes</span> <span class="buttons"><CreateRecipe /></span>
    </div>
    <hr class="my-3 mx-0 border-top-1 border-none surface-border" />
    <RecipeTable :recipes="recipes" class="recipeTable"/>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import { Recipes } from '@/lib/http/http';
import RecipeTable from '@/components/RecipeTable.vue';
import CreateRecipe from '@/views/CreateRecipe.vue';

@Options({
  components: {
    CreateRecipe,
    RecipeTable,
  },
  data() {
    return {
      recipes: {
        type: Object,
      },
    };
  },
  mounted() {
    Recipes.getAllRecipes().then((data) => { this.recipes = data; });
  },
})
export default class RecipeView extends Vue {}
</script>

<style lang="scss" scoped>
  .surface-ground {
    /* mobile viewport bug fix */
    min-height: 100%;
  }

  .header {
    width: 100%;
  }

  .buttons {
    margin-right: 0;
  }

  .recipeTable {
    padding-top: 2vh;
  }
</style>
