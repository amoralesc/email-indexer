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

  static fromJSON = (json: any) => {
    if (!json.messageId) json.messageId = ''
    if (!json.date) json.date = new Date()
    if (!json.from) json.from = ''
    if (!json.to) json.to = []
    if (!json.cc) json.cc = []
    if (!json.bcc) json.bcc = []
    if (!json.subject) json.subject = ''
    if (!json.body) json.body = ''
    if (!json.isRead) json.isRead = false
    if (!json.isStarred) json.isStarred = false
    if (!json.isSelected) json.isSelected = false

    return new Email(
      json._id,
      json.messageId,
      new Date(json.date),
      json.from,
      json.to,
      json.cc,
      json.bcc,
      json.subject,
      json.body,
      json.isRead,
      json.isStarred,
      json.isSelected
    )
  }

  getFormattedDate = () => {
    const day = this.date.getDate()
    const month = this.date.getMonth() + 1
    const year = this.date.getFullYear()

    const formattedDay = day < 10 ? `0${day}` : day
    const formattedMonth = month < 10 ? `0${month}` : month

    return `${formattedDay}/${formattedMonth}/${year}`
  }

  toJSON() {
    return {
      _id: this.id,
      messageId: this.messageId,
      date: this.date,
      from: this.from,
      to: this.to,
      cc: this.cc,
      bcc: this.bcc,
      subject: this.subject,
      body: this.body,
      isRead: this.isRead,
      isStarred: this.isStarred,
      isSelected: this.isSelected
    }
  }
}

export default Email
