<template>
    <ul class="max-h-80 overflow-auto py-2">
        <li v-for="prize in prizes" :key="prize.id"
            class="flex h-8 items-center justify-between gap-2 rounded-md px-2 hover:text-teal-950 hover:cursor-pointer hover:bg-teal-100"
            @click="selectedPrize = prize">
            <Icon v-if="selectedPrize && prize.id === selectedPrize.id" name="heroicons:chevron-right"
                  class="h-5 w-5 text-teal-400"/>
            <span class="block flex-grow truncate">{{ prize.name }}</span>
            <Icon v-if="selectedPrize && prize.id === selectedPrize.id" name="heroicons:chevron-left"
                  class="h-5 w-5 text-teal-400"/>
        </li>
        <li v-if="prizes.length === 0" class="text-gray-400">
            Пусто
        </li>
    </ul>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { useStateStore } from "~/store/state"

const prizeStore = usePrizeStore()
const { prizes, selectedPrize } = storeToRefs(prizeStore)

const stateStore = useStateStore()

watch(selectedPrize, () => {
    if (selectedPrize.value) {
        stateStore.selectedPrize = selectedPrize.value.id
        stateStore.update()
    }
}, { immediate: true })
</script>
