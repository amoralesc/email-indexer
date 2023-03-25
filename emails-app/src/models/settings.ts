import Pagination from './pagination'

class Settings {
  starredOnly: boolean
  sortBy: [string] | null
  pagination: Pagination

  constructor(
    starredOnly: boolean = false,
    sortBy: [string] | null = null,
    pagination = new Pagination()
  ) {
    this.starredOnly = starredOnly
    this.sortBy = sortBy
    this.pagination = pagination
  }

  getFormattedSettings = (): string => {
    let str =
      'page=' +
      this.pagination.page +
      '&pageSize=' +
      this.pagination.pageSize * 2 +
      '&starredOnly=' +
      this.starredOnly
    if (this.sortBy !== null && this.sortBy.length > 0) {
      str += '&sortBy=' + this.sortBy.join(',')
    }
    return str
  }
}

export default Settings
