import Pagination from './pagination'

class Settings {
  starredOnly: boolean
  sortBy: [string] | null
  pagination: Pagination

  constructor(
    starredOnly: boolean = false,
    sortBy: [string] | null = null,
    page: number = 1,
    pageSize: number = 50,
    total: number = 0
  ) {
    this.starredOnly = starredOnly
    this.sortBy = sortBy
    this.pagination = new Pagination(page, pageSize, total)
  }

  getFormattedSettings = () => {
    let str =
      'page=' +
      this.pagination.page +
      '&pageSize=' +
      this.pagination.pageSize +
      '&starredOnly=' +
      this.starredOnly
    if (this.sortBy !== null && this.sortBy.length > 0) {
      str += '&sortBy=' + this.sortBy.join(',')
    }
    return str
  }
}

export default Settings
