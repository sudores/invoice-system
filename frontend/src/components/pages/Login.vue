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
  <div class="login-container">
    <div class="login-box">
      <h2>Login</h2>
      <input v-model="email" type="text" placeholder="Email" />
      <input v-model="password" type="password" placeholder="Password" />
      <button @click="login">Log In</button>
      <p class="error" v-if="error">{{ error }}</p>
    </div>
  </div>
</template>

