<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'

import { useEmailsStore } from './stores'

import SearchBarView from './views/SearchBarView.vue'

const router = useRouter()
const emailsStore = useEmailsStore()

const onClickTitle = async () => {
  await emailsStore.initialize()
  router.push({ name: 'home' })
}
</script>

<template>
  <header>
    <div class="title" @click="onClickTitle()">
      <img alt="App logo" class="logo" src="@/assets/logo.svg" width="30" height="30" />
      <h1>Email Indexer</h1>
    </div>
    <SearchBarView />
  </header>

  <RouterView />

  <footer>
    <p>Email Indexer</p>
  </footer>
</template>

<style scoped>
header {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  margin-bottom: 1rem;
  padding: 0 0.5rem;
}

footer {
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 0.75rem;
  font-family: 'IBM Plex Mono', monospace;
  font-size: 0.8rem;
}

.title {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.logo {
  display: block;
  margin-top: 0.2rem;
}

h1 {
  font-size: 1.5rem;
  letter-spacing: -0.01rem;
  margin-left: 0.5rem;
  color: var(--color-heading);
}
</style>
