<template>
  <div class="h-screen flex">
    <Sidebar v-if="!isLoginPage" @toggle-sidebar="toggleSidebar" :is-open="isSidebarOpen" class="w-52" />
    <div class="flex flex-col flex-1">
      <Topbar v-if="!isLoginPage" @toggle-sidebar="toggleSidebar" class="h-12" />
      <main class="flex-1 bg-gray-100 p-4 overflow-auto">
        <router-view />
      </main>

    </div>

  </div>
</template>

<script setup>
import Sidebar from './components/Sidebar.vue'
import Topbar from './components/Topbar.vue'
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const isSidebarOpen = ref(true)

const toggleSidebar = () => {
  console.log("Sidebar toggle")
  isSidebarOpen.value = !isSidebarOpen.value
}

const route = useRoute()
const isLoginPage = computed(() => route.name === 'Login')
</script>
