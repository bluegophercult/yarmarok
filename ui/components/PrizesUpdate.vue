<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Змінити приз</template>

        <form @submit.prevent="updatePrize">
            <div class="flex flex-col gap-2">
                <TheInput v-model="updatedPrize.name" :placeholder="prize.name" label="Назва" required/>
                <TheInput v-model="updatedPrize.ticketCost" :placeholder="prize.ticketCost" number :min="1"
                          label="Ціна купону" required :disabled="donations.length > 0"/>
                <TheTextArea v-model="updatedPrize.description" :placeholder="prize.description" label="Опис"/>
            </div>

            <transition name="m-fade">
                <p v-show="errorMsg" class="mt-4 flex items-center gap-2 text-sm text-red-500 transition duration-200">
                    <Icon name="heroicons:exclamation-triangle" class="h-5 w-5"/>
                    {{ errorMsg }}
                </p>
            </transition>

            <div class="mt-4 flex gap-4">
                <TheButton submit full-width>Зберегти</TheButton>
                <TheButton :click="closeModal" secondary full-width>Закрити</TheButton>
            </div>
        </form>
    </TheModal>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { newPrizeSchema, Prize } from "~/types/prize"
import { ValidationError } from "yup"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"
import { useDonationStore } from "~/store/donation"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const donationStore = useDonationStore()
const { donations } = storeToRefs(donationStore)

const { showError } = useNotificationStore()

const prizeStore = usePrizeStore()

const props = defineProps<{
    prize: Prize,
    isOpen: boolean,
    closeModal: () => void,
}>()

const errorMsg = ref("")
const updatedPrize = ref(<Prize>{ ...props.prize })
onBeforeUpdate(() => {
    setTimeout(() => {
        updatedPrize.value = { ...props.prize }
        errorMsg.value = ""
    }, 200)
})

function updatePrize() {
    newPrizeSchema.validate(updatedPrize.value)
        .then(() => {
            prizeStore.updatePrize(selectedRaffle.value!.id, updatedPrize.value).catch(e => {
                console.error(e)
                showError("Не вдалося змінити приз!")
            })
            props.closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
