<script setup lang="ts">
import { useRouter } from 'vue-router'

import { useEmailsStore } from '@/stores'

import SearchIcon from '../components/icons/IconSearch.vue'
import TuneIcon from '../components/icons/IconTune.vue'
import CloseIcon from '../components/icons/IconClose.vue'

const router = useRouter()
const emailsStore = useEmailsStore()

const onInput = (e: Event) => {
  emailsStore.setQueryString((e.target as HTMLInputElement).value)
  console.log(emailsStore.getQueryString())
}

const onSearch = () => {
  if (emailsStore.queryString === '') {
    return
  }
  emailsStore.fetch()
}

const onClose = async () => {
  emailsStore.setQueryString('')
  await emailsStore.initialize()
  if (router.currentRoute.value.name !== 'home') {
    router.push({ name: 'home' })
  }
}
</script>

<template>
  <div class="search-bar">
    <i>
      <SearchIcon @click="onSearch" />
    </i>
    <input
      type="text"
      placeholder="Search"
      v-model="emailsStore.queryString"
      @value="emailsStore.queryString"
      @input="onInput"
      @keydown.enter="onSearch"
    />
    <i v-if="emailsStore.getQueryString() !== ''" @click="onClose">
      <CloseIcon />
    </i>
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
