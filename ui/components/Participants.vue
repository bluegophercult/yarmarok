<template>
    <div class="rounded-lg bg-white p-3 shadow-md ring-1 ring-black ring-opacity-5">
        <h2 class="text-xl">Учасники</h2>
        <hr class="mt-2">
        <ParticipantsList/>
        <hr class="mb-4">
        <ParticipantsCreate/>
    </div>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"

const { showError } = useNotificationStore()

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const participantStore = useParticipantStore()
watch(selectedRaffle, () => {
    if (selectedRaffle.value) {
        participantStore.getParticipants(selectedRaffle.value.id)
            .catch(e => {
                console.error(e)
                showError("Не вдалося отримати сиписок учасників!")
            })
    } else {
        participantStore.clearParticipants()
    }
})
</script>
