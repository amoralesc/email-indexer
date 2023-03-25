<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

import { useEmailsStore } from '@/stores/emails'

import ArrowBackIcon from '@/components/icons/IconArrowBack.vue'
import DeleteIcon from '@/components/icons/IconDelete.vue'
import EmailReadIcon from '@/components/icons/IconEmailRead.vue'
import EmailUnreadIcon from '@/components/icons/IconEmailUnread.vue'
import StarIcon from '@/components/icons/IconStar.vue'
import StarFilledIcon from '@/components/icons/IconStarFilled.vue'

const router = useRouter()
const route = useRoute()
const emailsStore = useEmailsStore()
const emailId = route.params.id as string
const email = emailsStore.getEmailById(emailId)

const onBack = () => {
  // navigate back in history
  router.go(-1)
}
const onToggleRead = () => {
  emailsStore.toggleReadOne(emailId)
}
const onToggleStar = () => {
  emailsStore.toggleStarredOne(emailId)
}
const onDelete = () => {
  if (confirm('Are you sure you want to delete this email? This action cannot be undone.')) {
    emailsStore.deleteOne(emailId)
    router.push({ name: 'home' })
  }
}
</script>

<template>
  <div class="email-view" v-if="email">
    <div class="email-view__header">
      <div class="email-view__actions">
        <i @click="onBack()">
          <ArrowBackIcon />
        </i>
        <i @click="onDelete()">
          <DeleteIcon />
        </i>
        <i @click="onToggleRead()">
          <EmailReadIcon v-if="!email.isRead" />
          <EmailUnreadIcon v-else />
        </i>
        <i @click="onToggleStar()">
          <StarIcon v-if="!email.isStarred" />
          <StarFilledIcon v-else />
        </i>
      </div>
      <NavigationControls :is-previous-disabled="true" :is-next-disabled="false" label="1-50" />
    </div>

    <div class="email-view__content">
      <div class="email-view__content__headers">
        <h2>{{ email.subject }}</h2>
        <p>Date: {{ email.getFormattedDate() }}</p>
        <p>From: {{ email.from }}</p>
        <p>To: {{ email.to.join(', ') }}</p>
        <p>CC: {{ email.cc.join(', ') }}</p>
        <p>BCC: {{ email.bcc.join(', ') }}</p>
        <p></p>
      </div>
      <div class="email-view__content__body">
        <p>{{ email.body }}</p>
      </div>
    </div>
  </div>
  <div v-else>
    <h1>404</h1>
    <p>Sorry, we couldn't find the email you were looking for.</p>
  </div>
</template>

<style scoped>
.email-view {
  height: 84vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding-bottom: 1rem;
  border-radius: 0.5rem;
  background-color: var(--color-background-2);
}

.email-view__header {
  flex-shrink: 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1rem 0.5rem 1rem;
}

.email-view__actions {
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

.email-view__content {
  display: flex;
  flex-direction: column;
  overflow: auto;
  padding: 0rem 1.5rem;
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
