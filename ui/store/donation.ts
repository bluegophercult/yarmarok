import { NewDonation, Donation, Donations } from "~/types/donation"

export const useDonationStore = defineStore({
    id: "donation-store",
    state: () => ({
        donations: <Donations>[],
        selectedDonation: <Donation | null>null,
    }),
    actions: {
        async getDonations(raffleId: string, prizeId: string) {
            const { data, error } = await useApiFetch<{
                items: Donations
            }>(`/api/raffles/${ raffleId }/prizes/${ prizeId }/donations`)
            if (error.value) {
                throw error.value
            }

            this.donations = data.value!.items || <Donations>[]
            this.selectFirstDonation()
        },
        clearDonations() {
            this.donations = []
        },
        async addDonation(raffleId: string, prizeId: string, newDonation: NewDonation, ticketCost: number) {
            const { data, error } = await useApiFetch<{
                id: string,
            }>(`/api/raffles/${ raffleId }/prizes/${ prizeId }/donations`, {
                method: "POST",
                body: newDonation,
            })
            if (error.value) {
                throw error.value
            }

            this.donations.push(<Donation>{
                id: data.value!.id,
                ticketsNumber: Math.floor(newDonation.amount / ticketCost),
                ...newDonation,
            })
            this.selectLastDonation()
        },
        async updateDonation(raffleId: string, prizeId: string, updatedDonation: Donation) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/prizes/${ prizeId }/donations/${ updatedDonation.id }`, {
                    method: "PUT",
                    body: updatedDonation,
                })
            if (error.value) {
                throw error.value
            }

            this.donations[this.donations.findIndex(donation => donation.id == updatedDonation.id)] = updatedDonation
            if (this.selectedDonation && this.selectedDonation.id === updatedDonation.id) {
                this.selectedDonation = updatedDonation
            }
        },
        async deleteDonation(raffleId: string, prizeId: string, id: string) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/prizes/${ prizeId }/donations/${ id }`, {
                    method: "DELETE",
                })
            if (error.value) {
                throw error.value
            }

            this.donations = this.donations.filter(donation => donation.id !== id)
            this.selectLastDonation()
        },
        selectFirstDonation() {
            this.selectedDonation = this.donations.length === 0 ? null : this.donations[0]
        },
        selectLastDonation() {
            this.selectedDonation = this.donations.length === 0 ? null : this.donations[this.donations.length - 1]
        },
    },
})
