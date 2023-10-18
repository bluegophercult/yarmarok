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
        async updateParticipant(raffleId: string, updatedParticipant: Participant) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/participants/${ updatedParticipant.id }`, {
                    method: "PUT",
                    body: updatedParticipant,
                })
            if (error.value) {
                throw error.value
            }

            this.participants[this.participants.findIndex(participant => participant.id == updatedParticipant.id)] = updatedParticipant
        },
        async deleteParticipant(raffleId: string, id: string) {
            const { error } = await useApiFetch(
                `/api/raffles/${ raffleId }/participants/${ id }`, {
                    method: "DELETE",
                })
            if (error.value) {
                throw error.value
            }

            this.participants = this.participants.filter(participant => participant.id !== id)
        },
        participantById(id: string): Participant | undefined {
            return this.participants.find(p => p.id == id)
        },
    },
})
