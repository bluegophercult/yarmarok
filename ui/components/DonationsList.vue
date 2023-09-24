<template>
    {{ donationStore.donations }}
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"
import { usePrizeStore } from "~/store/prize"
import { useDonationStore } from "~/store/donation"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const donationStore = useDonationStore()

watch([ selectedRaffle, selectedPrize ], () => {
    if (selectedRaffle.value && selectedPrize.value) {
        donationStore.getDonations(selectedRaffle.value.id, selectedPrize.value.id)
    } else {
        donationStore.clearDonations()
    }
}, { immediate: true })
</script>
