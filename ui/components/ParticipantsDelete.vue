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

const props = defineProps<{
    participant: Participant
    isOpen: boolean,
    closeModal: () => void,
}>()

const participantStore = useParticipantStore()

function deleteParticipant() {
    props.closeModal()
    setTimeout(() => {
        participantStore.deleteParticipant(props.participant.id)
    }, 200)
}
</script>
