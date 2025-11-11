package clientlibrary

type Services struct {
	Users       *UsersService
	VhostLimits *VhostLimitsService
	Vhosts      *VhostsService
	Queues      *QueuesService
	Policies    *PoliciesService
	Exchanges   *ExchangesService
	Permissions *PermissionsService
	Parameters  *ParametersService
	Bindings    *BindingsService
	Messages    *MessagesService
}

func NewServices(client *Client) *Services {
	return &Services{
		Users:       &UsersService{service: service{client: client, name: "users"}},
		VhostLimits: &VhostLimitsService{service: service{client: client, name: "vhostlimits"}},
		Vhosts:      &VhostsService{service: service{client: client, name: "vhosts"}},
		Queues:      &QueuesService{service: service{client: client, name: "queues"}},
		Policies:    &PoliciesService{service: service{client: client, name: "policies"}},
		Exchanges:   &ExchangesService{service: service{client: client, name: "exchanges"}},
		Permissions: &PermissionsService{service: service{client: client, name: "permissions"}},
		Parameters:  &ParametersService{service: service{client: client, name: "parameters"}},
		Bindings:    &BindingsService{service: service{client: client, name: "bindings"}},
		Messages:    &MessagesService{service: service{client: client, name: "messages"}},
	}
}
