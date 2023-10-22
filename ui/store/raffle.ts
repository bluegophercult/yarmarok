import { NewRaffle, Raffle, Raffles } from "~/types/raffle"
import { useStateStore } from "~/store/state"

export const useRaffleStore = defineStore({
    id: "raffle-store",
    state: () => ({
        raffles: <Raffles>[],
        selectedRaffle: <Raffle | null>null,
        rafflesLoaded: false,
    }),
    actions: {
        async getRaffles() {
            this.rafflesLoaded = false
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

            const stateStore = useStateStore()
            if (stateStore.selectedRaffle) {
                const selected = this.raffles.find(r => r.id == stateStore.selectedRaffle)
                if (selected) {
                    this.selectedRaffle = selected
                } else {
                    this.selectFirstRaffle()
                }
            } else {
                this.selectFirstRaffle()
            }

            this.rafflesLoaded = true
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
        async deleteRaffle(id: string) {
            const { error } = await useApiFetch(`/api/raffles/${ id }`, {
                method: "DELETE",
            })
            if (error.value) {
                throw error.value
            }

            this.raffles = this.raffles.filter(raffle => raffle.id !== id)
            this.selectFirstRaffle()
        },
        selectFirstRaffle() {
            this.selectedRaffle = this.raffles.length === 0 ? null : this.raffles[0]
        },
    },
})
