<script setup lang="ts">
import EmailItem from '@/components/EmailItem.vue'
import TabBar from '@/components/TabBar.vue'
import NavigationControls from '@/components/NavigationControls.vue'

import CheckboxIcon from '@/components/icons/IconCheckbox.vue'
import ChecboxCheckedIcon from '@/components/icons/IconCheckboxChecked.vue'
import InboxIcon from '@/components/icons/IconInbox.vue'
import InboxFilledIcon from '@/components/icons/IconInboxFilled.vue'
import StarIcon from '@/components/icons/IconStar.vue'
import StarFilledIcon from '@/components/icons/IconStarFilled.vue'
import EmailReadIcon from '@/components/icons/IconEmailRead.vue'
import EmailUnreadIcon from '@/components/icons/IconEmailUnread.vue'
import DeleteIcon from '@/components/icons/IconDelete.vue'

import type { Tab } from '@/models/tab'

import { useEmailsStore } from '@/stores/emails'
import { useTabStore } from '@/stores/tab'

const emailsStore = useEmailsStore()
const tabStore = useTabStore()

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
const onToggleSelect = (emailId: string) => {
  if (tabStore.selectedTab === 'all') {
    emailsStore.toggleSelectedOneOfAll(emailId)
  } else {
    emailsStore.toggleSelectedOneOfStarred(emailId)
  }
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
const onToggleSelected = () => {
  if (tabStore.selectedTab === 'all') {
    emailsStore.toggleSelectedAll()
  } else {
    emailsStore.toggleSelectedStarred()
  }
}
const onToggleReadSelected = () => {
  if (tabStore.selectedTab === 'all') {
    emailsStore.toggleReadSelectedAll()
  } else {
    emailsStore.toggleReadSelectedStarred()
  }
}
const onDeleteSelected = () => {
  if (
    confirm('Are you sure you want to delete the selected emails? This action cannot be undone.')
  ) {
    if (tabStore.selectedTab === 'all') {
      emailsStore.deleteSelectedAll()
    } else {
      emailsStore.deleteSelectedStarred()
    }
  }
}

const getTabEmails = () => {
  if (tabStore.selectedTab === 'all') {
    return emailsStore.all
  } else {
    return emailsStore.starred
  }
}
const getTabIsSelected = () => {
  if (tabStore.selectedTab === 'all') {
    return emailsStore.isSelectedAll
  } else {
    return emailsStore.isSelectedStarred
  }
}
const getTabIsRead = () => {
  if (tabStore.selectedTab === 'all') {
    return emailsStore.isReadAll
  } else {
    return emailsStore.isReadStarred
  }
}
</script>

<template>
  <div class="list-view">
    <div class="list-view__header">
      <div class="list-view__actions">
        <i class="list-view__checkbox" @click="onToggleSelected()">
          <CheckboxIcon v-if="!getTabIsSelected()" />
          <ChecboxCheckedIcon v-else />
        </i>
        <i class="list-view__delete" @click="onDeleteSelected()">
          <DeleteIcon />
        </i>
        <i class="list-view__read" @click="onToggleReadSelected()">
          <EmailReadIcon v-if="!getTabIsRead()" />
          <EmailUnreadIcon v-else />
        </i>
      </div>
      <NavigationControls :is-previous-disabled="false" :is-next-disabled="false" label="1-50" />
    </div>

    <TabBar
      :tabs="tabs"
      :selectedTab="tabStore.selectedTab"
      @selectTab="(tabName: string) => onTabSelect(tabName)"
      id="tab-bar"
    />

    <div class="list-view__emails">
      <EmailItem
        v-for="email in getTabEmails()"
        :key="email.id"
        :email="email"
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
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding-bottom: 5rem;
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
  background: var(--color-background);
}
::-webkit-scrollbar-thumb {
  background: var(--color-text);
}
::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary);
}
</style>
