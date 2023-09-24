<template>
    <div v-if="donations.length === 0" class="text-gray-400">
        Немає внесків
    </div>
    <div class="overflow-auto max-h-[50vh] overflow-y-auto">
        <table class="table-auto w-full">
            <thead class="sticky top-0 bg-white">
            <tr>
                <th class="px-2">Учасник</th>
                <th class="px-2">Купони</th>
                <th class="px-2">Сума</th>
                <th class="px-2" colspan="2">Дії</th>
                <th class="w-2"></th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(donation, i) in donations" :key="donation.id" :class="i % 2 === 1 ? 'bg-gray-100' : ''"
                :set="participant = participantById(donation.participantId)">
                <td class="border-y px-2">{{ participant.name }}</td>
                <td class="border-y px-2">{{ donation.ticketsNumber }}</td>
                <td class="border-y px-2">{{ donation.amount }} грн</td>
                <td class="border-y w-7 text-center">
                    <Icon name="heroicons:pencil"/>
                </td>
                <td class="border-y w-7 text-center">
                    <Icon name="heroicons:trash"/>
                </td>
                <td></td>
            </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup lang="ts">
import { useRaffleStore } from "~/store/raffle"
import { usePrizeStore } from "~/store/prize"
import { useDonationStore } from "~/store/donation"
import { useParticipantStore } from "~/store/participant"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const participantStore = useParticipantStore()
const { participantById } = participantStore

const donationStore = useDonationStore()
const { donations } = storeToRefs(donationStore)

watch([ selectedRaffle, selectedPrize ], () => {
    if (selectedRaffle.value && selectedPrize.value) {
        donationStore.getDonations(selectedRaffle.value.id, selectedPrize.value.id)
    } else {
        donationStore.clearDonations()
    }
}, { immediate: true })
</script>
