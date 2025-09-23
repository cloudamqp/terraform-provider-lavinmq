package clientlibrary

type Services struct {
	// TODO: not needed
	client *Client

	Users       *UsersService
	VhostLimits *VhostLimitsService
	Vhosts      *VhostsService
	Queues      *QueuesService
}

func NewServices(client *Client) *Services {
	return &Services{
		client:      client,
		Users:       (*UsersService)(&service{client: client}),
		VhostLimits: (*VhostLimitsService)(&service{client: client}),
		Vhosts:      (*VhostsService)(&service{client: client}),
		Queues:      (*QueuesService)(&service{client: client}),
	}
}
