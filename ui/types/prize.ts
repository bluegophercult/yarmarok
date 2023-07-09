import { object, string, number, InferType } from "yup"

export type Prize = {
    id: string,
    name: string,
    ticketCost: number,
    description: string,
}

export type Prizes = Prize[]

export const newPrizeSchema = object({
    name: string()
        .required("Назва обов'язкова"),
    ticketCost: number()
        .required("Ціна купону обов'язкова")
        .typeError("Ціна купону повинна бути числом")
        .integer("Ціна купону повинна бути цілим числом")
        .positive("Ціна купону повинна бути більша нуля"),
    note: string(),
})

export type NewPrize = InferType<typeof newPrizeSchema>