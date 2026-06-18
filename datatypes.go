package main

type ApplicationDecision string

const (
	ApplicationDecisionAccept  ApplicationDecision = "ACCEPT"
	ApplicationDecisionReject  ApplicationDecision = "REJECT"
	ApplicationDecisionPending ApplicationDecision = "PENDING"
)

type ApplicationData struct {
	id             string
	user_id        string
	user_ip        string
	age            uint16
	about          string
	inviter        string
	ai_categories  *map[string]any
	ai_decision    ApplicationDecision
	ai_answer      string
	ai_comment     string
	admin_decision ApplicationDecision
	admin_id       string
}

type ApplicationDataShort struct {
	id             string
	user_id        string
	age            uint16
	about          string
	inviter        string
	ai_categories  *map[string]any
	ai_decision    ApplicationDecision
	ai_answer      string
	admin_decision ApplicationDecision
	admin_id       string
}

type ApplicationRequest struct {
	age     uint16
	about   string
	inviter string
}
