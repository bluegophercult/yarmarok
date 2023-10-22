const selectedRaffleKey = "selected-raffle"
const selectedPrizeKey = "selected-prize"

export const useStateStore = defineStore({
    id: "state-store",
    state: () => ({
        selectedRaffle: <string | null>null,
        selectedPrize: <string | null>null,
    }),
    actions: {
        init() {
            if (!process.client) {
                return
            }

            this.selectedRaffle = localStorage.getItem(selectedRaffleKey)
            this.selectedPrize = localStorage.getItem(selectedPrizeKey)
        },
        update() {
            if (!process.client) {
                return
            }

            if (this.selectedRaffle) {
                localStorage.setItem(selectedRaffleKey, this.selectedRaffle)
            } else {
                localStorage.removeItem(selectedRaffleKey)
            }
            if (this.selectedPrize) {
                localStorage.setItem(selectedPrizeKey, this.selectedPrize)
            } else {
                localStorage.removeItem(selectedPrizeKey)
            }
        },
    },
})
