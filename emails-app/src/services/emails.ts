import axios from 'axios'

import type Email from '@/models/email'
import type Settings from '@/models/settings'
import type { Response } from '@/models/response'
import type { Query } from '@/models/query'

const baseUrl = 'http://localhost:8080/api/emails'

const getAll = async (settings: Settings): Promise<Response> => {
  const response = await axios.get(baseUrl + '?' + settings.getFormattedSettings())
  return response.data
}

const searchByQuery = async (query: Query, settings: Settings): Promise<Response> => {
  const response = await axios.post(
    `${baseUrl}/search` + '?' + settings.getFormattedSettings(),
    query
  )
  return response.data
}

const searchByQueryString = async (query: string, settings: Settings): Promise<Response> => {
  const response = await axios.get(`${baseUrl}/query?q=${query}&` + settings.getFormattedSettings())
  return response.data
}

const get = async (id: string): Promise<Email> => {
  const response = await axios.get(`${baseUrl}/${id}`)
  return response.data
}

const getByMessageId = async (messageId: string): Promise<Email> => {
  const response = await axios.get(`${baseUrl}/messageId/${messageId}`)
  return response.data
}

const update = async (id: string, email: Email): Promise<Email> => {
  const response = await axios.put(`${baseUrl}/${id}`, email)
  return response.data
}

const updateMany = async (emails: Email[]): Promise<Email[]> => {
  const response = await axios.put(`${baseUrl}`, emails)
  return response.data
}

const remove = async (id: string): Promise<void> => {
  await axios.delete(`${baseUrl}/${id}`)
}

const removeMany = async (ids: string[]): Promise<void> => {
  await axios.delete(baseUrl + '?ids=' + ids.join(','))
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
