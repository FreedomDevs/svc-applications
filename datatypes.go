package main

type ApplicationDecision string

const (
	ApplicationDecisionAccept  ApplicationDecision = "ACCEPT"
	ApplicationDecisionReject  ApplicationDecision = "REJECT"
	ApplicationDecisionPending ApplicationDecision = "PENDING"
)

type ApplicationData struct {
	id             uint64
	userid         string
	age            uint8
	about          string
	join_reason    string
	inviter        string
	submitted_at   string
	ai_categories  *map[string]any
	ai_decision    ApplicationDecision
	admin_decision ApplicationDecision
	ai_answer      string
	ai_comment     string
}

type ApplicationRequest struct {
	age         uint8
	about       string
	join_reason string
	inviter     string
}
