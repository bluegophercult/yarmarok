<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>{{ participant.name }}</template>

        <div  v-if="participant.phone !== ''" class="mb-2">
            <div class="mb-1">Телефон</div>
            <div class="rounded-md bg-gray-100 p-2 shadow ring-1 ring-gray-600 ring-opacity-5">
                {{ participant.phone }}
            </div>
        </div>
        <div v-if="participant.note !== ''">
            <div class="mb-1">Опис</div>
            <p class="whitespace-pre rounded-md bg-gray-100 p-2 shadow ring-1 ring-gray-600 ring-opacity-5">
                {{ participant.note }}</p>
        </div>

        <div class="mt-4 flex flex-wrap gap-4">
            <TheButton :click="() => { closeModal(); isOpenDelete = true }" danger class="flex-1">Видалити</TheButton>
            <TheButton :click="() => { closeModal(); isOpenUpdate = true} " class="flex-1">Змінити</TheButton>
            <TheButton :click="closeModal" secondary class="flex-1">Закрити</TheButton>
        </div>
    </TheModal>

    <ParticipantsDelete :participant="participant" :is-open="isOpenDelete" :close-modal="() => isOpenDelete = false"/>
    <ParticipantsUpdate :participant="participant" :is-open="isOpenUpdate" :close-modal="() => isOpenUpdate = false"/>
</template>

<script setup lang="ts">
import { Participant } from "~/types/participant"

defineProps<{
    participant: Participant,
    isOpen: boolean,
    closeModal: () => void,
}>()

const isOpenDelete = ref(false)
const isOpenUpdate = ref(false)
</script>
