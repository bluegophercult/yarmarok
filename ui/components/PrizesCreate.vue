<template>
    <TheButton @click="alert('Add')" full-width>Додати приз</TheButton>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"

const prizeStore = usePrizeStore()
</script>
