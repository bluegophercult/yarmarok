<template>
    <div class="w-full md:w-1/4 sm:w-1/2 rounded-lg bg-white p-3 shadow-md ring-1 ring-black ring-opacity-5">
        <h2 class="text-xl">Призи</h2>
        <hr class="mt-2">
        <PrizesList/>
        <hr class="mb-4">
        <PrizesCreate/>
    </div>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
watch(selectedRaffle, () => {
    if (selectedRaffle.value) {
        prizeStore.getPrizes(selectedRaffle.value.id)
    } else {
        prizeStore.clearPrizes()
    }
})
</script>
