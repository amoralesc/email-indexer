import { ref } from 'vue'
import { defineStore } from 'pinia'

import Email from '@/models/email'
import Settings from '@/models/settings'

import service from '@/services'

interface EmailState {
  emails: {
    current: Email[]
    next: Email[]
  }
  isSelected: boolean
  isRead: boolean
  settings: Settings
}

const waitTime = 1000

export const useEmailsStore = defineStore('emails', () => {
  const all = ref({
    emails: {
      current: [],
      next: []
    },
    isSelected: false,
    isRead: false,
    settings: new Settings()
  } as EmailState)

  const starred = ref({
    emails: {
      previous: [],
      current: [],
      next: []
    },
    isSelected: false,
    isRead: false,
    settings: new Settings(true)
  } as EmailState)

  // getters

  const getEmailsOfTab = (tab: string): Email[] => {
    if (tab === 'all') {
      return all.value.emails.current
    }
    return starred.value.emails.current
  }

  const getFormattedPaginationOfTab = (tab: string): string => {
    if (tab === 'all') {
      return all.value.settings.pagination.getFormattedPagination()
    }
    return starred.value.settings.pagination.getFormattedPagination()
  }

  const getIsSelectedOfTab = (tab: string): boolean => {
    if (tab === 'all') {
      return all.value.isSelected
    }
    return starred.value.isSelected
  }

  const getIsReadOfTab = (tab: string): boolean => {
    if (tab === 'all') {
      return all.value.isRead
    }
    return starred.value.isRead
  }

  const getEmailById = (emailId: string): Email | undefined => {
    const email = all.value.emails.current.find((email) => email.id === emailId)
    if (email) {
      return email
    }
    return starred.value.emails.current.find((email) => email.id === emailId)
  }

  // actions

  async function fetchEmailsOfAll() {
    const response = await service.getAll(all.value.settings)

    // the above response returns double the amount of emails
    // the first half is the current page, the second half is the next page
    // so we need to split the response into two arrays
    const emails: Email[] = []
    for (const email of response.emails) {
      emails.push(Email.fromJSON(email))
    }

    all.value.emails.current = emails.slice(0, all.value.settings.pagination.pageSize)
    all.value.emails.next = emails.slice(all.value.settings.pagination.pageSize)
    all.value.settings.pagination.total = response.total
  }

  async function fetchEmailsOfStarred() {
    const response = await service.getAll(starred.value.settings)

    const emails: Email[] = []
    for (const email of response.emails) {
      emails.push(Email.fromJSON(email))
    }

    starred.value.emails.current = emails.slice(0, starred.value.settings.pagination.pageSize)
    starred.value.emails.next = emails.slice(starred.value.settings.pagination.pageSize)
    starred.value.settings.pagination.total = response.total
  }

  async function initialize() {
    await fetchEmailsOfAll()
    await fetchEmailsOfStarred()
  }

  function toggleSelectedOfAll() {
    all.value.isSelected = !all.value.isSelected
    all.value.emails.current.forEach((email) => {
      email.isSelected = all.value.isSelected
    })
  }

  function toggleSelectedOfStarred() {
    starred.value.isSelected = !starred.value.isSelected
    starred.value.emails.current.forEach((email) => {
      email.isSelected = starred.value.isSelected
    })
  }

  function toggleSelectedOfTab(tab: string) {
    if (tab === 'all') {
      toggleSelectedOfAll()
    } else {
      toggleSelectedOfStarred()
    }
  }

  async function toggleReadOfSelectedOfAll() {
    all.value.isRead = !all.value.isRead
    all.value.emails.current.forEach((email) => {
      if (email.isSelected) {
        email.isRead = all.value.isRead

        const emailOfStarred = starred.value.emails.current.find(
          (emailOfStarred) => emailOfStarred.id === email.id
        )
        if (emailOfStarred) {
          emailOfStarred.isRead = all.value.isRead
        }
      }
    })

    service.updateMany(all.value.emails.current.filter((email) => email.isSelected))
  }

  async function toggleReadOfSelectedOfStarred() {
    starred.value.isRead = !starred.value.isRead
    starred.value.emails.current.forEach((email) => {
      if (email.isSelected) {
        email.isRead = starred.value.isRead

        const emailOfAll = all.value.emails.current.find((emailOfAll) => emailOfAll.id === email.id)
        if (emailOfAll) {
          emailOfAll.isRead = starred.value.isRead
        }
      }
    })

    service.updateMany(starred.value.emails.current.filter((email) => email.isSelected))
  }

  async function toggleReadOfSelectedOfTab(tab: string) {
    if (tab === 'all') {
      await toggleReadOfSelectedOfAll()
    } else {
      await toggleReadOfSelectedOfStarred()
    }
  }

  async function deleteSelectedOfAll() {
    // filter out the emails that are selected
    all.value.emails.current = all.value.emails.current.filter((email) => !email.isSelected)
    // fill current page with next page top emails until current page is full
    while (all.value.emails.current.length < all.value.settings.pagination.pageSize) {
      const nextEmail = all.value.emails.next.shift()
      if (nextEmail) {
        all.value.emails.current.push(nextEmail)
      } else {
        break
      }
    }

    // now call the services
    // this is done in this order because the service is slow
    // to register the delete changes and needs to be waited on
    // so the re-fetching of the emails works properly
    const ids = all.value.emails.current
      .filter((email) => email.isSelected)
      .map((email) => email.id)

    await service.removeMany(ids)
    await new Promise((resolve) => setTimeout(resolve, waitTime))
    await fetchEmailsOfAll()
    await fetchEmailsOfStarred()
  }

  async function deleteSelectedOfStarred() {
    starred.value.emails.current = starred.value.emails.current.filter((email) => !email.isSelected)
    while (starred.value.emails.current.length < starred.value.settings.pagination.pageSize) {
      const nextEmail = starred.value.emails.next.shift()
      if (nextEmail) {
        starred.value.emails.current.push(nextEmail)
      } else {
        break
      }
    }

    const ids = starred.value.emails.current
      .filter((email) => email.isSelected)
      .map((email) => email.id)

    await service.removeMany(ids)
    await new Promise((resolve) => setTimeout(resolve, waitTime))
    await fetchEmailsOfAll()
    await fetchEmailsOfStarred()
  }

  async function deleteSelectedOfTab(tab: string) {
    if (tab === 'all') {
      await deleteSelectedOfAll()
    } else {
      await deleteSelectedOfStarred()
    }
  }

  function toggleSelectedOneOfAll(emailId: string) {
    const email = all.value.emails.current.find((email) => email.id === emailId)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  function toggleSelectedOneOfStarred(emailId: string) {
    const email = starred.value.emails.current.find((email) => email.id === emailId)
    if (email) {
      email.isSelected = !email.isSelected
    }
  }

  function toggleSelectedOneOfTab(tab: string, emailId: string) {
    if (tab === 'all') {
      toggleSelectedOneOfAll(emailId)
    } else {
      toggleSelectedOneOfStarred(emailId)
    }
  }

  async function toggleReadOne(emailId: string) {
    const emailOfAll = all.value.emails.current.find((email) => email.id === emailId)
    const emailOfStarred = starred.value.emails.current.find((email) => email.id === emailId)

    if (emailOfAll) {
      emailOfAll.isRead = !emailOfAll.isRead
      if (emailOfStarred) {
        emailOfStarred.isRead = emailOfAll.isRead
      }

      await service.update(emailId, emailOfAll)
    } else {
      if (emailOfStarred) {
        emailOfStarred.isRead = !emailOfStarred.isRead

        await service.update(emailId, emailOfStarred)
      }
    }
  }

  async function toggleStarredOne(emailId: string) {
    const emailOfAll = all.value.emails.current.find((email) => email.id === emailId)
    const emailOfStarred = starred.value.emails.current.find((email) => email.id === emailId)

    if (emailOfAll) {
      emailOfAll.isStarred = !emailOfAll.isStarred
      if (emailOfStarred) {
        emailOfStarred.isStarred = emailOfAll.isStarred
      }

      await service.update(emailId, emailOfAll)
      await new Promise((resolve) => setTimeout(resolve, waitTime))
      await fetchEmailsOfStarred()
    } else {
      if (emailOfStarred) {
        // if the email is only in the starred list, it means we are in the starred tab
        // so we need to filter it out of the starred list and add the next email
        // from the next page to the current page
        starred.value.emails.current = starred.value.emails.current.filter(
          (email) => email.id !== emailId
        )
        const nextEmail = starred.value.emails.next.shift()
        if (nextEmail) {
          starred.value.emails.current.push(nextEmail)
        }

        await service.update(emailId, emailOfStarred)
        await new Promise((resolve) => setTimeout(resolve, waitTime))
        await fetchEmailsOfStarred()
      }
    }
  }

  async function deleteOne(emailId: string) {
    // filter out the email from the current page
    all.value.emails.current = all.value.emails.current.filter((email) => email.id !== emailId)
    // fill current page with next page top emails until current page is full
    // ONLY if current page has less than the page size
    if (all.value.emails.current.length < all.value.settings.pagination.pageSize) {
      const nextEmail = all.value.emails.next.shift()
      if (nextEmail) {
        all.value.emails.current.push(nextEmail)
      }
    }

    // same but for the starred emails
    starred.value.emails.current = starred.value.emails.current.filter(
      (email) => email.id !== emailId
    )
    if (starred.value.emails.current.length < starred.value.settings.pagination.pageSize) {
      const nextEmail = starred.value.emails.next.shift()
      if (nextEmail) {
        starred.value.emails.current.push(nextEmail)
      }
    }

    // now call the services
    await service.remove(emailId)
    await new Promise((resolve) => setTimeout(resolve, waitTime))
    await fetchEmailsOfAll()
    await fetchEmailsOfStarred()
  }

  async function setReadOne(emailId: string) {
    const emailOfAll = all.value.emails.current.find((email) => email.id === emailId)
    const emailOfStarred = all.value.emails.current.find((email) => email.id === emailId)

    if (emailOfAll) {
      emailOfAll.isRead = true
      if (emailOfStarred) {
        emailOfStarred.isRead = true
      }
      await service.update(emailId, emailOfAll)
    } else {
      if (emailOfStarred) {
        emailOfStarred.isRead = true
        await service.update(emailId, emailOfStarred)
      }
    }
  }

  return {
    all,
    starred,
    // export all getters
    getEmailsOfTab,
    getFormattedPaginationOfTab,
    getIsSelectedOfTab,
    getIsReadOfTab,
    getEmailById,
    // export all actions
    fetchEmailsOfAll,
    fetchEmailsOfStarred,
    initialize,
    toggleSelectedOfAll,
    toggleSelectedOfStarred,
    toggleSelectedOfTab,
    toggleReadOfSelectedOfAll,
    toggleReadOfSelectedOfStarred,
    toggleReadOfSelectedOfTab,
    deleteSelectedOfAll,
    deleteSelectedOfStarred,
    deleteSelectedOfTab,
    toggleSelectedOneOfAll,
    toggleSelectedOneOfStarred,
    toggleSelectedOneOfTab,
    toggleReadOne,
    toggleStarredOne,
    deleteOne,
    setReadOne
  }
})
