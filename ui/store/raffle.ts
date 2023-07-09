import { NewRaffle, Raffle, Raffles } from "~/types/raffle"

export const useRaffleStore = defineStore({
    id: "raffle-store",
    state: () => ({
        raffles: <Raffles>{},
        selectedRaffle: <Raffle | null>null,
    }),
    actions: {
        getRaffles() {
            this.raffles = <Raffles>[
                { id: "2", name: "Atlas weekend" },
                { id: "1", name: "Фестиваль їжі" },
            ]
            this.selectFirstRaffle()
        },
        addRaffle(newRaffle: NewRaffle) {
            // TODO: API call
            this.raffles.unshift(<Raffle>{
                id: `${ this.raffles.length + 1 }`,
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
            this.selectFirstRaffle()
        },
        selectFirstRaffle() {
            this.selectedRaffle = this.raffles.length === 0 ? null : this.raffles[0]
        },
    },
})
