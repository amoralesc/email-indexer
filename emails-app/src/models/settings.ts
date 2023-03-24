class Settings {
  page: number
  pageSize: number
  sortBy: [string] | null
  starredOnly: boolean

  constructor(
    page: number = 1,
    pageSize: number = 50,
    sortBy: [string] | null = null,
    starredOnly: boolean = false
  ) {
    this.page = page
    this.pageSize = pageSize
    this.sortBy = sortBy
    this.starredOnly = starredOnly
  }

  getFormattedSettings = () => {
    let str =
      'page=' + this.page + '&pageSize=' + this.pageSize + '&starredOnly=' + this.starredOnly
    if (this.sortBy !== null && this.sortBy.length > 0) {
      str += '&sortBy=' + this.sortBy.join(',')
    }
    return str
  }
}

export default Settings
