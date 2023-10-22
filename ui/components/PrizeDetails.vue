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
                <TheButton class="mb-2" :click="playPrize" :disabled="donations.length == 0">
                    {{ selectedPrize.playResults != null ? "Розіграти знову" : "Розіграти приз" }}
                </TheButton>

                <div v-if="selectedPrize.playResults != null">
                    Розіграші:
                    <div class="divide-y">
                        <div v-for="(result, resultIdx) in selectedPrize.playResults" :key="resultIdx">
                            <ul>
                                <li v-for="(winner, winnerIdx) in result.winners" :key="winnerIdx">
                                    {{ winner.participant.name }}
                                    <span v-if="winner.participant.phone">({{ winner.participant.phone }})</span>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>

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
</template>

<script setup lang="ts">
import { usePrizeStore } from "~/store/prize"
import { useDonationStore } from "~/store/donation"
import { useRaffleStore } from "~/store/raffle"

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
</script>
