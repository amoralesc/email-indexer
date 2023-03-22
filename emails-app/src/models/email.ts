interface Email {
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
    selected: boolean | false
}

export type { Email }
