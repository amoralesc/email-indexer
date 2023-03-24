import { ref } from 'vue'
import { defineStore } from 'pinia'

import Email from '@/models/email'

const dummyEmail = new Email(
  '1',
  '1',
  new Date(),
  'John Doe',
  ['Jane Doe'],
  ['Jane Doe'],
  ['Jane Doe'],
  'Lorem ipsum dolor sit amet',
  'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel ultricies lacinia, nisl nisl aliquam nisl.',
  false,
  false
)

const initialEmails: Email[] = [dummyEmail]
for (let i = 2; i < 50; i++) {
  initialEmails.push({
    ...dummyEmail,
    id: i.toString(),
    messageId: i.toString(),
    isStarred: i % 2 === 0
  })
}

export const useEmailsStore = defineStore('emails', () => {
  const all = ref(initialEmails)
  const starred = ref(
    initialEmails.filter((email) => email.isStarred).map((email) => ({ ...email }))
  )
  const isSelectedAll = ref(false)
  const isSelectedStarred = ref(false)
  const isReadAll = ref(false)
  const isReadStarred = ref(false)

  const getEmailById = (id: string) => {
    return all.value.find((email) => email.id === id)
  }

  function toggleStarredOne(id: string) {
    all.value.concat(starred.value).forEach((email) => {
      if (email.id === id) {
        email.isStarred = !email.isStarred
      }
    })
    starred.value = all.value.filter((email) => email.isStarred)
  }

  function toggleReadOne(id: string) {
    all.value.forEach((email) => {
      if (email.id === id) {
        email.isRead = !email.isRead
      }
    })
    starred.value.forEach((email) => {
      if (email.id === id) {
        email.isRead = !email.isRead
      }
    })
  }

  function toggleSelectedOneOfAll(id: string) {
    const email = all.value.find((email) => email.id === id)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  function toggleSelectedOneOfStarred(id: string) {
    const email = starred.value.find((email) => email.id === id)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  function toggleSelectedAll() {
    all.value.forEach((email) => {
      email.isSelected = !isSelectedAll.value
    })
    isSelectedAll.value = !isSelectedAll.value
  }

  function toggleSelectedStarred() {
    starred.value.forEach((email) => {
      email.isSelected = !isSelectedStarred.value
    })
    isSelectedStarred.value = !isSelectedStarred.value
  }

  function toggleReadSelectedAll() {
    all.value.forEach((email) => {
      if (email.isSelected) {
        email.isRead = !isReadAll.value
      }
    })
    isReadAll.value = !isReadAll.value
  }

  function toggleReadSelectedStarred() {
    starred.value.forEach((email) => {
      if (email.isSelected) {
        email.isRead = !isReadStarred.value
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
    all.value = all.value.filter((email) => !email.isSelected)
    starred.value = all.value.filter((email) => email.isStarred)
  }

  function deleteSelectedStarred() {
    starred.value = starred.value.filter((email) => !email.isSelected)
    all.value = all.value.filter((email) => !email.isStarred)
  }

  return {
    all,
    starred,
    isSelectedAll,
    isSelectedStarred,
    isReadAll,
    isReadStarred,
    getEmailById,
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
