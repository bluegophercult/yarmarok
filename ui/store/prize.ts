import { NewPrize, Prize, Prizes } from "~/types/prize"
import { useStateStore } from "~/store/state"

export const usePrizeStore = defineStore({
    id: "prize-store",
    state: () => ({
        prizes: <Prizes>[],
        selectedPrize: <Prize | null>null,
    }),
    actions: {
        async getPrizes(raffleId: string) {
            const { data, error } = await useApiFetch<{
                items: Prizes
            }>(`/api/raffles/${ raffleId }/prizes`)
            if (error.value) {
                throw error.value
            }

            this.prizes = data.value!.items || <Prizes>[]

            const stateStore = useStateStore()
            if (stateStore.selectedRaffle) {
                const selected = this.prizes.find(p => p.id == stateStore.selectedPrize)
                if (selected) {
                    this.selectedPrize = selected
                } else {
                    this.selectFirstPrize()
                }
            } else {
                this.selectFirstPrize()
            }
        },
        clearPrizes() {
            this.prizes = []
        },
        async addPrize(raffleId: string, newPrize: NewPrize) {
            const { data, error } = await useApiFetch<{
                id: string,
            }>(`/api/raffles/${ raffleId }/prizes`, {
                method: "POST",
                body: newPrize,
            })
            if (error.value) {
                throw error.value
            }

            this.prizes.push(<Prize>{
                id: data.value!.id,
                ...newPrize,
            })
            this.selectLastPrize()
        },
        async updatePrize(raffleId: string, updatedPrize: Prize) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/prizes/${ updatedPrize.id }`, {
                    method: "PUT",
                    body: updatedPrize,
                })
            if (error.value) {
                throw error.value
            }

            this.prizes[this.prizes.findIndex(prize => prize.id == updatedPrize.id)] = updatedPrize
            if (this.selectedPrize && this.selectedPrize.id === updatedPrize.id) {
                this.selectedPrize = updatedPrize
            }
        },
        async deletePrize(raffleId: string, id: string) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/prizes/${ id }`, {
                    method: "DELETE",
                })
            if (error.value) {
                throw error.value
            }

            this.prizes = this.prizes.filter(prize => prize.id !== id)
            this.selectLastPrize()
        },
        async playPrize(raffleId: string, id: string) {
            const { error } = await useApiFetch<{
                isPlayed: boolean
            }>(
                `/api/raffles/${ raffleId }/prizes/${ id }/play`, {
                    method: "GET",
                })
            if (error.value) {
                throw error.value
            }
        },
        selectFirstPrize() {
            this.selectedPrize = this.prizes.length === 0 ? null : this.prizes[0]
        },
        selectLastPrize() {
            this.selectedPrize = this.prizes.length === 0 ? null : this.prizes[this.prizes.length - 1]
        },
    },
})
