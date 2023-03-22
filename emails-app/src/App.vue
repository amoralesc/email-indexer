<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <header>
    <div class="title">
      <img alt="App logo" class="logo" src="@/assets/logo.svg" width="30" height="30" />
      <h1>Email Indexer</h1>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0rem 1rem;
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
