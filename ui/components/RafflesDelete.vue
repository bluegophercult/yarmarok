<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Видалити розіграш?</template>

        <div>
            {{ raffle.name }}
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="deleteRaffle" full-width>Видалили</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"
import { Raffle } from "~/types/raffle"

const props = defineProps<{
    raffle: Raffle
    isOpen: boolean,
    closeModal: () => void,
}>()

const raffleStore = useRaffleStore()

function deleteRaffle() {
    raffleStore.deleteRaffle(props.raffle.id)
    props.closeModal()
}
</script>
