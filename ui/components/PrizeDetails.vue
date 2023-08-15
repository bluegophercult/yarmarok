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
            <div class="flex justify-between">
                <div class="text-xl">{{ selectedPrize.name }}</div>
                <div>Ціна купону: {{ selectedPrize.ticketCost }}</div>
            </div>
            <div v-if="selectedPrize.description" class="whitespace-pre">{{ selectedPrize.description }}</div>
            <hr>
            <div class="flex justify-between items-center gap-2">
                <div>
                    <div class="text-xl">Внески</div>
                </div>
                <div>
                    <TheButton :click="() => { isOpenCreateDonation = true} ">Додати внесок</TheButton>
                </div>
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

const prizeStore = usePrizeStore()
const { selectedPrize } = storeToRefs(prizeStore)

const isOpenDelete = ref(false)
const isOpenUpdate = ref(false)
const isOpenCreateDonation = ref(false)
</script>
