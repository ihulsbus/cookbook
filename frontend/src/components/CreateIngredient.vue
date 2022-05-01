<template>
  <Button
    label="New"
    icon="pi pi-plus"
    class="p-button-outlined p-button-success mr-2"
    @click="openDialog"
  />
  <Dialog
    header="New ingredient"
    :breakpoints="{'960px': '75vw', '640px': '100vw'}"
    :style="{width: '50vw'}"
    :modal="true"
    v-model:visible="visible"
  >
  <div class="p-fluid">
    <div class="p-field">
      <label for="name">Ingredient Name</label>
      <InputText
        id="name"
        type="text"
        v-model="ingredient.name"
        :class="{'p-invalid': validationErrors.name && submitted}"
        autofocus
      />
      <small
        v-show="validationErrors.name && submitted"
        class="p-error">
          Ingredient name is required.
      </small>
    </div>
  </div>
  <template #footer>
    <Button class="p-button-outlined" label="Cancel" icon="pi pi-cross" @click="closeDialog()"/>
    <Button class="p-button-outlined" label="Create" icon="pi pi-check" @click="validateForm()"/>
  </template>
  </Dialog>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Dialog from 'primevue/dialog';
import { Ingredients } from '@/lib/http/http';

@Options({
  components: {
    Dialog,
  },
  emits: ['createdIngredient'],
  data() {
    return {
      visible: false,
      ingredient: {
      },
      validationErrors: {
        type: Object,
      },
      isFormValid: false,
    };
  },
  methods: {
    openDialog() {
      this.visible = true;
    },
    closeDialog() {
      this.visible = false;
    },
    createIngredient() {
      Ingredients.createIngredient(this.ingredient).then(
        this.$toast.add({
          severity: 'success', summary: 'Ingredient created', life: 3000,
        }),
      );
      this.ingredient = {};
      this.$emit('createdIngredient');
    },
    validateForm() {
      if (!this.ingredient.name.trim()) {
        this.validationErrors.name = true;
      } else {
        delete this.validationErrors.name;
        this.isFormValid = true;
        this.visible = false;
        this.createIngredient();
      }
    },
  },
})
export default class CreateIngredient extends Vue {}
</script>

<style lang="scss" scoped>

</style>
