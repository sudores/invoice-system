<script setup lang="ts">
import { ref } from 'vue'
import api from '../../api/client'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'

const email = ref('')
const password = ref('')
const error = ref('')

const router = useRouter()
const user = useUserStore()

async function login() {
  error.value = ''
  try {
    const res = await api.post('/api/v1/login', {
      email: email.value,
      password: password.value,
    })
    const token = res.data.jwt
    if (!token) {
      error.value = 'Invalid response from server'
      return
    }
    user.setToken(token)
    await router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Login failed'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="w-full max-w-sm p-8 bg-white rounded-lg shadow-md">
      <h2 class="text-2xl font-semibold text-center mb-6">Login</h2>
      
      <input
        v-model="email"
        type="text"
        placeholder="Email"
        class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <input
        v-model="password"
        type="password"
        placeholder="Password"
        class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <button
        @click="login"
        class="w-full py-2 bg-blue-500 text-white font-semibold rounded hover:bg-blue-600 transition-colors"
      >
        Log In
      </button>

      <p v-if="error" class="mt-4 text-red-500 text-sm text-center">{{ error }}</p>
    </div>
  </div>
</template>

