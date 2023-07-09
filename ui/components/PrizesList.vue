<template>
    <ul class="py-2">
        <li v-for="prize in prizes" :key="prize.id" class="flex h-8 items-center rounded-md px-2 hover:bg-teal-100">
            {{ prize.name }}
        </li>
    </ul>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"

const prizeStore = usePrizeStore()
const { prizes } = storeToRefs(prizeStore)
</script>
