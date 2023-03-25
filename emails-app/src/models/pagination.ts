class Pagination {
  page: number
  pageSize: number
  total: number

  constructor(page: number = 1, pageSize: number = 50, total: number = 0) {
    this.page = page
    this.pageSize = pageSize
    this.total = total
  }

  // pagination needs to translate the page number and page size
  // into start and size parameters for the API
  // also, the API calls need pagination to send twice the size
  // because 2 pages are loaded at a time
  getFormattedSettings = (): string => {
    const start = (this.page - 1) * this.pageSize
    const size = this.pageSize * 2

    return 'start=' + start + '&size=' + size
  }

  getFormattedPagination = (): string => {
    const lowerBound = (this.page - 1) * this.pageSize + 1
    const upperBound =
      this.page * this.pageSize > this.total ? this.total : this.page * this.pageSize

    return lowerBound + '-' + upperBound + ' of ' + this.total
  }

  getMaxPage = () => {
    return Math.ceil(this.total / this.pageSize)
  }
}

export default Pagination
