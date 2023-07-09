<template>
    <div class="w-1/4 rounded-lg bg-white p-3 shadow-md ring-1 ring-black ring-opacity-5">
        <h2 class="text-xl">Призи</h2>
        <hr class="mt-2">
        <PrizesList/>
        <hr class="mb-4">
        <PrizesCreate/>
    </div>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"

const prizeStore = usePrizeStore()
prizeStore.getPrizes()
</script>
