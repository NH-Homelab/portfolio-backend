package publichandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/NH-Homelab/portfolio-backend/internal/models"
	portfoliodao "github.com/NH-Homelab/portfolio-backend/internal/portfolio_dao"
)

type PublicHandler struct {
	dao *portfoliodao.PortfolioDao
}

func NewPublicHandler(dao *portfoliodao.PortfolioDao) *PublicHandler {
	return &PublicHandler{dao}
}

func (ph *PublicHandler) RegisterHandlers(mux *http.ServeMux) {
	// retrieves the project by id
	mux.HandleFunc("GET /api/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		project, err := ph.dao.GetProjectById(id)
		if err != nil {
			log.Printf("Failed to retrieve project %d: %v", id, err)
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}

		// Filter to only include published milestones
		publishedMilestones := make([]models.Milestone, 0)
		for _, m := range project.Milestones {
			if m.Status == "published" {
				publishedMilestones = append(publishedMilestones, m)
			}
		}
		project.Milestones = publishedMilestones

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(project); err != nil {
			log.Printf("Failed to encode project response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// retrieves all published projects
	mux.HandleFunc("GET /api/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, err := ph.dao.GetAllProjects()
		if err != nil {
			log.Printf("Failed to retrieve projects: %v", err)
			http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			log.Printf("Failed to encode projects response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// retrieves the milestone by id if the milestone is published
	mux.HandleFunc("GET /api/milestones/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid milestone ID", http.StatusBadRequest)
			return
		}

		milestone, err := ph.dao.GetMilestoneById(id)
		if err != nil {
			log.Printf("Failed to retrieve milestone %d: %v", id, err)
			http.Error(w, "Milestone not found", http.StatusNotFound)
			return
		}

		// Only return if published
		if milestone.Status != "published" {
			http.Error(w, "Milestone not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(milestone); err != nil {
			log.Printf("Failed to encode milestone response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// retrieves all published milestones
	mux.HandleFunc("GET /api/milestones", func(w http.ResponseWriter, r *http.Request) {
		milestones, err := ph.dao.GetAllPublishedMilestones()
		if err != nil {
			log.Printf("Failed to retrieve milestones: %v", err)
			http.Error(w, "Failed to retrieve milestones", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(milestones); err != nil {
			log.Printf("Failed to encode milestones response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}
