import type Email from './email'

interface Response {
  total: number
  took: number
  emails: Email[]
}

export type { Response }
