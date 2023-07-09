import { NewPrize, Prize, Prizes } from "~/types/prize"

export const usePrizeStore = defineStore({
    id: "prize-store",
    state: () => ({
        prizes: <Prizes>{},
    }),
    actions: {
        getPrizes() {
            this.prizes = <Prizes>[
                { id: "1", name: "Прапор" },
                { id: "2", name: "Печенька" },
                { id: "3", name: "Шкарпетки" },
            ]
        },
        addPrize(newPrize: NewPrize) {
            // TODO: API call
            this.prizes.push(<Prize>{
                id: `${ this.prizes.length + 1 }`,
                ...newPrize,
            })
        },
        updatePrize(updatedPrize: Prize) {
            // TODO: API call
            this.prizes[this.prizes.findIndex(prize => prize.id == updatedPrize.id)] = updatedPrize
        },
        deletePrize(id: string) {
            // TODO: API call
            this.prizes = this.prizes.filter(prize => prize.id !== id)
        },
    },
})
