<template>
  <!-- eslint-disable max-len -->
  <DataTable ref="ingredientamounts" :value="ingredientamounts" v-model:selection="selectedIngredients" dataKey="ingredientid" :paginator="true" :rows="10" :rowsPerPageOptions="[5,10,25]" responsiveLayout="scroll">
    <template #header>
      <div class="table-header flex flex-column md:flex-row md:justiify-content-between">
        <Button label="New" icon="pi pi-plus" class="p-button-outlined p-button-success mr-2" @click="openNew" />
        <Button label="Delete" icon="pi pi-trash" class="p-button-outlined p-button-danger" @click="confirmDeleteSelected" :disabled="!selectedIngredients || !selectedIngredients.length" />
      </div>
    </template>

    <Column selectionMode="multiple" style="width: 3rem" :exportable="false">
    </Column>

    <Column field="ingredientid" header="Name" :sortable="true" :exportable="false" >
      <template #body="slotProps">
        {{ ingredients.find(x => x.id === slotProps.data.ingredientid).name }}
      </template>
    </Column>

    <Column field="amount" header="Amount" :sortable="false" >
    </Column>

    <Column field="unit" header="Unit" :sortable="false" >
    </Column>

    <Column :exportable="false" style="min-width:3rem;max-width:3rem">
      <template #body="slotProps">
        <Button icon="pi pi-pencil" class="p-button-outlined p-button-rounded p-button-text" @click="editIngredient(slotProps.data)" />
        <Button icon="pi pi-trash" class="p-button-outlined p-button-rounded p-button-text" @click="confirmDeleteIngredient(slotProps.data)" />
      </template>
    </Column>
  </DataTable>

  <Dialog v-model:visible="ingredientDialog" :style="{width: '450px'}" header="Ingredient Details" :modal="true" class="p-fluid">
    <span class="p-float-label" style="margin-top:30px; margin-bottom:30px">
      <Dropdown
        id="ingredient"
        v-model="this.ingredient.ingredientid"
        :options="ingredients"
        dataKey="id"
        optionLabel="name"
        optionValue="id"
      />
      <label for="ingredient">Ingredient</label>
    </span>

    <span class="p-float-label" style="margin-bottom:30px">
      <InputNumber
        id="amount"
        class="p-inputnumber-sm"
        showButtons
        v-model="this.ingredient.amount"
      />
      <label for="amount">Amount</label>
    </span>

    <span class="p-float-label">
      <InputText
        id="unit"
        class="p-inputtext-sm"
        v-model="this.ingredient.unit"
      />
      <label for="unit">Unit</label>
    </span>

    <template #footer>
      <CreateIngredient @createdIngredient="getIngredients"/>
      <Button label="Cancel" icon="pi pi-times" class="p-button-outlined p-button-text" @click="hideDialog"/>
      <Button label="Save" icon="pi pi-check" class="p-button-outlined p-button-text" @click="saveIngredient" />
    </template>
  </Dialog>

  <Dialog v-model:visible="deleteIngredientDialog" :style="{width: '450px'}" header="Confirm" :modal="true">
    <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="ingredient">Are you sure you want to delete <b>{{ ingredients.find(x => x.id === ingredient.ingredientid).name }}</b>?</span>
    </div>
    <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-text" @click="deleteIngredientDialog = false"/>
        <Button label="Yes" icon="pi pi-check" class="p-button-text" @click="deleteIngredient" />
    </template>
  </Dialog>

  <Dialog v-model:visible="deleteIngredientsDialog" :style="{width: '450px'}" header="Confirm" :modal="true">
    <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span v-if="selectedIngredients">Are you sure you want to delete the selected ingredients?</span>
    </div>
    <template #footer>
        <Button label="No" icon="pi pi-times" class="p-button-outlined p-button-text" @click="deleteIngredientsDialog = false"/>
        <Button label="Yes" icon="pi pi-check" class="p-button-outlined p-button-text" @click="deleteSelectedIngredients" />
    </template>
  </Dialog>
  <!-- eslint-enable max-len -->
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import { FilterMatchMode } from 'primevue/api';
import { IngredientAmount, Ingredients } from '@/lib/http/http';
import CreateIngredient from '@/components/CreateIngredient.vue';

@Options({
  components: {
    CreateIngredient,
    DataTable,
    Column,
    FilterMatchMode,
  },
  props: {
    ingredientamounts: {},
  },
  emits: ['updateIngredientAmounts'],
  data() {
    return {
      submitted: false,
      ingredient: null,
      ingredients: null,
      ingredientDialog: false,
      deleteIngredientDialog: false,
      deleteIngredientsDialog: false,
      selectedIngredients: null,
    };
  },
  mounted() {
    this.getIngredients();
  },
  methods: {
    getIngredients() {
      Ingredients.getAllIngredients().then((data) => { this.ingredients = data; });
    },
    openNew() {
      this.ingredient = {};
      this.submitted = false;
      this.ingredientDialog = true;
    },

    hideDialog() {
      this.ingredientDialog = false;
      this.submitted = false;
    },

    confirmDeleteIngredient(ingredient:IngredientAmount) {
      this.ingredient = ingredient;
      this.deleteIngredientDialog = true;
    },

    confirmDeleteSelected() {
      this.deleteIngredientsDialog = true;
    },

    editIngredient(ingredient:IngredientAmount) {
      this.ingredient = { ...ingredient };
      this.ingredientDialog = true;
    },

    saveIngredient() {
      this.submitted = true;
      const i = this.ingredientamounts.findIndex(
        (x:Record<string, unknown>) => x.ingredientid === this.ingredient.ingredientid,
      );
      if (i > -1) this.ingredientamounts[i] = this.ingredient;
      else this.ingredientamounts.push(this.ingredient);

      this.ingredientDialog = false;
      this.ingredient = {};
    },

    deleteIngredient() {
      this.$emit('updateIngredientAmounts', this.ingredientamounts.filter(
        (object:Record<string, unknown>) => object.ingredientid !== this.ingredient.ingredientid,
      ));

      this.deleteIngredientDialog = false;
      this.ingredient = {};
      this.$toast.add({
        severity: 'success', summary: 'Successful', detail: 'Ingredient Deleted', life: 3000,
      });
    },

    deleteSelectedIngredients() {
      this.$emit('updateIngredientAmounts', this.ingredientamounts.filter(
        (val:Record<string, unknown>) => !this.selectedIngredients.includes(val),
      ));
      this.deleteIngredientsDialog = false;
      this.selectedIngredients = null;
      this.$toast.add({
        severity: 'success', summary: 'Successful', detail: 'Ingredients Deleted', life: 3000,
      });
    },
  },
})
export default class IngredientEditor extends Vue {}
</script>

<style lang="scss" scoped>
  .element {
    margin-top: 30px;
    margin-bottom: 60px;
  }
</style>
