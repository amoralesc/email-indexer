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
    let str = this.pagination.getFormattedSettings() + '&starredOnly=' + this.starredOnly
    if (this.sortBy !== null && this.sortBy.length > 0) {
      str += '&sortBy=' + this.sortBy.join(',')
    }
    return str
  }
}

export default Settings
