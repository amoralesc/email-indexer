import axios from 'axios'

import type Email from '@/models/email'
import type Settings from '@/models/settings'
import type { Query } from '@/models/query'

const baseUrl = 'http://localhost:8080/api/emails'

const getAll = async (settings: Settings) => {
  const response = await axios.get(baseUrl + '?' + settings.getFormattedSettings())
  return response.data
}

const searchByQuery = async (query: Query, settings: Settings) => {
  const response = await axios.post(
    `${baseUrl}/search` + '?' + settings.getFormattedSettings(),
    query
  )
  return response.data
}

const searchByQueryString = async (query: string, settings: Settings) => {
  const response = await axios.get(`${baseUrl}/query?q=${query}&` + settings.getFormattedSettings())
  return response.data
}

const get = async (id: string) => {
  const response = await axios.get(`${baseUrl}/${id}`)
  return response.data
}

const getByMessageId = async (messageId: string) => {
  const response = await axios.get(`${baseUrl}/messageId/${messageId}`)
  return response.data
}

const update = async (id: string, email: Email) => {
  const response = await axios.put(`${baseUrl}/${id}`, email)
  return response.data
}

const updateMany = async (emails: Email[]) => {
  const response = await axios.put(`${baseUrl}`, emails)
  return response.data
}

const remove = async (id: string) => {
  const response = await axios.delete(`${baseUrl}/${id}`)
  return response.data
}

const removeMany = async (ids: string[]) => {
  const response = await axios.delete(`${baseUrl}/bulk/${ids.join(',')}`)
  return response.data
}

export default {
  getAll,
  searchByQuery,
  searchByQueryString,
  get,
  getByMessageId,
  update,
  updateMany,
  remove,
  removeMany
}
