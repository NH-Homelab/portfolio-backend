package models

import (
	"sort"
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`

	Milestones []Milestone
}

// Returns a slice of tagnames (strings) sorted by number of occurrences in descending order
func (p Project) tags() []string {
	tags := make(map[string]int)

	for _, milestone := range p.Milestones {
		for _, tag := range milestone.Tags {
			tags[tag]++
		}
	}

	// Create a slice of tag names and sort by count (descending)
	tagNames := make([]string, 0, len(tags))
	for tag := range tags {
		tagNames = append(tagNames, tag)
	}

	sort.Slice(tagNames, func(i, j int) bool {
		return tags[tagNames[i]] > tags[tagNames[j]]
	})

	return tagNames
}
