<template>
  <div class="surface-ground px-4 py-4 md:px-6 lg:px-8">
    <div class="flex flex-row flex-wrap">
      <div class="flex align-items-center cancelButton">
        <Button
          class="p-button-rounded p-button-text"
          icon="pi pi-times"
          @click="cancelPopup=true" />
      </div>
      <div class="flex align-items-center saveButton">
        <Button
          class="p-button-rounded p-button-text"
          icon="pi pi-check"
          @click="saveRecipeUpdate" />
      </div>
    </div>
    <div class="grid grid-nogutter surface-section text-800">
      <div class="col-12 md:col-6 p-6 text-center md:text-left flex align-items-center ">
        <section>
          <span class="block text-6xl font-bold mb-1">
            <InputText
              type="text"
              class="p-inputtext-lg"
              placeholder="Recipe title"
              v-model="recipe.title"
            />
          </span>
          <p class="mt-0 mb-4 text-700 line-height-3">
            <InputText
              type="text"
              class="p-inputtext-lg inputBox"
              placeholder="Recipe description"
              v-model="recipe.description"
            />
          </p>
          <div class="p-grid p-formgrid p-fluid mt-0 mb-4 text-700 line-height-3">
            <span class="p-col-1">
              <img :src="require(`@/assets/icons/knife1.svg`)" class="p-mr-5" style="height: 20px"/>
              <InputNumber
                class="p-inputnumber-sm"
                placeholder="Preparation time"
                showButtons
                suffix=" Minutes"
                v-model="recipe.preptime"
              />
            </span>
            <span class="p-col-1">
              <img :src="require(`@/assets/icons/pot.svg`)" class="p-mr-5" style="height: 20px"/>
              <InputNumber
                class="p-inputnumber-sm"
                placeholder="Cooking time"
                showButtons
                suffix=" Minutes"
                v-model="recipe.cooktime"
              />
            </span>
            <span class="p-col-1">
              <i class="pi pi-clock p-mr-5"/>
              {{ recipe.preptime+recipe.cooktime }} Minutes
            </span>
            <span class="p-col-1">
              <i class="pi pi-user p-mr-5"/>
              <InputNumber
                class="p-inputnumber-sm"
                placeholder="Amount of persons"
                showButtons
                suffix=" Persons"
                v-model="recipe.persons"
              />
            </span>
          </div>
          {{ recipe.tags }}
          <Chips v-model="recipe.tags">
            <template #chip="slotProps">
              <div>
                <span>{{slotProps.value.name}}</span>
              </div>
            </template>
          </Chips>
        </section>
      </div>
      <div class="col-12 md:col-6 overflow-hidden">
        <img
          :src="require(`@/assets/images/placeholder/${image}`)"
          alt="Image"
          class="md:ml-auto block md:h-full cover-image"
          style="clip-path: polygon(8% 0, 100% 0%, 100% 100%, 0 100%)">
      </div>
    </div>

    <div class="grid">
      <div class="col-12 lg:col-6">
        <div class="p-3 h-full">
          <div class="shadow-2 p-3 h-full flex flex-column surface-card" style="border-radius: 6px">
            <div class="text-900 font-medium text-xl mb-2">Ingredients</div>
            <hr class="my-3 mx-0 border-top-1 border-none surface-border" />
            <IngredientEditor
              :ingredientamounts="recipe.ingredientamounts"
              @updateIngredientAmounts="updateIngredientAmounts"
            />
          </div>
        </div>
      </div>

      <div class="col-12 lg:col-6">
        <div class="p-3 h-full">
          <div class="shadow-2 p-3 h-full flex flex-column surface-card" style="border-radius: 6px">
            <div class="text-900 font-medium text-xl mb-2">Method</div>
            <hr class="my-3 mx-0 border-top-1 border-none surface-border" />
            <Editor v-model="recipe.method" editorStyle="height: 320px"/>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog v-model:visible="cancelPopup" :style="{width: '450px'}" header="Confirm" :modal="true">
    <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span>Are you sure you want to cancel?</span>
    </div>
    <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="cancelPopup = false"/>
        <Button label="Yes" icon="pi pi-check" class="p-button-text"
          @click="$router.push({ name: 'RecipeDetailView', params: { id: this.$route.params.id } })"
        />
    </template>
  </Dialog>

</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Editor from 'primevue/editor';
import Chips from 'primevue/chips';
import { Recipes } from '@/lib/http/http';
import IngredientEditor from '@/components/IngredientEditor.vue';

@Options({
  components: {
    Editor,
    Chips,
    IngredientEditor,
  },
  props: {
    recipeID: 0,
  },
  data() {
    return {
      recipe: {},
      cancelPopup: false,
      image: 'chicken.jpg',
    };
  },
  mounted() {
    Recipes.getSingleRecipe(this.$route.params.id).then((data) => { this.recipe = data; });
  },
  methods: {
    updateIngredientAmounts(arg:Record<string, unknown>) {
      this.recipe.ingredientamounts = arg;
    },

    cancelRecipeUpdate() {
      this.cancelPopup = true;
    },

    saveRecipeUpdate() {
      Recipes.updateRecipe(this.recipe).then(
        this.saveSuccess(),
        this.saveFailed(),
      );
    },
    saveSuccess() {
      this.$toast.add({
        severity: 'success', summary: 'Update successful', detail: 'Recipe updated', life: 3000,
      });
      this.$router.push({ name: 'RecipeDetailView', params: { id: this.$route.params.id } });
    },
    saveFailed() {
      this.$toast.add({
        severity: 'error', summary: 'Update failed', detail: 'Recipe update failed', life: 3000,
      });
    },
  },
})
export default class RecipeEdit extends Vue {}
</script>

<style lang="scss" scoped>
  .surface-ground {
    min-height: 100%;
  }

  .p-inputnumber-sm {
    max-width: 150px;
    margin-left: 2px;
    margin-right: 10px;
  }

  input {
    width: 450px;
  }
</style>
