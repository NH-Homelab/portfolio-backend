package models

import "time"

// Enum for the different types of possible milestones
type Milestone_Type string
const (
	Major     Milestone_Type = "project_major"	// Major implies it's a project
	Minor     Milestone_Type = "project_minor"	// Minor implies it's a project
	Education Milestone_Type = "education"
	Career    Milestone_Type = "career"
)

type Milestone struct {
	ID             int            `json:"id"`
	Title          string         `json:"title"`
	Milestone_date time.Time      `json:"milestone_date"`
	Description    string         `json:"description"`
	Body_url       string         `json:"body_url"`
	Github_url     string         `json:"github_url"`
	Image_url      string         `json:"image_url"`
	Milestone_type Milestone_Type `json:"milestone_type"`
	Status         string         `json:"status"`
	Project_id     int            `json:"project_id"`
	
	Tags           []string
}

