<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Змінити розіграш</template>

        <form @submit.prevent="updateRaffle">
            <div class="flex flex-col gap-2">
                <TheInput v-model="updatedRaffle.name" label="Назва" :placeholder="raffle.name" required/>
                <TheTextArea v-model="updatedRaffle.note" label="Опис" :placeholder="raffle.note"/>
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
