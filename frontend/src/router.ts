import { createRouter, createWebHistory } from 'vue-router'
import Login from './components/pages/Login.vue'
import Dashboard from './components/pages/Dashboard.vue'
import CreateInvoice from './components/pages/CreateInvoice.vue'
import ListInvoice from './components/pages/ListInvoice.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login, name: 'Login' },
    { path: '/', component: Dashboard, name: 'Dashboard' },
    { path: '/invoice/create', component: CreateInvoice, name: 'CreateInvoice' },
    { path: '/invoice/list', component: ListInvoice, name: 'ListInvoice' },
  ]
})

export default router
