<template>
    <ul class="max-h-80 overflow-auto py-2">
        <li v-for="participant in participants" :key="participant.id" @click="openParticipantView(participant)"
            class="flex h-8 items-center justify-between gap-2 rounded-md px-2 hover:text-teal-950 hover:cursor-pointer hover:bg-teal-100">
            <span class="block flex-grow truncate">{{ participant.name }}</span>
        </li>
        <li v-if="participants.length === 0" class="text-gray-400">
            Пусто
        </li>
    </ul>

    <ParticipantsView v-if="selectedParticipant" :participant="selectedParticipant" :is-open="isOpenView"
                      :close-modal="() => isOpenView = false"/>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { Participant } from "~/types/participant"
import { Ref } from "@vue/reactivity"

const participantStore = useParticipantStore()
const { participants } = storeToRefs(participantStore)

const isOpenView = ref(false)
const selectedParticipant: Ref<Participant | undefined> = ref(undefined)

function openParticipantView(participant: Participant) {
    selectedParticipant.value = participant
    isOpenView.value = true
}
</script>
