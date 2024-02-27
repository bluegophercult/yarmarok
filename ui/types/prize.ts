import { object, string, number, InferType } from "yup"
import { Participant } from "~/types/participant"
import { Donation } from "~/types/donation"

export type Prize = {
    id: string,
    name: string,
    ticketCost: number,
    description: string,
    playResult: PlayResult | null,
}

export type PlayResult = {
    winners: PlayParticipant[]
    participants: PlayParticipant[]
}

export type PlayParticipant = {
    participant: Participant
    totalDonation: number
    totalTicketsNumber: number
    donations: Donation
}

export type Prizes = Prize[]

export const newPrizeSchema = object({
    name: string()
        .required("Назва обов'язкова"),
    ticketCost: number()
        .when("name", { is: () => true, then: (schema) => schema.required() })
        .required("Ціна купону обов'язкова")
        .typeError("Ціна купону повинна бути числом")
        .integer("Ціна купону повинна бути цілим числом")
        .positive("Ціна купону повинна бути більша нуля"),
    description: string(),
})

export type NewPrize = InferType<typeof newPrizeSchema>
