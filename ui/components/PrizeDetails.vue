<template>
    <div class="rounded-lg bg-white p-3 shadow-md ring-1 ring-black ring-opacity-5">
        <div class="flex justify-between items-center gap-2">
            <div class="text-xl">Деталі призу</div>
            <div class="flex gap-2">
                <TheButton :click="() => { isOpenUpdate = true} " :disabled="!selectedPrize">Змінити</TheButton>
                <TheButton :click="() => { isOpenDelete = true }" danger :disabled="!selectedPrize">Видалити</TheButton>
            </div>
        </div>
        <hr class="mt-2">
        <div v-if="selectedPrize" class="flex flex-col gap-2 mt-2">
            <div class="flex justify-between items-center">
                <div class="text-xl">{{ selectedPrize.name }}</div>
                <div>Ціна купону: {{ selectedPrize.ticketCost }} грн</div>
            </div>
            <div v-if="selectedPrize.description" class="whitespace-pre">{{ selectedPrize.description }}</div>
            <hr>
            <div class="flex justify-between items-center gap-2">
                <div>
                    <div class="text-xl">Внески</div>
                </div>
                <div>
                    <DonationsCreate/>
                </div>
            </div>
            <div>
                <div v-if="selectedPrize.playResults != null" class="mb-2">
                    <span class="mr-2">Переможці:</span>
                    <span v-for="(result, resultIdx) in selectedPrize.playResults" :key="resultIdx">
                        <span v-for="(winner, winnerIdx) in result.winners" :key="winnerIdx"
                              @click="openParticipantView(winner.participant)"
                              class="rounded-md bg-gray-100 px-1 shadow ring-1 ring-gray-600 ring-opacity-5 mr-2 cursor-pointer">
                            {{ winner.participant.name }}
                        </span>
                        <span v-if="selectedPrize.playResults!.length - 1 != resultIdx" class="mr-2">|</span>
                    </span>
                </div>

                <TheButton class="mb-1" :click="playPrize" :disabled="donations.length == 0">
                    {{ selectedPrize.playResults != null ? "Розіграти знову" : "Розіграти приз" }}
                </TheButton>

                <p>Загальна кількість купонів: {{ totalTickets }}</p>
            </div>
            <DonationsList/>
        </div>
        <div v-else class="mt-2 text-gray-400">
            Не вибрано приз
        </div>
    </div>

    <PrizesUpdate v-if="selectedPrize" :prize="selectedPrize" :is-open="isOpenUpdate"
                  :close-modal="() => isOpenUpdate = false"/>
    <PrizesDelete v-if="selectedPrize" :prize="selectedPrize" :is-open="isOpenDelete"
                  :close-modal="() => isOpenDelete = false"/>

    <ParticipantsView v-if="selectedParticipant" :participant="selectedParticipant" :is-open="isOpenParticipantView"
                      :close-modal="() => isOpenParticipantView = false" hide-controls/>
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { useDonationStore } from "~/store/donation"
import { useRaffleStore } from "~/store/raffle"
import { Ref } from "@vue/reactivity"
import { Participant } from "~/types/participant"

const raffleStore = useRaffleStore()
const { selectedRaffle } = storeToRefs(raffleStore)

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const donationStore = useDonationStore()
const { donations } = storeToRefs(donationStore)

const totalTickets = computed(() => donations.value.reduce((acc, d) => {
    return acc + d.ticketsNumber
}, 0))

async function playPrize() {
    await prizeStore.playPrize(selectedRaffle.value!.id, selectedPrize.value!.id)
    await prizeStore.getPrizes(selectedRaffle.value!.id)
}

const isOpenDelete = ref(false)
const isOpenUpdate = ref(false)

const isOpenParticipantView = ref(false)
const selectedParticipant: Ref<Participant | undefined> = ref(undefined)

function openParticipantView(participant: Participant) {
    selectedParticipant.value = participant
    isOpenParticipantView.value = true
}
</script>
