import { ref } from 'vue'
import { defineStore } from 'pinia'

import type { Email } from '@/models/email'

const dummyEmail: Email = {
  id: '1',
  message_id: '1',
  from: 'John Doe',
  to: ['Jane Doe'],
  cc: [],
  bcc: [],
  subject: 'Hello',
  body: 'Hello Jane, how are you?',
  date: new Date(),
  read: false,
  starred: false,
  selected: false
}

const initialEmails: Email[] = [dummyEmail]

// generate 49 dummy emails for testing
for (let i = 2; i < 50; i++) {
  initialEmails.push({
    ...dummyEmail,
    id: i.toString(),
    message_id: i.toString(),
    starred: i % 2 === 0
  })
}

export const useEmailsStore = defineStore('emails', () => {
  const all = ref(initialEmails)
  const starred = ref(initialEmails.filter((email) => email.starred).map((email) => ({ ...email })))
  const isSelectedAll = ref(false)
  const isSelectedStarred = ref(false)
  const isReadAll = ref(false)
  const isReadStarred = ref(false)

  function toggleStarredOne(id: string) {
    all.value.concat(starred.value).forEach((email) => {
      if (email.id === id) {
        email.starred = !email.starred
      }
    })
    starred.value = all.value.filter((email) => email.starred)
  }

  function toggleReadOne(id: string) {
    all.value.concat(starred.value).forEach((email) => {
      if (email.id === id) {
        email.read = !email.read
      }
    })
  }

  function toggleSelectedOneOfAll(id: string) {
    const email = all.value.find((email) => email.id === id)
    if (email) {
      email.selected = !email.selected
    }
  }

  function toggleSelectedOneOfStarred(id: string) {
    const email = starred.value.find((email) => email.id === id)
    if (email) {
      email.selected = !email.selected
    }
  }

  function toggleSelectedAll() {
    all.value.forEach((email) => {
      email.selected = !isSelectedAll.value
    })
    isSelectedAll.value = !isSelectedAll.value
  }

  function toggleSelectedStarred() {
    starred.value.forEach((email) => {
      email.selected = !isSelectedStarred.value
    })
    isSelectedStarred.value = !isSelectedStarred.value
  }

  function toggleReadSelectedAll() {
    all.value.forEach((email) => {
      if (email.selected) {
        email.read = !isReadAll.value
      }
    })
    isReadAll.value = !isReadAll.value
  }

  function toggleReadSelectedStarred() {
    starred.value.forEach((email) => {
      if (email.selected) {
        email.read = !isReadStarred.value
      }
    })
    isReadStarred.value = !isReadStarred.value
  }

  function deleteOne(id: string) {
    const index = all.value.findIndex((email) => email.id === id)
    if (index > -1) {
      all.value.splice(index, 1)
    }
    const starredIndex = starred.value.findIndex((email) => email.id === id)
    if (starredIndex > -1) {
      starred.value.splice(starredIndex, 1)
    }
  }

  function deleteSelectedAll() {
    all.value = all.value.filter((email) => !email.selected)
    starred.value = all.value.filter((email) => email.starred)
  }

  function deleteSelectedStarred() {
    starred.value = starred.value.filter((email) => !email.selected)
    all.value = all.value.filter((email) => !email.starred)
  }

  return {
    all,
    starred,
    isSelectedAll,
    isSelectedStarred,
    isReadAll,
    isReadStarred,
    toggleStarredOne,
    toggleReadOne,
    toggleSelectedOneOfAll,
    toggleSelectedOneOfStarred,
    toggleSelectedAll,
    toggleSelectedStarred,
    toggleReadSelectedAll,
    toggleReadSelectedStarred,
    deleteOne,
    deleteSelectedAll,
    deleteSelectedStarred
  }
})
