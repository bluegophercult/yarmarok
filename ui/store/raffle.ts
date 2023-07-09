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
                { id: "1", name: "Фестиваль їжі" },
                { id: "2", name: "Atlas weekend" },
            ]
            this.selectLastRaffle()
        },
        addRaffle(newRaffle: NewRaffle) {
            // TODO: API call
            this.raffles.push(<Raffle>{
                id: `${ this.raffles.length + 1 }`,
                ...newRaffle,
            })
            this.selectLastRaffle()
        },
        updateRaffle(updatedRaffle: Raffle) {
            // TODO: API call
            this.raffles[this.raffles.findIndex(raffle => raffle.id == updatedRaffle.id)] = updatedRaffle
            this.selectedRaffle = updatedRaffle
        },
        deleteRaffle(id: string) {
            // TODO: API call
            this.raffles = this.raffles.filter(raffle => raffle.id !== id)
            this.selectLastRaffle()
        },
        selectLastRaffle() {
            this.selectedRaffle = this.raffles.length === 0 ? null : this.raffles[this.raffles.length - 1]
        },
    },
})
