<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Точно видалити розіграш?</template>

        <div>
            {{ raffle.name }}
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="deleteRaffle" danger full-width>Видалили</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"
import { Raffle } from "~/types/raffle"
import { useNotificationStore } from "~/store/notification"

const props = defineProps<{
    raffle: Raffle
    isOpen: boolean,
    closeModal: () => void,
}>()

const { showError } = useNotificationStore()

const raffleStore = useRaffleStore()

function deleteRaffle() {
    props.closeModal()
    setTimeout(() => {
        raffleStore.deleteRaffle(props.raffle.id).catch(e => {
            console.error(e)
            showError("Не вдалося видалити розіграш!")
        })
    }, 200)
}
</script>
