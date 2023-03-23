<script setup lang="ts">
import ArrowBackIcon from '@/components/icons/IconArrowBack.vue'
import DeleteIcon from '@/components/icons/IconDelete.vue'
import EmailReadIcon from '@/components/icons/IconEmailRead.vue'
import EmailUnreadIcon from '@/components/icons/IconEmailUnread.vue'
import StarIcon from '@/components/icons/IconStar.vue'
import StarFilledIcon from '@/components/icons/IconStarFilled.vue'

import type Email from '@/models/email'
import { useEmailsStore } from '@/stores/emails'

import { useRouter } from 'vue-router'

const props = defineProps<{
  email: Email
}>()

const emailsStore = useEmailsStore()

const onBack = () => {
  useRouter().push({ name: 'home' })
}
const onDelete = () => {
  if (confirm('Are you sure you want to delete this email? This action cannot be undone.')) {
    emailsStore.deleteOne(props.email.id)
    useRouter().push({ name: 'home' })
  }
}
const onToggleRead = () => {
  emailsStore.toggleReadOne(props.email.id)
}
const onToggleStar = () => {
  emailsStore.toggleStarredOne(props.email.id)
}
</script>

<template>
  <div class="email-view">
    <div class="email-view__header">
      <div class="email-view__actions">
        <i @click="onBack()">
          <ArrowBackIcon />
        </i>
        <i @click="onDelete()">
          <DeleteIcon />
        </i>
        <i @click="onToggleRead()">
          <EmailReadIcon v-if="!email.read" />
          <EmailUnreadIcon v-else />
        </i>
        <i @click="onToggleStar()">
          <StarIcon v-if="!email.starred" />
          <StarFilledIcon v-else />
        </i>
      </div>
      <NavigationControls :is-previous-disabled="true" :is-next-disabled="false" label="1-50" />
    </div>

    <div class="email-view__content">
      <EmailItem :key="email.id" :email="email" />
    </div>
  </div>
</template>
