package request

type UpdateTeamRequest struct {
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
}
