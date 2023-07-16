import { NewRaffle, Raffle, Raffles } from "~/types/raffle"

export const useRaffleStore = defineStore({
    id: "raffle-store",
    state: () => ({
        raffles: <Raffles>[],
        selectedRaffle: <Raffle | null>null,
    }),
    actions: {
        async getRaffles() {
            const { data, error } = await useApiFetch<{
                raffles: Raffles
            }>("/raffles")
            if (error.value) {
                throw error.value
            }

            this.raffles = data.value!.raffles || <Raffles>[]
            this.selectFirstRaffle()
        },
        async addRaffle(newRaffle: NewRaffle) {
            const { data, error } = await useApiFetch<{
                id: string,
            }>("/raffles", {
                method: "POST",
                body: newRaffle,
            })
            if (error.value) {
                throw error.value
            }

            this.raffles.unshift(<Raffle>{
                id: data.value!.id,
                ...newRaffle,
            })
            this.selectFirstRaffle()
        },
        updateRaffle(updatedRaffle: Raffle) {
            // TODO: API call
            this.raffles[this.raffles.findIndex(raffle => raffle.id == updatedRaffle.id)] = updatedRaffle
            this.selectedRaffle = updatedRaffle
        },
        deleteRaffle(id: string) {
            // TODO: API call
            this.raffles = this.raffles.filter(raffle => raffle.id !== id)
        },
        selectFirstRaffle() {
            this.selectedRaffle = this.raffles.length === 0 ? null : this.raffles[0]
        },
    },
})
