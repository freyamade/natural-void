export class Alert {
    constructor(type: string, message: string) {
        const notif = document.createElement('div')
        notif.classList.add('notification', `is-${type}`)
        const text = document.createTextNode(message)
        notif.appendChild(text)
        document.getElementById('notif-center')!.appendChild(notif)
    }
}
