<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Точно видалити приз?</template>

        <div>
            {{ prize.name }}
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="deletePrize" danger full-width>Видалили</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { Prize } from "~/types/prize"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"

const props = defineProps<{
    prize: Prize
    isOpen: boolean,
    closeModal: () => void,
}>()

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const { showError } = useNotificationStore()

const prizeStore = usePrizeStore()

function deletePrize() {
    props.closeModal()
    setTimeout(() => {
        prizeStore.deletePrize(selectedRaffle.value!.id, props.prize.id).catch(e => {
            console.error(e)
            showError("Не вдалося видалити приз!")
        })
    }, 200)
}
</script>
