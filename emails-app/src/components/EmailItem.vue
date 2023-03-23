<script setup lang="ts">
import type Email from '@/models/email'

import CheckboxIcon from './icons/IconCheckbox.vue'
import ChecboxCheckedIcon from './icons/IconCheckboxChecked.vue'
import StarIcon from './icons/IconStar.vue'
import StarFilledIcon from './icons/IconStarFilled.vue'
import DeleteIcon from './icons/IconDelete.vue'
import EmailReadIcon from './icons/IconEmailRead.vue'
import EmailUnreadIcon from './icons/IconEmailUnread.vue'

defineProps<{
  email: Email
}>()

defineEmits<{
  (event: 'open'): void
  (event: 'toggleSelect'): void
  (event: 'toggleStar'): void
  (event: 'delete'): void
  (event: 'toggleRead'): void
}>()
</script>

<template>
  <div
    class="email-item"
    :class="{
      'email-item--selected': email.selected,
      'email-item--starred': email.starred,
      'email-item--read': email.read,
      'email-item--unread': !email.read
    }"
    @click="$emit('open')"
  >
    <div class="email-item__actions">
      <i class="email-item__checkbox" @click="$emit('toggleSelect')" v-on:click.stop>
        <CheckboxIcon v-if="!email.selected" />
        <ChecboxCheckedIcon v-else />
      </i>
      <i class="email-item__star" @click="$emit('toggleStar')" v-on:click.stop>
        <StarIcon v-if="!email.starred" />
        <StarFilledIcon v-else />
      </i>
    </div>

    <div class="email-item__subject">{{ email.subject }}</div>
    <div class="email-item__from">{{ email.from }}</div>
    <div class="email-item__to">{{ email.to.join(', ') }}</div>
    <div class="email-item__date">{{ email.getFormattedDate() }}</div>

    <div class="email-item__hover_actions">
      <i class="email-item__delete" @click="$emit('delete')" v-on:click.stop>
        <DeleteIcon />
      </i>
      <i class="email-item__read" @click="$emit('toggleRead')" v-on:click.stop>
        <EmailUnreadIcon v-if="email.read" />
        <EmailReadIcon v-else />
      </i>
    </div>
  </div>
</template>

<style scoped>
.email-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
}

.email-item:hover {
  box-shadow: 0 0 0 1px var(--color-border-hover);
}

.email-item:hover i {
  color: var(--color-heading);
}

.email-item--unread {
  color: var(--color-heading);
}

.email-item--selected {
  background-color: var(--color-primary-soft);
}

.email-item__actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.email-item__star {
  width: 1.5rem;
  height: 1.5rem;
  margin-top: -0.1rem;
}

.email-item--starred .email-item__star {
  color: var(--color-primary);
}

i {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  cursor: pointer;
}

i:hover {
  color: var(--color-primary);
}

.email-item:hover i:hover {
  color: var(--color-primary);
}

.email-item--unread i {
  color: var(--color-text);
}

.email-item__subject {
  flex: 2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-item__from,
.email-item__to {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-item__date {
  flex-shrink: 0;
  font-family: 'IBM Plex Mono', monospace;
  font-size: 0.8rem;
  text-align: right;
}

.email-item__hover_actions {
  position: absolute;
  right: 1rem;
  z-index: 999;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  opacity: 0;
}

.email-item:hover .email-item__hover_actions {
  opacity: 1;
}

.email-item:hover .email-item__date {
  opacity: 0;
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
