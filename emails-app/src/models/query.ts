interface DateRange {
  from?: Date
  to?: Date
}

interface Query {
  from?: string
  to?: string[]
  cc?: string[]
  bcc?: string[]
  subjectIncludes?: string
  bodyIncludes?: string
  bodyExcludes?: string
  dateRange?: DateRange
}

export type { DateRange, Query }
