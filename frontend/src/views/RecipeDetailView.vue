<template>
  <div class="surface-ground px-4 py-4 md:px-6 lg:px-8">
    <div class="flex flex-row flex-wrap">
      <div class="flex align-items-center backButton">
        <Button
          class="p-button-rounded p-button-text"
          icon="pi pi-arrow-left"
          @click="this.$router.push({ name: 'RecipeView' });" />
      </div>
      <div class="flex align-items-center  editButton">
        <Button
          class="p-button-rounded p-button-text"
          icon="pi pi-pencil"
          @click="this.$router.push({ name: 'RecipeEdit', params: { id: recipe.id } });" />
      </div>
    </div>
    <div class="grid grid-nogutter surface-section text-800">
      <div class="col-12 md:col-6 p-6 text-center md:text-left flex align-items-center ">
        <section>
          <span class="block text-6xl font-bold mb-1">{{ recipe.title }}</span>
          <p class="mt-0 mb-4 text-700 line-height-3">{{ recipe.description }}</p>
          <div class="p-grid p-formgrid p-fluid mt-0 mb-4 text-700 line-height-3">
            <span class="p-col-1">
              <img :src="require(`@/assets/icons/knife1.svg`)" class="p-mr-5" style="height: 20px"/>
              {{ recipe.preptime }} Minutes
            </span>
            <span class="p-col-1">
              <img :src="require(`@/assets/icons/pot.svg`)" class="p-mr-5" style="height: 20px"/>
              {{ recipe.cooktime }} Minutes
            </span>
            <span class="p-col-1">
              <i class="pi pi-clock p-mr-5"/>
              {{ recipe.preptime+recipe.cooktime }} Minutes
            </span>
          </div>
        </section>
      </div>
      <div class="col-12 md:col-6 overflow-hidden">
        <img
          :src="require(`@/assets/images/placeholder/${image}`)"
          alt="Image"
          class="md:ml-auto block md:h-full"
          style="clip-path: polygon(8% 0, 100% 0%, 100% 100%, 0 100%)">
      </div>
    </div>

    <div class="grid">
      <div class="col-12 lg:col-3">
        <div class="p-3 h-full">
          <div class="shadow-2 p-3 h-full flex flex-column surface-card" style="border-radius: 6px">
            <div class="text-900 font-medium text-xl mb-2">Ingredients</div>
            <hr class="my-3 mx-0 border-top-1 border-none surface-border" />
            <table>
              <div v-for="ingredient in recipe.ingredientamounts" :key="ingredient.ingredientid">
                  <tr>
                    <td class="name">
                      {{ recipe.ingredients.find(x => x.id === ingredient.ingredientid).name }}
                    </td>
                    <td class="amount">
                      {{ ingredient.amount }}
                    </td>
                    <td class="unit">
                      {{ ingredient.unit }}
                    </td>
                  </tr>
              </div>
            </table>
          </div>
        </div>
      </div>

      <div class="col-12 lg:col-9">
        <div class="p-3 h-full">
          <div class="shadow-2 p-3 h-full flex flex-column surface-card" style="border-radius: 6px">
            <div class="text-900 font-medium text-xl mb-2">Method</div>
            <hr class="my-3 mx-0 border-top-1 border-none surface-border" />
            <div v-html="recipe.method" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Image from 'primevue/image';
import {
  Recipes,
} from '@/lib/http/http';

@Options({
  components: {
    Image,
  },
  props: {
    recipeID: 0,
  },
  data() {
    return {
      recipe: {},
      image: 'chicken.jpg',
    };
  },
  mounted() {
    Recipes.getSingleRecipe(this.$route.params.id).then((data) => { this.recipe = data; });
  },
})
export default class RecipeDetailView extends Vue {}
</script>

<style lang="scss" scoped>
  .surface-ground {
    min-height: 100%;
  }

  .recipeDetail {
    padding-top: 3vh;
    padding-left: 3vw;
    padding-right: 3vw;
  }

  .headerTitle {
    color: white;
    position: absolute;
    top: 50%;
    left: 2%;
  }

  .name {
    text-align:left;
  }

  .amount {
    text-align: right;
  }

</style>
