<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { useEmailsStore } from '@/stores'
import { useTabStore } from '@/stores'
import type { Tab } from '@/models/tab'

import EmailItem from '@/components/EmailItem.vue'
import TabBar from '@/components/TabBar.vue'
import NavigationControls from '@/components/NavigationControls.vue'

import InboxIcon from '@/components/icons/IconInbox.vue'
import InboxFilledIcon from '@/components/icons/IconInboxFilled.vue'
import StarIcon from '@/components/icons/IconStar.vue'
import StarFilledIcon from '@/components/icons/IconStarFilled.vue'
import CheckboxIcon from '@/components/icons/IconCheckbox.vue'
import ChecboxCheckedIcon from '@/components/icons/IconCheckboxChecked.vue'
import EmailReadIcon from '@/components/icons/IconEmailRead.vue'
import EmailUnreadIcon from '@/components/icons/IconEmailUnread.vue'
import DeleteIcon from '@/components/icons/IconDelete.vue'

const router = useRouter()
const emailsStore = useEmailsStore()
const tabStore = useTabStore()

onMounted(() => {
  emailsStore.initialize()
})

const tabs: Tab[] = [
  {
    label: 'All',
    tabName: 'all',
    icon: InboxIcon,
    iconFilled: InboxFilledIcon
  },
  {
    label: 'Starred',
    tabName: 'starred',
    icon: StarIcon,
    iconFilled: StarFilledIcon
  }
]

const onTabSelect = (tabName: string) => {
  tabStore.setTab(tabName)
}

const onOpen = (emailId: string) => {
  router.push({ name: 'email', params: { id: emailId } })
  emailsStore.setReadOne(emailId)
}

const onToggleSelect = (emailId: string) => {
  emailsStore.toggleSelectedOneOfTab(emailId, tabStore.selectedTab)
}

const onToggleRead = (emailId: string) => {
  emailsStore.toggleReadOne(emailId)
}

const onToggleStar = (emailId: string) => {
  emailsStore.toggleStarredOne(emailId)
}

const onDelete = (emailId: string) => {
  if (confirm('Are you sure you want to delete this email? This action cannot be undone.')) {
    emailsStore.deleteOne(emailId)
  }
}

const onToggleSelectedTab = () => {
  emailsStore.toggleSelectedOfTab(tabStore.selectedTab)
}

const onToggleReadSelectedTab = () => {
  emailsStore.toggleReadOfSelectedOfTab(tabStore.selectedTab)
}

const onDeleteSelectedTab = () => {
  if (
    confirm('Are you sure you want to delete the selected emails? This action cannot be undone.')
  ) {
    emailsStore.deleteSelectedOfTab(tabStore.selectedTab)
  }
}

const onPreviousPageTab = () => {
  emailsStore.previousPageOfTab(tabStore.selectedTab)
}

const onNextPageTab = () => {
  emailsStore.nextPageOfTab(tabStore.selectedTab)
}
</script>

<template>
  <div class="list-view">
    <div class="list-view__header">
      <div class="list-view__actions">
        <i @click="onToggleSelectedTab()">
          <CheckboxIcon v-if="!emailsStore.getIsSelectedOfTab(tabStore.selectedTab)" />
          <ChecboxCheckedIcon v-else />
        </i>
        <i @click="onDeleteSelectedTab()">
          <DeleteIcon />
        </i>
        <i @click="onToggleReadSelectedTab()">
          <EmailReadIcon v-if="!emailsStore.getIsReadOfTab(tabStore.selectedTab)" />
          <EmailUnreadIcon v-else />
        </i>
      </div>
      <NavigationControls
        :is-previous-disabled="!emailsStore.getHasPreviousPageOfTab(tabStore.selectedTab)"
        :is-next-disabled="!emailsStore.getHasNextPageOfTab(tabStore.selectedTab)"
        :label="emailsStore.getFormattedPaginationOfTab(tabStore.selectedTab)"
        @previous="onPreviousPageTab()"
        @next="onNextPageTab()"
      />
    </div>

    <TabBar
      :tabs="tabs"
      :selectedTab="tabStore.selectedTab"
      @selectTab="(tabName: string) => onTabSelect(tabName)"
      id="tab-bar"
    />

    <div class="list-view__emails">
      <EmailItem
        v-for="email in emailsStore.getEmailsOfTab(tabStore.selectedTab)"
        :key="email.id"
        :email="email"
        @open="onOpen(email.id)"
        @toggleSelect="onToggleSelect(email.id)"
        @toggleRead="onToggleRead(email.id)"
        @toggleStar="onToggleStar(email.id)"
        @delete="onDelete(email.id)"
      />
    </div>
  </div>
</template>

<style scoped>
.list-view {
  height: 84vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding-bottom: 1rem;
  border-radius: 0.5rem;
  background-color: var(--color-background-2);
}

.list-view__header {
  flex-shrink: 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1rem 0.5rem 1rem;
}

.list-view__actions {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  gap: 0.5rem;
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

#tab-bar {
  flex-shrink: 0;
}

.list-view__emails {
  overflow-y: auto;
}

/* Scrollbar */
::-webkit-scrollbar {
  width: 0.75rem;
}
::-webkit-scrollbar-track {
  background-color: var(--color-background-2);
}
::-webkit-scrollbar-thumb {
  background: var(--color-text);
}
::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary);
}
</style>
