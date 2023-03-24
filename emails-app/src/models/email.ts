class Email {
  id: string
  messageId: string
  date: Date
  from: string
  to: string[]
  cc: string[]
  bcc: string[]
  subject: string
  body: string
  isRead: boolean
  isStarred: boolean
  isSelected: boolean

  constructor(
    id: string,
    messageId: string,
    date: Date,
    from: string,
    to: string[],
    cc: string[],
    bcc: string[],
    subject: string,
    body: string,
    isRead: boolean,
    isStarred: boolean,
    isSelected: boolean = false
  ) {
    this.id = id
    this.messageId = messageId
    this.date = date
    this.from = from
    this.to = to
    this.cc = cc
    this.bcc = bcc
    this.subject = subject
    this.body = body
    this.isRead = isRead
    this.isStarred = isStarred
    this.isSelected = isSelected
  }

  getFormattedDate = () => {
    const day = this.date.getDate()
    const month = this.date.getMonth() + 1
    const year = this.date.getFullYear()

    const formattedDay = day < 10 ? `0${day}` : day
    const formattedMonth = month < 10 ? `0${month}` : month

    return `${formattedDay}/${formattedMonth}/${year}`
  }
}

export default Email
