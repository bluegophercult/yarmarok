import { object, string, InferType } from "yup"

export type Raffle = {
    id: string,
    name: string,
    note: string,
}

export type Raffles = Raffle[]

export const newRaffleSchema = object({
    name: string().required("Назва обов'язкова"),
    note: string(),
})

export type NewRaffle = InferType<typeof newRaffleSchema>
