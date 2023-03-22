import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useTabStore = defineStore('tab', () => {
  const selectedTab = ref('all')
  function setTab(tabName: string) {
    selectedTab.value = tabName
  }

  return { selectedTab, setTab }
})
