import { NewRaffle, Raffle, Raffles } from "~/types/raffle"

export const useRaffleStore = defineStore({
    id: "raffle-store",
    state: () => ({
        raffles: <Raffles>[
            { id: "1", name: "Фестиваль їжі" },
            { id: "2", name: "Atlas weekend" },
        ],
        selectedRaffle: <Raffle>{},
    }),
    actions: {
        addRaffle(newRaffle: NewRaffle) {
            // TODO: API call
            this.raffles.push(<Raffle>{
                id: `${ this.raffles.length + 1 }`,
                ...newRaffle,
            })
        },
        deleteRaffle(id: string) {
            // TODO: API call
            this.raffles = this.raffles.filter(raffle => raffle.id !== id)
        },
    },
})
