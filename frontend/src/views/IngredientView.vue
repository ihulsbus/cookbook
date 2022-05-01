<template>
  <div>
    <h1>Ingredients</h1>
    <DataTable
      :value="ingredients"
      :paginator="true"
      :rows="10"
      v-model:filters="filters"
      v-model:selection="selectedIngredients"
      filterDisplay="menu"
      :loading="loading"
      responsiveLayout="scroll"
      :globalFilterFields="['id','name']"
      dataKey="id">
      <template #header>
        <div class="flex justify-content-between">
            <span>
              <CreateIngredient @createdIngredient="getIngredients"/>
              <Button
                label="Delete"
                icon="pi pi-trash"
                class="p-button-outlined p-button-danger mr-2"
                @click="confirmDeleteSelected"
                :disabled="!selectedIngredients || !selectedIngredients.length"
              />
              <Button
                type="button"
                icon="pi pi-filter-slash"
                label="Clear"
                class="p-button-outlined"
                @click="clearFilter()"
              />
            </span>
            <span class="p-input-icon-left">
                <i class="pi pi-search" />
                <InputText v-model="filters['global'].value" placeholder="Keyword Search" />
            </span>
        </div>
      </template>
      <template #empty>
        No ingredients found.
      </template>
      <template #loading>
        Loading ingredient data. Please wait.
      </template>
      <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
      <Column field="id" header="ID"></Column>
      <Column field="name" header="Name"></Column>
      <Column :exportable="false" style="min-width:3rem;max-width:3rem">
      <template #body="slotProps">
        <Button
          icon="pi pi-trash"
          class="p-button-outlined p-button-rounded p-button-text"
          @click="confirmDeleteIngredient(slotProps.data)"
        />
      </template>
    </Column>
    </DataTable>

    <Dialog
      v-model:visible="deleteIngredientDialog"
      :style="{width: '450px'}"
      header="Confirm"
      :modal="true">
      <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
          <span v-if="ingredient">
            Are you sure you want to delete
            <b>
              {{ ingredient.name }}
            </b>?
          </span>
      </div>
      <template #footer>
          <Button
            label="No"
            icon="pi pi-times"
            class="p-button-text"
            @click="deleteIngredientDialog = false"
          />
          <Button
            label="Yes"
            icon="pi pi-check"
            class="p-button-text"
            @click="deleteIngredient"
          />
      </template>
    </Dialog>

    <Dialog
      v-model:visible="deleteIngredientsDialog"
      :style="{width: '450px'}"
      header="Confirm"
      :modal="true"
    >
      <div class="confirmation-content">
          <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
          <span
            v-if="selectedIngredients"
          >
            Are you sure you want to delete the selected ingredients?
          </span>
      </div>
      <template #footer>
          <Button
            label="No"
            icon="pi pi-times"
            class="p-button-outlined p-button-text"
            @click="deleteIngredientsDialog = false"
          />
          <Button
            label="Yes"
            icon="pi pi-check"
            class="p-button-outlined p-button-text"
            @click="deleteIngredient"
          />
      </template>
    </Dialog>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import { FilterMatchMode } from 'primevue/api';
import { Ingredient, Ingredients } from '@/lib/http/http';
import CreateIngredient from '@/components/CreateIngredient.vue';

@Options({
  components: {
    DataTable,
    Column,
    ColumnGroup,
    CreateIngredient,
    FilterMatchMode,
  },
  data() {
    return {
      ingredients: [],
      loading: true,
      filters: null,
      deleteIngredientDialog: false,
      deleteIngredientsDialog: false,
      selectedIngredients: null,
    };
  },
  created() {
    this.initFilters();
  },
  mounted() {
    this.getIngredients();
  },
  methods: {
    getIngredients() {
      Ingredients.getAllIngredients().then(
        (data) => { this.ingredients = data; this.loading = false; },
      );
    },
    clearFilter() {
      this.initFilters();
    },
    initFilters() {
      this.filters = {
        global: { value: null, matchMode: FilterMatchMode.CONTAINS },
        id: { value: null, matchMode: FilterMatchMode.CONTAINS },
        name: { value: null, matchMode: FilterMatchMode.CONTAINS },
      };
    },
    confirmDeleteIngredient(ingredient:Ingredient) {
      this.ingredient = ingredient;
      this.deleteIngredientDialog = true;
    },

    confirmDeleteSelected() {
      this.deleteIngredientsDialog = true;
    },
    deleteIngredient() {
      Ingredients.deleteIngredient(this.ingredient).then(
        this.updateSuccess(),
        this.updateFailed(),
      );
    },
    updateSuccess() {
      this.$toast.add({
        severity: 'success', summary: 'Update successful', life: 3000,
      });
    },
    updateFailed() {
      this.$toast.add({
        severity: 'error', summary: 'Update failed', life: 3000,
      });
    },
  },
})
export default class IngredientView extends Vue {}
</script>
