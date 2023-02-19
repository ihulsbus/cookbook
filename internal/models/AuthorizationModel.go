package models

const (
	// ACTIONS
	get    = "GET"
	put    = "PUT"
	post   = "POST"
	delete = "DELETE"

	// ROLES
	member = "members"
	admin  = "administrators"
)

var (
	AuthorizationModel = map[string]map[string][]string{
		"recipe": {
			get:    {member, admin},
			put:    {member, admin},
			post:   {member, admin},
			delete: {admin},
		},
		"ingredient": {
			get:    {member, admin},
			put:    {member, admin},
			post:   {member, admin},
			delete: {admin},
		},
		"tag": {
			get:    {member, admin},
			put:    {member, admin},
			post:   {member, admin},
			delete: {admin},
		},
		"category": {
			get:    {member, admin},
			put:    {member, admin},
			post:   {member, admin},
			delete: {admin},
		},
	}
)
