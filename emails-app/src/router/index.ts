import { createRouter, createWebHistory } from 'vue-router'
import ListView from '../views/ListView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/emails',
      name: 'home',
      component: ListView
    },
    {
      path: '/emails/:id',
      name: 'email',
      component: () => import('../views/EmailView.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/emails'
    }
  ]
})

export default router
