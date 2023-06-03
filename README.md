# Yarmarok

An app to make charity auctions smoother
The project aimed at organizing local street fairs or bazaars to collect funds for charity or military
support.
It provides a platform for managing lots, participants, donations, and raffles.

## Entities and value objects.

### Organizer

- **ID**: Unique identifier of user account.

### Raffle

A raffle represents the overall event and the collection of lots and participants. It has the following properties:

- `ID`: The unique identifier of the raffle.
- `Organizer`: The org of raffle.
- `Name`: The name of the raffle.
- `Description`: Optional description of the raffle.
- `Lots`: The list of lots associated with the raffle.
- `Participants`: The list of participants associated with the raffle.
- `CreatedAt`: The date of raffle creation.

### Lot

- `ID`: Unique identifier for the lot.
- `Name`: Name of the lot.
- `TicketCost`: The cost of a single ticket for the lot.
- `Description`:  Optional description of the lot.
- `Donations`: List of participants who have made contributions for this lot.
- `Winner`: The winner of the raffle for this lot. **Field is not finalized.**

### Participant

- `ID`: Unique identifier for the participant.
- `Name`: Name of the participant.
- `Phone`: Phone number of the participant.
- `Note`: Optional note about the participant.

### Donation

- `ID`: Unique identifier for the donation.
- `LotID`: ID of the lot for which the donation was made.
- `ParticipantID`: ID of the participant who made the donation.
- `Amount`: Amount of money transferred for the donation.

## Running the project

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Task](https://taskfile.dev/#/installation)
- [Docker](https://docs.docker.com/get-docker/)
- [Mockgen](https://github.com/golang/mock#installation)


Quick install of go Task:
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Run tasks

```bash
task test # run tests
task integration-test # run heavy integration tests
task generate # generate mocks
```
