import axios, { AxiosResponse } from 'axios';
import idsrvAuth from '@/lib/auth/auth';

const defaultOptions = {
  baseURL: `${process.env.VUE_APP_COOKBOOK_BACKEND_URL}/v1`,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
};

const instance = axios.create(defaultOptions);

instance.interceptors.request.use((config) => {
  const conf = config;
  const token = idsrvAuth.accessToken;
  conf.headers.Authorization = token ? `Bearer ${token}` : '';
  return conf;
});

export interface Ingredient {
  id: number;
  name: string;
}

export interface Tag {
  id: number;
  name: string;
}

export interface IngredientAmount {
  recipeid: number;
  ingredientid: number;
  amount: number;
  unit: string;
}

export interface Recipe {
  ID: number;
  Title: string;
  Description: string;
  Method: string;
  Preptime: number;
  Cooktime: number;
  Persons: number;
  Ingredients: Array<Ingredient>;
  Ingredientamounts: Array<IngredientAmount>;
  Tags: Array<Tag>;
}

const responseBody = (response: AxiosResponse) => response.data;

const requests = {
  get: (url: string) => instance.get(url).then(responseBody),
  post: (url: string, body: string) => instance.post(url, body).then(responseBody),
  put: (url: string, body: string) => instance.put(url, body).then(responseBody),
  delete: (url: string, body: string) => instance.delete(url, { data: body }).then(responseBody),
};

export const Recipes = {
  getAllRecipes: (): Promise<Recipe[]> => requests.get('/recipe'),
  getSingleRecipe: (id: number): Promise<Recipe> => requests.get(`/recipe/${id}`),
  createRecipe: (item: Recipe): Promise<Recipe> => requests.post('/recipe', JSON.stringify(item)),
  updateRecipe: (item: Recipe): Promise<Recipe> => requests.put('/recipe', JSON.stringify(item)),
};

export const Ingredients = {
  getAllIngredients: (): Promise<Ingredient[]> => requests.get('/ingredients'),
  getSingleIngredient: (ingredientID: number): Promise<Ingredient[]> => requests.get(`/ingredients/${ingredientID}`),
  createIngredient: (item: Ingredient): Promise<Ingredient> => requests.post('/ingredients', JSON.stringify(item)),
  deleteIngredient: (item: Ingredient): Promise<Ingredient> => requests.delete('/ingredients', JSON.stringify(item)),
};
