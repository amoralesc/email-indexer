import { ref } from 'vue'
import { defineStore } from 'pinia'

import Email from '@/models/email'
import type { Response } from '@/models/response'
import Settings from '@/models/settings'

import service from '@/services/emails'

interface EmailStore {
  all: {
    emails: Email[]
    isRead: boolean
    isSelected: boolean
    settings: Settings
  }
  starred: {
    emails: Email[]
    isRead: boolean
    isSelected: boolean
    settings: Settings
  }
}

export const useEmailsStore = defineStore('emails', () => {
  const data = ref<EmailStore>({
    all: {
      emails: [],
      isSelected: false,
      isRead: false,
      settings: new Settings()
    },
    starred: {
      emails: [],
      isSelected: false,
      isRead: false,
      settings: new Settings(true)
    }
  })

  // getters
  const getAllEmails = () => {
    return data.value.all.emails
  }
  const getStarredEmails = () => {
    return data.value.starred.emails
  }
  const getAllIsSelected = () => {
    return data.value.all.isSelected
  }
  const getStarredIsSelected = () => {
    return data.value.starred.isSelected
  }
  const getAllIsRead = () => {
    return data.value.all.isRead
  }
  const getStarredIsRead = () => {
    return data.value.starred.isRead
  }
  const getAllFormattedPagination = () => {
    return data.value.all.settings.pagination.getFormattedPagination()
  }
  const getStarredFormattedPagination = () => {
    return data.value.starred.settings.pagination.getFormattedPagination()
  }
  const getEmailById = (emailId: string) => {
    const email = data.value.all.emails.find((email) => email.id === emailId)
    if (email) {
      return email
    }
    return data.value.starred.emails.find((email) => email.id === emailId)
  }

  // actions
  async function fetchAllEmails() {
    const response = (await service.getAll(data.value.all.settings)) as Response

    const emails: Email[] = []
    for (const email of response.emails) {
      emails.push(Email.fromJSON(email))
    }

    data.value.all.emails = emails
    data.value.all.settings.pagination.total = response.total

    console.log('fetchAllEmails', data.value.all.emails)
  }

  async function fetchStarredEmails() {
    const response = (await service.getAll(data.value.starred.settings)) as Response

    const emails: Email[] = []
    for (const email of response.emails) {
      emails.push(Email.fromJSON(email))
    }

    data.value.starred.emails = emails
    data.value.starred.settings.pagination.total = response.total
  }

  async function initializeData() {
    fetchAllEmails()
    fetchStarredEmails()
  }

  function toggleSelectedOfAll() {
    data.value.all.isSelected = !data.value.all.isSelected
    data.value.all.emails.forEach((email) => {
      email.isSelected = data.value.all.isSelected
    })
  }

  function toggleSelectedOfStarred() {
    data.value.starred.isSelected = !data.value.starred.isSelected
    data.value.starred.emails.forEach((email) => {
      email.isSelected = data.value.starred.isSelected
    })
  }

  async function toggleReadOfAllSelected() {
    data.value.all.isRead = !data.value.all.isRead
    data.value.all.emails.forEach((email) => {
      if (email.isSelected) {
        email.isRead = data.value.all.isRead

        const starredEmail = data.value.starred.emails.find(
          (starredEmail) => starredEmail.id === email.id
        )
        if (starredEmail) {
          starredEmail.isRead = data.value.all.isRead
        }
      }
    })

    service.updateMany(data.value.all.emails.filter((email) => email.isSelected))
  }

  async function toggleReadOfStarredSelected() {
    data.value.starred.isRead = !data.value.starred.isRead
    data.value.starred.emails.forEach((email) => {
      if (email.isSelected) {
        email.isRead = data.value.starred.isRead

        const allEmail = data.value.all.emails.find((allEmail) => allEmail.id === email.id)
        if (allEmail) {
          allEmail.isRead = data.value.starred.isRead
        }
      }
    })

    service.updateMany(data.value.starred.emails.filter((email) => email.isSelected))
  }

  async function deleteSelectedOfAll() {
    const selectedEmails = data.value.all.emails.filter((email) => email.isSelected)
    await service.removeMany(selectedEmails.map((email) => email.id))
    await fetchAllEmails()
    await fetchStarredEmails()
  }

  async function deleteSelectedOfStarred() {
    const selectedEmails = data.value.starred.emails.filter((email) => email.isSelected)
    await service.removeMany(selectedEmails.map((email) => email.id))
    await fetchAllEmails()
    await fetchStarredEmails()
  }

  function toggleSelectedOneOfAll(emailId: string) {
    const email = data.value.all.emails.find((email) => email.id === emailId)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  function toggleSelectedOneOfStarred(emailId: string) {
    const email = data.value.starred.emails.find((email) => email.id === emailId)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  async function toggleReadOne(emailId: string) {
    const allEmail = data.value.all.emails.find((email) => email.id === emailId)
    const starredEmail = data.value.starred.emails.find((email) => email.id === emailId)

    if (allEmail) {
      allEmail.isRead = !allEmail.isRead
      if (starredEmail) {
        starredEmail.isRead = allEmail.isRead
      }
      await service.update(emailId, allEmail)
    } else {
      if (starredEmail) {
        starredEmail.isRead = !starredEmail.isRead
        await service.update(emailId, starredEmail)
      }
    }
  }

  async function toggleStarredOne(emailId: string) {
    const allEmail = data.value.all.emails.find((email) => email.id === emailId)
    const starredEmail = data.value.starred.emails.find((email) => email.id === emailId)

    if (allEmail) {
      allEmail.isStarred = !allEmail.isStarred
      if (starredEmail) {
        starredEmail.isStarred = allEmail.isStarred
      }
      await service.update(emailId, allEmail)
      await fetchStarredEmails()
    } else {
      if (starredEmail) {
        starredEmail.isStarred = !starredEmail.isStarred
        await service.update(emailId, starredEmail)
        await fetchStarredEmails()
      }
    }
  }

  async function deleteOne(emailId: string) {
    await service.remove(emailId)
    await fetchAllEmails()
    await fetchStarredEmails()
  }

  async function setReadOne(emailId: string) {
    const allEmail = data.value.all.emails.find((email) => email.id === emailId)
    const starredEmail = data.value.starred.emails.find((email) => email.id === emailId)

    if (allEmail) {
      allEmail.isRead = true
      if (starredEmail) {
        starredEmail.isRead = true
      }
      await service.update(emailId, allEmail)
    } else {
      if (starredEmail) {
        starredEmail.isRead = true
        await service.update(emailId, starredEmail)
      }
    }
  }

  return {
    data,
    getAllEmails,
    getStarredEmails,
    getAllIsSelected,
    getStarredIsSelected,
    getAllIsRead,
    getStarredIsRead,
    getAllFormattedPagination,
    getStarredFormattedPagination,
    getEmailById,
    fetchAllEmails,
    fetchStarredEmails,
    initializeData,
    toggleSelectedOfAll,
    toggleSelectedOfStarred,
    toggleReadOfAllSelected,
    toggleReadOfStarredSelected,
    deleteSelectedOfAll,
    deleteSelectedOfStarred,
    toggleSelectedOneOfAll,
    toggleSelectedOneOfStarred,
    toggleReadOne,
    toggleStarredOne,
    deleteOne,
    setReadOne
  }
})
