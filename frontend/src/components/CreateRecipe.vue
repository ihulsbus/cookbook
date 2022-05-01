<template>
  <Button class="p-button-outlined" label="New" icon="pi pi-plus" @click="toggleDialog()"/>
  <Dialog
    header="New recipe"
    :breakpoints="{'960px': '75vw', '640px': '100vw'}"
    :style="{width: '50vw'}"
    :modal="true"
    v-model:visible="visible"
  >
    <h5 class="text-center">Create Recipe</h5>
    <form class="p-fluid" @submit.prevent="onSubmit">
      <div class="field">
          <div class="p-float-label">
            <InputText
              id="recipetitle"
              type="text"
              v-model.trim="recipe.title"
              autofocus
            />
          </div>
          <small
            class="error"
            v-for="(error, index) of v$.title.$errors"
            :key="index"
          >
          {{ error.$message }}
        </small>
      </div>
      <Button type="cancel" label="Cancel" class="mt-2" />
      <Button type="submit" label="Submit" class="mt-2" />
    </form>
  </Dialog>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import { required, minLength, between } from '@vuelidate/validators';

@Options({
  components: {
  },
  data() {
    return {
      visible: false,
      recipe: {
        title: '',
        description: '',
        preptime: 0,
        cooktime: 0,
        persons: 0,
      },
    };
  },
  validations() {
    return {
      visible: false,
      recipe: {
        title: {
          required,
          minLength: minLength(2),
        },
        description: {
          required,
        },
        preptime: {
          required,
          between: between(1, 120),
        },
        cooktime: {
          required,
          between: between(1, 10080),
        },
        persons: {
          required,
          between: between(1, 10),
        },
      },
    };
  },
  methods: {
    onSubmit() {
      this.v$.$touch();
      if (this.v$.$error) return;
      console.log('Form is valid');
    },
    toggleDialog() {
      this.visible = !this.visible;
      if (!this.visible) {
        this.resetForm();
      }
    },
    resetForm() {
      this.recipe = {};
    },
  },
})

export default class RecipeCreate extends Vue {}
</script>

<style lang="scss" scoped>
</style>
