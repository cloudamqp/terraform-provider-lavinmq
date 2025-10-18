package clientlibrary

type Services struct {
	// TODO: not needed
	client *Client

	Users       *UsersService
	VhostLimits *VhostLimitsService
	Vhosts      *VhostsService
	Queues      *QueuesService
	Policies    *PoliciesService
	Exchanges   *ExchangesService
	Permissions *PermissionsService
	Bindings    *BindingsService
}

func NewServices(client *Client) *Services {
	return &Services{
		client:      client,
		Users:       (*UsersService)(&service{client: client}),
		VhostLimits: (*VhostLimitsService)(&service{client: client}),
		Vhosts:      (*VhostsService)(&service{client: client}),
		Queues:      (*QueuesService)(&service{client: client}),
		Policies:    (*PoliciesService)(&service{client: client}),
		Exchanges:   (*ExchangesService)(&service{client: client}),
		Permissions: (*PermissionsService)(&service{client: client}),
		Bindings:    (*BindingsService)(&service{client: client}),
	}
}
