class Email {
  id: string
  message_id: string
  date: Date
  from: string
  to: string[]
  cc: string[]
  bcc: string[]
  subject: string
  body: string
  read: boolean
  starred: boolean
  selected: boolean

  constructor(
    id: string,
    message_id: string,
    date: Date,
    from: string,
    to: string[],
    cc: string[],
    bcc: string[],
    subject: string,
    body: string,
    read: boolean,
    starred: boolean,
    selected: boolean = false
  ) {
    this.id = id
    this.message_id = message_id
    this.date = date
    this.from = from
    this.to = to
    this.cc = cc
    this.bcc = bcc
    this.subject = subject
    this.body = body
    this.read = read
    this.starred = starred
    this.selected = selected
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

// export class
export default Email
