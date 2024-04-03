package model

type VoteDTO struct {
	ParticipantId string `json:"participantId"`
	BarId         string `json:"barId"`
}

type VoteStatsDTO = []VoteStatsRow

type VoteStatsRow struct {
	BarId     uint   `json:"bar_id"`
	BarName   string `json:"bar_name"`
	VoteCount uint   `json:"vote_count"`
}
