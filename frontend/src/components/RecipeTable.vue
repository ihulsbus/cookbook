<template>
  <div class="container">
    <div class="grid">
      <div
      class="card"
      v-for="recipe in recipes"
      :key="recipe.id"
      v-on:click="onClick(recipe.id)">
          <img alt="Header image" :src="require(`@/assets/images/placeholder/${image}`)">
          <span class="title">{{ recipe.title }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';

@Options({
  props: {
    recipes: {
      type: Object,
    },
  },
  data() {
    return {
      image: 'chicken.jpg',
      value: [],
    };
  },
  methods: {
    onClick(recipeID:number) {
      this.$router.push({ name: 'RecipeDetailView', params: { id: recipeID } });
    },
  },
})
export default class RecipeTable extends Vue {}
</script>

<style lang="scss" scoped>
  .grid {
  display: grid;
  grid-gap: 1rem;

  .card {
    &:hover {
      img {
        transform: scale(1.1);
        filter: blur(0px);
      }
    }
    cursor: pointer;
    border-radius: 0.5rem;
    position: relative;
    overflow: hidden;
    border: 0px solid #ccc;
    box-sizing: border-box;
    img {
      transition: transform 0.25s ease-in-out, filter 0.25s ease-in-out;
      filter: blur(3px);
      display: block;
      width: 100%;
      box-sizing: border-box;
    }
    span.title {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;
      width: 100%;
      padding: 1rem;
      color: #fff;
      font-size: 1.7rem;
      text-shadow: 0px 0px 7px rgba(0, 0, 0, 0.9);
      box-sizing: border-box;
    }
  }
}

.container {
  margin: 0 auto;
}

@media (min-width: 480px) {
  .container {
    max-width: 100%;
  }
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (min-width: 600px) {
  .container {
    max-width: 100%;
  }
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (min-width: 1025px) {
  .container {
    max-width: 1024px;
  }
  .grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (min-width: 1400px) {
  .container {
    max-width: 1400px;
  }
  .grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

html {
  margin: 0;
}

body {
  padding-top: 6rem;
  padding-bottom: 6rem;
}

</style>
