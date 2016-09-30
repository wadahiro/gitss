package model

type Search struct {
	ID            string `json:"_id"`
	Rev           string `json:"_rev"`
	Name          string `json:"name"`
	SyncSettingID string `json:"syncSettingId"`
	IssueID       string `json:"issueId"`
	Summary       string `json:"summary"`
	Description   string `json:"description"`
	Created       string `json:"created"`
	Updated       string `json:"updated"`
	Checked       bool   `json:"checked"`
	CheckedDate   string `json:"checkedDate"`
	Memo          string `json:"memo"`
}
