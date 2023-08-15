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
                items: Raffles
            }>("/api/raffles")
            if (error.value) {
                if (error.value.statusCode && error.value.statusCode === 500 && !window.location.href.endsWith("/api/login")) {
                    navigateTo("/api/login", { external: true, replace: true, redirectCode: 302 })
                    return
                }
                throw error.value
            }

            this.raffles = data.value!.items || <Raffles>[]
            this.selectFirstRaffle()
        },
        async addRaffle(newRaffle: NewRaffle) {
            const { data, error } = await useApiFetch<{
                id: string,
            }>("/api/raffles", {
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
        async updateRaffle(updatedRaffle: Raffle) {
            const { error } = await useApiFetch(`/api/raffles/${ updatedRaffle.id }`, {
                method: "PUT",
                body: updatedRaffle,
            })
            if (error.value) {
                throw error.value
            }

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
