<template>
    <div class="flex gap-2">
        <RafflesSelect/>
        <RafflesMore/>
        <RafflesCreate/>
    </div>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"

const raffleStore = useRaffleStore()
raffleStore.getRaffles()
</script>
