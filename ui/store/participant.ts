import { NewParticipant, Participant, Participants } from "~/types/participant"

export const useParticipantStore = defineStore({
    id: "participant-store",
    state: () => ({
        participants: <Participants>[],
    }),
    actions: {
        async getParticipants(raffleId: string) {
            const { data, error } = await useApiFetch<{
                items: Participants
            }>(`/api/raffles/${ raffleId }/participants`)
            if (error.value) {
                throw error.value
            }

            this.participants = data.value!.items || <Participants>[]
        },
        clearParticipants() {
            this.participants = []
        },
        async addParticipant(raffleId: string, newParticipant: NewParticipant) {
            const { data, error } = await useApiFetch<{
                id: string,
            }>(`/api/raffles/${ raffleId }/participants`, {
                method: "POST",
                body: newParticipant,
            })
            if (error.value) {
                throw error.value
            }

            this.participants.push(<Participant>{
                id: data.value!.id,
                ...newParticipant,
            })
        },
        updateParticipant(updatedParticipant: Participant) {
            // TODO: API call
            this.participants[this.participants.findIndex(participant => participant.id == updatedParticipant.id)] = updatedParticipant
        },
        deleteParticipant(id: string) {
            // TODO: API call
            this.participants = this.participants.filter(participant => participant.id !== id)
        },
    },
})
