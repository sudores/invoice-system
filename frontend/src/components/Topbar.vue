<template>
  <header class="h-12 bg-gray-700 text-white flex items-center px-4 l-border">
    <button class="bg-gray-600 px-2 py-1 rounded mr-2" @click="$emit('toggle-sidebar')">
      ☰
    </button>
    <div>
        <p class="text-lg font-semibold"> Simple invoices </p>
    </div>
    <div class="mlauto logout bg-red-500 text-white ml-auto px-1 py-1 rounded">
      <button class="" @click="logout">Logout</button>
    </div>
  </header>
</template>

<script lang="ts">
</script>

<style scoped>
</style>
