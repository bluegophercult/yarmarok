export type Raffle = {
    id: string,
    name: string,
    note: string,
}

export type Raffles = Raffle[]

export type NewRaffle = {
    name: string,
    note: string,
}
