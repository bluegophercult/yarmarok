export type NotificationType = "info" | "error"

export const useNotificationStore = defineStore({
    id: "notification-store",
    state: () => ({
        shown: false,
        type: <NotificationType>"info",
        message: <any>"",
    }),
    actions: {
        showNotification(message: any, type: NotificationType) {
            this.message = message
            this.type = type
            this.shown = true
        },
        showInfo(message: any) {
            this.showNotification(message, "info")
        },
        showError(message: any) {
            this.showNotification(message, "error")
        },
        hideNotification() {
            this.shown = false
        },
    },
})
