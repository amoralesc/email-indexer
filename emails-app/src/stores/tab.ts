import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useTabStore = defineStore('tab', () => {
  const tab = ref('all')
  function selectTab(tabName: string) {
    tab.value = tabName
  }

  return { tab, selectTab }
})
