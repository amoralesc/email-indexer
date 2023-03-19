<script setup lang="ts">
import TabItem from './TabItem.vue'
import InboxIcon from './icons/IconInbox.vue'
import InboxFilledIcon from './icons/IconInboxFilled.vue'
import StarIcon from './icons/IconStar.vue'
import StarFilledIcon from './icons/IconStarFilled.vue'
import { useTabStore } from '../stores/tab'

const store = useTabStore()

const tabs = [
  {
    label: 'All',
    icon: InboxIcon,
    iconFilled: InboxFilledIcon,
    tab: 'all'
  },
  {
    label: 'Starred',
    icon: StarIcon,
    iconFilled: StarFilledIcon,
    tab: 'starred'
  }
]

const isSelected = (tabName: string) => {
  return store.tab === tabName
}

const getIcon = (tabName: string) => {
  const { icon, iconFilled } = tabs.find((t) => t.tab === tabName) || tabs[0]
  return isSelected(tabName) ? iconFilled : icon
}
</script>

<template>
  <div class="tab-bar">
    <TabItem
      v-for="t in tabs"
      :key="t.tab"
      :label="t.label"
      :selected="isSelected(t.tab)"
      @click="() => store.selectTab(t.tab)"
    >
      <component :is="getIcon(t.tab)" />
    </TabItem>
  </div>
</template>

<style scoped>
.tab-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  height: 48px;
  border-bottom: 1px solid var(--color-border);
}
</style>
