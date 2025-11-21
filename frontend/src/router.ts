import { createRouter, createWebHistory } from 'vue-router'
import Login from './components/pages/Login.vue'
//import Dashboard from './pages/Dashboard.vue' // create placeholder

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    //{ path: '/', component: Dashboard },
  ]
})

export default router
