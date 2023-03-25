<script setup lang="ts">
import { ref } from 'vue'

import SearchIcon from './icons/IconSearch.vue'
import TuneIcon from './icons/IconTune.vue'

// get text from input field
const input = ref('')
const onInput = (e: Event) => {
  input.value = (e.target as HTMLInputElement).value
}
</script>

<template>
  <!-- This is a search bar component. It renders a search logo and an input field. -->

  <div class="search-bar">
    <i @click="$emit('search', input)">
      <SearchIcon />
    </i>
    <input
      type="text"
      placeholder="Search"
      v-model="input"
      @input="onInput"
      @keyup.enter="$emit('search', input)"
    />
    <i>
      <TuneIcon />
    </i>
  </div>
</template>

<style scoped>
.search-bar {
  flex-grow: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  height: 2.5rem;
  background-color: var(--color-background-soft);
  border-radius: 0.5rem;
  padding: 0 0.6rem;
  max-width: 50rem;
}

.search-bar input {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  background-color: transparent;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-heading);
}

i {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  cursor: pointer;
}

i:hover {
  color: var(--color-primary);
}
</style>
