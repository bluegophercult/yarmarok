<template>
    <TheButton @click="openModal" :disabled="!selectedRaffle || !selectedPrize" full-width>Додати внесок</TheButton>

    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Додати новий внесок</template>

        <form @submit.prevent="addDonation">
            <div class="flex flex-col gap-2">
                <TheInput v-model="newDonation.amount" label="Сума" :min="1" number required/>

                <label class="flex flex-col gap-1">
                    <span>Учасник <span class="text-sm text-red-400">*</span></span>
                    <HeadlessCombobox v-model="selectedParticipant">
                        <div class="relative mt-1">
                            <div class="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm">
                                <HeadlessComboboxInput
                                        class="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
                                        :displayValue="participant => participant ? participant.name : ''"
                                        @change="query = $event.target.value"
                                />
                                <HeadlessComboboxButton class="absolute inset-y-0 right-0 flex items-center pr-2">
                                    <Icon name="heroicons:chevron-up-down"
                                          class="h-5 w-5 text-gray-600 transition duration-200 group-hover:text-teal-400"/>
                                </HeadlessComboboxButton>
                            </div>

                            <HeadlessTransitionRoot
                                    leave="transition ease-in duration-100"
                                    leaveFrom="opacity-100"
                                    leaveTo="opacity-0"
                                    @after-leave="query = ''"
                            >
                                <HeadlessComboboxOptions
                                        class="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                                    <div v-if="filteredParticipants.length === 0 && query !== ''"
                                         class="relative cursor-default select-none py-2 px-4 text-gray-700">
                                        Нікого не знайдено
                                    </div>

                                    <HeadlessComboboxOption
                                            v-for="participant in filteredParticipants"
                                            as="template"
                                            :key="participant.id"
                                            :value="participant"
                                            v-slot="{ selected, active }"
                                    >
                                        <li class="relative cursor-default select-none py-2 pl-10 pr-4"
                                            :class="{ 'bg-teal-600 text-white': active, 'text-gray-900': !active }">
                                        <span class="block truncate"
                                              :class="{ 'font-medium': selected, 'font-normal': !selected }">
                                          {{ participant.name }}
                                        </span>
                                            <span v-if="selected"
                                                  class="absolute inset-y-0 left-0 flex items-center pl-3"
                                                  :class="{ 'text-white': active, 'text-teal-600': !active }"
                                            >
                                            <Icon name="heroicons:chevron-right" class="h-5 w-5 text-teal-400"/>
                                        </span>
                                        </li>
                                    </HeadlessComboboxOption>
                                </HeadlessComboboxOptions>
                            </HeadlessTransitionRoot>
                        </div>
                    </HeadlessCombobox>
                </label>
            </div>

            <transition name="m-fade">
                <p v-show="errorMsg" class="mt-4 flex items-center gap-2 text-sm text-red-500 transition duration-200">
                    <Icon name="heroicons:exclamation-triangle" class="h-5 w-5"/>
                    {{ errorMsg }}
                </p>
            </transition>

            <div class="mt-4 flex gap-4">
                <TheButton submit full-width>Додати</TheButton>
                <TheButton :click="closeModal" secondary full-width>Закрити</TheButton>
            </div>
        </form>
    </TheModal>
</template>

<script setup lang="ts">
import { useParticipantStore } from "~/store/participant"
import { Ref } from "@vue/reactivity"
import { ValidationError } from "yup"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"
import { usePrizeStore } from "~/store/prize"
import { useDonationStore } from "~/store/donation"
import { NewDonation, newDonationSchema } from "~/types/donation"
import { Participant } from "~/types/participant"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const participantStore = useParticipantStore()
const { participants } = storeToRefs(participantStore)

const { showError } = useNotificationStore()

const donationStore = useDonationStore()
const newDonation: Ref<NewDonation> = ref(<NewDonation>{
    amount: 0,
    participantId: "",
})

const isOpen = ref(false)
const errorMsg = ref("")

const query = ref("")
const selectedParticipant: Ref<Participant | null> = ref(null)
const filteredParticipants = computed(() =>
    query.value === ""
        ? participants.value
        : participants.value.filter(participant => {
            return participant.name.toLowerCase().includes(query.value.toLowerCase())
        }),
)

function openModal() {
    errorMsg.value = ""
    isOpen.value = true
}

function closeModal() {
    isOpen.value = false
    setTimeout(() => {
        newDonation.value = <NewDonation>{
            amount: 0,
            participantId: "",
        }
        selectedParticipant.value = null
    }, 200)
}

function addDonation() {
    if (selectedParticipant.value) {
        newDonation.value.participantId = selectedParticipant.value.id
    }
    newDonationSchema.validate(newDonation.value)
        .then(() => {
            donationStore.addDonation(
                selectedRaffle.value!.id,
                selectedPrize.value!.id,
                newDonation.value,
                selectedPrize.value!.ticketCost,
            ).catch(e => {
                console.error(e)
                showError("Не вдалося створити внесок!")
            })
            closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
