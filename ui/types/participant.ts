import { object, string, InferType } from "yup"
import { sequence } from "nth-check"

export type Participant = {
    id: string,
    name: string,
    phone: string,
    note: string,
}

export type Participants = Participant[]

export const newParticipantSchema = object({
    name: string()
        .required("Ім'я обов'язкове"),
    phone: string()
        .when("name", { is: () => true, then: (schema) => schema.required() })
        .required("Номер телефону обов'язковий"),
    note: string(),
})

export type NewParticipant = InferType<typeof newParticipantSchema>
