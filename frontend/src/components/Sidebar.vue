<script lang="ts" setup>
import { ref, defineProps, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// Receive the prop from parent
const props = defineProps<{ isOpen: boolean }>()

// Router instance
const router = useRouter()

// Menu pages and selected page
const pages = [
  { label: 'Dashboard', route: '/' },
  { label: 'Create Invoice', route: '/invoice/create' },
  { label: 'Received Invoices', route: '/invoice/list' }
]
const selectedPage = ref('Dashboard')

// Emit for mobile toggle
const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
}>()

function selectPage(page: { label: string; route: string }) {
  selectedPage.value = page.label
  console.log('Selected page:', page.label)

  // Navigate to the route
  router.push(page.route)

  // Auto-close sidebar on mobile
  if (window.innerWidth < 768) {
    emit('toggle-sidebar')
  }
}
onMounted(() => {
  if (window.innerWidth < 768 && props.isOpen) {
    emit('toggle-sidebar')
  }
})
</script>

<template>
  <aside
    v-show="props.isOpen"
    class="w-52 bg-gray-700 text-white p-4 flex flex-col"
  >
    <ul>
      <li
        v-for="page in pages"
        :key="page.label"
        @click="selectPage(page)"
        :class="[
          'cursor-pointer my-2 px-2 py-1 rounded',
          page.label === selectedPage ? 'bg-gray-500' : 'hover:bg-gray-600'
        ]"
      >
        {{ page.label }}
      </li>
    </ul>
  </aside>
</template>

