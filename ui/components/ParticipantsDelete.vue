<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Точно видалити учасника?</template>

        <div>
            {{ participant.name }}
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="deleteParticipant" danger full-width>Видалили</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { Participant } from "~/types/participant"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"

const props = defineProps<{
    participant: Participant
    isOpen: boolean,
    closeModal: () => void,
}>()

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const { showError } = useNotificationStore()

const participantStore = useParticipantStore()

function deleteParticipant() {
    props.closeModal()
    setTimeout(() => {
        participantStore.deleteParticipant(selectedRaffle.value!.id, props.participant.id).catch(e => {
            console.error(e)
            showError("Не вдалося видалити учасника!")
        })
    }, 200)
}
</script>
