<template>
    <div class="flex gap-2">
        <RafflesSelect/>
        <RafflesMore/>
        <RafflesCreate/>
    </div>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"

const { showError } = useNotificationStore()

const raffleStore = useRaffleStore()
onMounted(() => {
    raffleStore.getRaffles()
        .catch(e => {
            console.error(e)
            showError("Не вдалося отримати сиписок розіграшів!")
        })
})
</script>
