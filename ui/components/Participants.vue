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

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const ParticipantStore = useParticipantStore()
watch(selectedRaffle, () => {
    if (selectedRaffle.value) {
        ParticipantStore.getParticipants(selectedRaffle.value.id)
    } else {
        ParticipantStore.clearParticipants()
    }
})
</script>
