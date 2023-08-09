<template>
    <div class="rounded-lg bg-white p-3 shadow-md ring-1 ring-black ring-opacity-5">
        <h2 class="text-xl">Деталі призу</h2>
        <hr class="mt-2">
        <div v-if="selectedPrize">
            {{ selectedPrize.name }}
        </div>
        <div v-else>
            Не вибрано
        </div>
    </div>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)
</script>
