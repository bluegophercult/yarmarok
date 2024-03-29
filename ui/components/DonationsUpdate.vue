<template>
    <TheModal :is-open="isOpen" :close-modal="closeModal">
        <template #title>Змінити внесок</template>

        <form @submit.prevent="updateDonation">
            <div class="flex flex-col gap-2">
                <TheInput v-model="updatedDonation.amount" :placeholder="donation.amount" label="Сума" :min="1" number
                          required/>

                <div class="flex flex-col gap-1">
                    <label for="donation-update-participant">Учасник <span class="text-sm text-red-400">*</span></label>
                    <HeadlessCombobox v-model="selectedParticipant">
                        <div class="relative mt-1">
                            <div class="relative w-full">
                                <HeadlessComboboxInput
                                        class="w-full border-0 py-2 pl-3 pr-10 rounded-md ring-1 ring-black ring-opacity-30 transition duration-100 focus:outline-none focus:ring-1 focus:ring-teal-400"
                                        :displayValue="participant => participant ? participant.name : ''"
                                        @change="query = $event.target.value" id="donation-update-participant"
                                        :placeholder="participantById(donation.participantId)?.name"
                                />
                                <HeadlessComboboxButton class="absolute inset-y-0 right-0 flex items-center pr-2">
                                    <Icon name="heroicons:chevron-up-down"
                                          class="h-5 w-5 text-gray-600 transition duration-200 hover:text-teal-400"/>
                                </HeadlessComboboxButton>
                            </div>

                            <HeadlessTransitionRoot
                                    enter="transition ease-out duration-200"
                                    enter-from="opacity-0"
                                    enter-to="opacity-100"
                                    leave="transition ease-in duration-200"
                                    leave-from="opacity-100"
                                    leave-to="opacity-0"
                                    @after-leave="query = ''"
                            >
                                <HeadlessComboboxOptions
                                        class="mt-2 max-h-36 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                    <div v-if="filteredParticipants.length === 0 && query !== ''"
                                         class="select-none py-2 px-4 text-gray-700">
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
                                            :class="{ 'bg-teal-100 text-teal-950': active, 'text-gray-900': !active }">
                                            <span class="block truncate"
                                                  :class="{ 'font-medium': selected, 'font-normal': !selected }">
                                              {{ participant.name }}
                                            </span>
                                            <span v-if="selected"
                                                  class="absolute inset-y-0 left-0 flex items-center pl-3"
                                                  :class="{ 'text-white': active, 'text-teal-600': !active }">
                                                <Icon name="heroicons:chevron-right" class="h-5 w-5 text-teal-400"/>
                                            </span>
                                        </li>
                                    </HeadlessComboboxOption>
                                </HeadlessComboboxOptions>
                            </HeadlessTransitionRoot>
                        </div>
                    </HeadlessCombobox>
                </div>
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
import { useDonationStore } from "~/store/donation"
import { newDonationSchema, Donation } from "~/types/donation"
import { ValidationError } from "yup"
import { useRaffleStore } from "~/store/raffle"
import { useNotificationStore } from "~/store/notification"
import { usePrizeStore } from "~/store/prize"
import { Ref } from "@vue/reactivity"
import { Participant } from "~/types/participant"
import { useParticipantStore } from "~/store/participant"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const participantStore = useParticipantStore()
const { participants } = storeToRefs(participantStore)
const { participantById } = participantStore

const { showError } = useNotificationStore()

const donationStore = useDonationStore()

const props = defineProps<{
    donation: Donation,
    isOpen: boolean,
    closeModal: () => void,
}>()

const errorMsg = ref("")
const updatedDonation = ref(<Donation>{ ...props.donation })
onBeforeUpdate(() => {
    setTimeout(() => {
        updatedDonation.value = { ...props.donation }
        selectedParticipant.value = participantById(updatedDonation.value.participantId)!
        errorMsg.value = ""
    }, 200)
})

const query = ref("")
const selectedParticipant: Ref<Participant | null> = ref(participantById(updatedDonation.value.participantId)!)
const filteredParticipants = computed(() =>
    query.value === ""
        ? participants.value
        : participants.value.filter(participant => {
            return participant.name.toLowerCase().includes(query.value.toLowerCase())
        }),
)

function updateDonation() {
    if (selectedParticipant.value) {
        updatedDonation.value.participantId = selectedParticipant.value.id
    }
    newDonationSchema.validate(updatedDonation.value)
        .then(() => {
            donationStore.updateDonation(
                selectedRaffle.value!.id,
                selectedPrize.value!.id,
                updatedDonation.value,
                selectedPrize.value!.ticketCost,
            ).catch(e => {
                console.error(e)
                showError("Не вдалося змінити внесок!")
            })
            props.closeModal()
        })
        .catch((e: ValidationError) => {
            errorMsg.value = e.message
        })
}
</script>
