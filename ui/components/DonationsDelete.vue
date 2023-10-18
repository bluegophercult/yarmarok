<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Точно видалити внесок?</template>

        <div>
            {{ participantById(donation.participantId).name }} - {{ donation.amount }} грн
        </div>

        <div class="mt-4 flex gap-4">
            <TheButton :click="deleteDonation" danger full-width>Видалили</TheButton>
            <TheButton :click="closeModal" full-width secondary>Закрити</TheButton>
        </div>
    </TheModal>
</template>

<script setup lang="ts">
import { useDonationStore } from "~/store/donation"
import { Donation } from "~/types/donation"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"
import { useParticipantStore } from "~/store/participant"
import { usePrizeStore } from "~/store/prize"

const props = defineProps<{
    donation: Donation
    isOpen: boolean,
    closeModal: () => void,
}>()

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const participantStore = useParticipantStore()
const { participantById } = participantStore

const { showError } = useNotificationStore()

const donationStore = useDonationStore()

function deleteDonation() {
    props.closeModal()
    setTimeout(() => {
        donationStore.deleteDonation(selectedRaffle.value!.id, selectedPrize.value!.id, props.donation.id).catch(e => {
            console.error(e)
            showError("Не вдалося видалити внесок!")
        })
    }, 200)
}
</script>
