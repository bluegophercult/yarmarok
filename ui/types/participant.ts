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
        .required("Номер телефону обов'язковий")
        .matches(/^\+380\d{9,10}$/, "Номер телефону повинен бути вигляду: +380112233444"),
    note: string(),
})

export type NewParticipant = InferType<typeof newParticipantSchema>
