import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import idsrvAuth from '@/lib/auth/auth';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Splashscreen',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "splashscreen" */ '@/views/Splashscreen.vue'),
  },
  {
    path: '/app',
    name: 'Cookbook',
    component: () => import(/* webpackChunkName: "cookbook" */ '@/views/Cookbook.vue'),
    children: [
      {
        path: '/app/home',
        name: 'HomeView',
        component: () => import(/* webpackChunkName: "home" */ '@/views/HomeView.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/recipes',
        name: 'RecipeView',
        component: () => import(/* webpackChunkName: "recipes" */ '@/views/RecipeView.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/recipes/:id',
        name: 'RecipeDetailView',
        component: () => import(/* webpackChunkName: "recipe" */ '@/views/RecipeDetailView.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/recipes/:id/edit',
        name: 'RecipeEdit',
        component: () => import(/* webpackChunkName: "recipeEdit" */ '@/views/RecipeEdit.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/ingredients',
        name: 'IngredientView',
        component: () => import(/* webpackChunkName: "ingredients" */ '@/views/IngredientView.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/shopping',
        name: 'ShoppingView',
        component: () => import(/* webpackChunkName: "shopping" */ '@/views/ShoppingView.vue'),
        meta: {
          authName: 'main',
        },
      },
      {
        path: '/app/planner',
        name: 'PlannerView',
        component: () => import(/* webpackChunkName: "planner" */ '@/views/PlannerView.vue'),
        meta: {
          authName: 'main',
        },
      },
    ],
    meta: {
      authName: 'main',
    },
  },
  {
    path: '/logout',
    name: 'LogoutView',
    component: () => import(/* webpackChunkName: "logout" */ '@/views/LogoutView.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

idsrvAuth.useRouter(router);

export default router;
