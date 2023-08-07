<template>
    <ul class="max-h-80 overflow-auto py-2">
        <li v-for="participant in participants" :key="participant.id"
            class="flex h-8 items-center justify-between gap-2 rounded-md px-2 hover:text-teal-950 hover:cursor-pointer hover:bg-teal-100">
            <span class="block flex-grow truncate">{{ participant.name }}</span>
        </li>
        <li v-if="participants.length === 0" class="text-gray-400">
            Пусто
        </li>
    </ul>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"

const participantStore = useParticipantStore()
const { participants } = storeToRefs(participantStore)
</script>
