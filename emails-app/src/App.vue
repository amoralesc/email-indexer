<script setup lang="ts">
import { RouterView } from 'vue-router'
import SearchBar from './components/SearchBar.vue'
</script>

<template>
  <header>
    <div class="title">
      <img alt="App logo" class="logo" src="@/assets/logo.svg" width="30" height="30" />
      <h1>Email Indexer</h1>
    </div>
    <SearchBar />
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
