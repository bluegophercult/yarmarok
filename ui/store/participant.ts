import { NewParticipant, Participant, Participants } from "~/types/participant"

export const useParticipantStore = defineStore({
    id: "participant-store",
    state: () => ({
        participants: <Participants>[],
    }),
    actions: {
        getParticipants(raffleId: string) {
            // TODO: API call
            this.participants = <Participants>[
                { id: "1", name: "Оксана" },
                { id: "2", name: "Ярослав" },
            ]
        },
        clearParticipants() {
            this.participants = []
        },
        addParticipant(newParticipant: NewParticipant) {
            // TODO: API call
            this.participants.push(<Participant>{
                id: `${ this.participants.length + 1 }`,
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
