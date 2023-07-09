<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Змінити розіграш</template>

        <form @submit.prevent="updateRaffle">
            <div class="flex flex-col gap-2">
                <TheInput v-model="updatedRaffle.name" label="Назва" required/>
                <TheTextArea v-model="updatedRaffle.note" label="Опис"/>
            </div>

            <transition name="m-fade">
                <p v-show="errorMsg" class="flex gap-2 items-center transition duration-200 mt-4 text-sm text-red-500">
                    <Icon name="heroicons:exclamation-triangle" class="w-5 h-5"/>
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
import { useRaffleStore } from "~/store/raffle"
import { newRaffleSchema, Raffle } from "~/types/raffle"
import { ValidationError } from "yup"

const props = defineProps<{
    raffle: Raffle
    isOpen: boolean,
    closeModal: () => void,
}>()

const raffleStore = useRaffleStore()

const errorMsg = ref("")
const updatedRaffle = ref(<Raffle>{ ...props.raffle })
onBeforeUpdate(() => {
    setTimeout(() => {
        updatedRaffle.value = { ...props.raffle }
        errorMsg.value = ""
    }, 200)
})

function updateRaffle() {
    newRaffleSchema.validate(updatedRaffle.value)
        .then(() => {
            raffleStore.updateRaffle(updatedRaffle.value)
            props.closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
