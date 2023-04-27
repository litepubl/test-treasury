// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Person -.
type Person struct {
	Uid       int    `json:"uid"       example:2674`
	FirstName string `json:"first_name"       example:"Abu"`
	LastName  string `json:"last_name"       example:"ABBAS"`
}
