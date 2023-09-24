import { object, string, number, InferType } from "yup"

export type Donation = {
    id: string,
    participantId: string,
    amount: number,
    ticketsNumber: number,
    createdAt: string,
}

export type Donations = Donation[]

export const newDonationSchema = object({
    amount: number()
        .required("Сума внесоку обов'язкова")
        .typeError("Сума внесоку повинна бути числом")
        .integer("Сума внесоку повинна бути цілим числом")
        .positive("Сума внесоку повинна бути більша нуля"),
    participantId: string()
        .when("amount", { is: () => true, then: (schema) => schema.required() })
        .required("Учасник обов'язковий"),
})

export type NewDonation = InferType<typeof newDonationSchema>
