class Pagination {
  page: number
  pageSize: number
  total: number

  constructor(page: number = 1, pageSize: number = 50, total: number = 0) {
    this.page = page
    this.pageSize = pageSize
    this.total = total
  }

  getFormattedPagination = () => {
    const lowerBound = (this.page - 1) * this.pageSize + 1
    const upperBound =
      this.page * this.pageSize > this.total ? this.total : this.page * this.pageSize

    return lowerBound + '-' + upperBound + ' of ' + this.total
  }
}

export default Pagination
