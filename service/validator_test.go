package service

import (
	"testing"
)

func Test_validateRaffle(t *testing.T) {
	type args struct {
		raf *RaffleRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid RaffleRequest",
			args: args{
				raf: &RaffleRequest{
					Name: "Example Raffle123Ї",
					Note: "ExampleЇ",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid RaffleRequest (Name too short)",
			args: args{
				raf: &RaffleRequest{
					Name: "Ra",
					Note: "Example",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid RaffleRequest (Invalid symbols)",
			args: args{
				raf: &RaffleRequest{
					Name: "Example RaffleЇ世",
					Note: "Example世",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.raf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validateRaffle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePrize(t *testing.T) {
	type args struct {
		p *PrizeRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid PrizeRequest",
			args: args{
				p: &PrizeRequest{
					Name:        "Example PrizeЇ",
					TicketCost:  100,
					Description: "Example",
				},
			},
			wantErr: false,
		}, {
			name: "Invalid PrizeRequest (Negative Value)",
			args: args{
				p: &PrizeRequest{
					Name:        "Invalid Prize",
					TicketCost:  -50,
					Description: "Example",
				},
			},
			wantErr: true,
		}, {
			name: "Invalid PrizeRequest (Name too short)",
			args: args{
				p: &PrizeRequest{
					Name:        "Ra",
					TicketCost:  100,
					Description: "Example",
				},
			},
			wantErr: true,
		}, {
			name: "Invalid PrizeRequest (Invalid symbols)",
			args: args{
				p: &PrizeRequest{
					Name:        "Example PrizeЇ世",
					TicketCost:  100,
					Description: "Example",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validatePrize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateParticipant(t *testing.T) {
	type args struct {
		p *ParticipantRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid ParticipantRequest",
			args: args{
				p: &ParticipantRequest{
					Name:  "John DoeЇ",
					Phone: "+380123456789",
					Note:  "Example",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid ParticipantRequest (Invalid Phone)",
			args: args{
				p: &ParticipantRequest{
					Name:  "John Doe",
					Phone: "invalidphone",
					Note:  "Example",
				},
			},
			wantErr: true,
		}, {
			name: "Invalid ParticipantRequest (Name too short)",
			args: args{
				p: &ParticipantRequest{
					Name:  "J",
					Phone: "+380123456789",
					Note:  "Example",
				},
			},
			wantErr: true,
		}, {
			name: "Invalid ParticipantRequest (Invalid symbols)",
			args: args{
				p: &ParticipantRequest{
					Name:  "John DoeЇ世",
					Phone: "+380123456789",
					Note:  "Example世",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParticipant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
