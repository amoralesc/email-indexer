<script setup lang="ts">
import TabItem from './TabItem.vue'

import type { Tab } from '@/models/tab'

const props = defineProps<{
  tabs: Tab[]
  selectedTab: string
}>()

const isSelected = (tabName: string) => {
  return props.selectedTab === tabName
}
const getIcon = (tab: Tab) => {
  return isSelected(tab.tabName) ? tab.iconFilled : tab.icon
}
</script>

<template>
  <div class="tab-bar">
    <TabItem
      v-for="t in tabs"
      :key="t.tabName"
      :label="t.label"
      :selected="isSelected(t.tabName)"
      @click="$emit('selectTab', t.tabName)"
    >
      <component :is="getIcon(t)" />
    </TabItem>
  </div>
</template>

<style scoped>
.tab-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  width: 100%;
  height: 3rem;
  border-bottom: 1px solid var(--color-border);
}
</style>
