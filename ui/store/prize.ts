import { NewPrize, Prize, Prizes } from "~/types/prize"

export const usePrizeStore = defineStore({
    id: "prize-store",
    state: () => ({
        prizes: <Prizes>{},
        selectedPrize: <Prize | null>null,
    }),
    actions: {
        getPrizes() {
            this.prizes = <Prizes>[
                { id: "1", name: "Прапор" },
                { id: "2", name: "Печенька" },
                { id: "3", name: "Шкарпетки" },
            ]
            this.selectFirstPrize()
        },
        addPrize(newPrize: NewPrize) {
            // TODO: API call
            this.prizes.push(<Prize>{
                id: `${ this.prizes.length + 1 }`,
                ...newPrize,
            })
            this.selectLastPrize()
        },
        updatePrize(updatedPrize: Prize) {
            // TODO: API call
            this.prizes[this.prizes.findIndex(prize => prize.id == updatedPrize.id)] = updatedPrize
        },
        deletePrize(id: string) {
            // TODO: API call
            this.prizes = this.prizes.filter(prize => prize.id !== id)
            this.selectFirstPrize()
        },
        selectFirstPrize() {
            this.selectedPrize = this.prizes.length === 0 ? null : this.prizes[0]
        },
        selectLastPrize() {
            this.selectedPrize = this.prizes.length === 0 ? null : this.prizes[this.prizes.length - 1]
        },
    },
})
