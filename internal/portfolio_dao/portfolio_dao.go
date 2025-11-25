package portfoliodao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/NH-Homelab/portfolio-backend/internal/database"
	"github.com/NH-Homelab/portfolio-backend/internal/models"
)

const (
	getProjectById = `
		SELECT 
			p.id, p.name, p.description, p.created_at,
			m.id, m.title, m.milestone_date, m.description, 
			m.body_url, m.github_url, m.image_url, 
			m.milestone_type, m.status, m.project_id
		FROM projects p
		LEFT JOIN milestones m ON p.id = m.project_id
		WHERE p.id = $1
		ORDER BY m.milestone_date`
	getAllProjects = `
		SELECT id, name, description, created_at
		FROM projects
		ORDER BY id`
	getMilestoneById = `
		SELECT id, title, milestone_date, description, body_url, 
			   github_url, image_url, milestone_type, status, project_id
		FROM milestones
		WHERE id = $1`
	getAllPublishedMilestones = `
		SELECT id, title, milestone_date, description, body_url, 
			   github_url, image_url, milestone_type, status, project_id
		FROM milestones
		WHERE status = 'published'
		ORDER BY milestone_date DESC`
	createProject = `
		INSERT INTO projects (name, description)
		VALUES ($1, $2)
		RETURNING id`
	createMilestone = `
		INSERT INTO milestones (
			title, milestone_date, description, body_url, 
			github_url, image_url, milestone_type, status, project_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	deleteProject = `
		DELETE FROM projects
		WHERE id = $1`
	deleteMilestone = `
		DELETE FROM milestones
		WHERE id = $1`
)

type PortfolioDao struct {
	db database.Database
}

// Update structs for partial updates
type ProjectUpdate struct {
	Name        *string `db:"name"`
	Description *string `db:"description"`
}

type MilestoneUpdate struct {
	Title         *string                `db:"title"`
	MilestoneDate *time.Time             `db:"milestone_date"`
	Description   *string                `db:"description"`
	BodyURL       *string                `db:"body_url"`
	GithubURL     *string                `db:"github_url"`
	ImageURL      *string                `db:"image_url"`
	MilestoneType *models.Milestone_Type `db:"milestone_type"`
	Status        *string                `db:"status"`
	ProjectID     *int                   `db:"project_id"`
}

// Create new instance of PortfolioDao
func NewPortfolioDao(db database.Database) *PortfolioDao {
	return &PortfolioDao{db}
}

// buildUpdateQuery dynamically builds an UPDATE query from a struct with pointer fields
// Only non-nil pointer fields will be included in the update
func buildUpdateQuery(table string, id int, update interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(update)
	t := v.Type()

	var setClauses []string
	var args []interface{}
	argCount := 1

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the db tag for the column name
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			continue // Skip fields without db tag
		}

		// Check if the pointer is non-nil
		if !field.IsNil() {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", dbTag, argCount))
			args = append(args, field.Elem().Interface())
			argCount++
		}
	}

	if len(setClauses) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		table,
		strings.Join(setClauses, ", "),
		argCount)
	args = append(args, id)

	return query, args, nil
}

// Returns all projects without their milestones
func (dao *PortfolioDao) GetAllProjects() ([]models.Project, error) {
	rows, err := dao.db.Query(getAllProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Created_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		p.Milestones = make([]models.Milestone, 0)
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

// Returns a single project with its milestones by ID
func (dao *PortfolioDao) GetProjectById(id int) (*models.Project, error) {
	projects, err := dao.queryProjectsWithMilestones(getProjectById, id)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project with id %d not found", id)
	}

	return &projects[0], nil
}

// Helper function to query projects with milestones and handle row scanning
func (dao *PortfolioDao) queryProjectsWithMilestones(query string, args ...interface{}) ([]models.Project, error) {
	rows, err := dao.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects with milestones: %w", err)
	}
	defer rows.Close()

	projectsMap := make(map[int]*models.Project)
	var projectOrder []int

	for rows.Next() {
		var p models.Project
		var m models.Milestone
		var milestoneID sql.NullInt64
		var milestoneTitle, milestoneDesc sql.NullString
		var milestoneDate sql.NullTime
		var bodyURL, githubURL, imageURL sql.NullString
		var milestoneType, milestoneStatus sql.NullString
		var milestoneProjectID sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Created_at,
			&milestoneID, &milestoneTitle, &milestoneDate, &milestoneDesc,
			&bodyURL, &githubURL, &imageURL,
			&milestoneType, &milestoneStatus, &milestoneProjectID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get or create project
		project, exists := projectsMap[p.ID]
		if !exists {
			project = &p
			project.Milestones = make([]models.Milestone, 0)
			projectsMap[p.ID] = project
			projectOrder = append(projectOrder, p.ID)
		}

		// Add milestone if it exists
		if milestoneID.Valid {
			m.ID = int(milestoneID.Int64)
			m.Title = milestoneTitle.String
			m.Milestone_date = milestoneDate.Time
			m.Description = milestoneDesc.String
			m.Body_url = bodyURL.String
			m.Github_url = githubURL.String
			m.Image_url = imageURL.String
			m.Milestone_type = models.Milestone_Type(milestoneType.String)
			m.Status = milestoneStatus.String
			m.Project_id = int(milestoneProjectID.Int64)

			project.Milestones = append(project.Milestones, m)
		}
	}

	// Convert map to slice in original order
	projects := make([]models.Project, 0, len(projectsMap))
	for _, id := range projectOrder {
		projects = append(projects, *projectsMap[id])
	}

	return projects, rows.Err()
}

// UpdateProject performs a partial update on a project
func (dao *PortfolioDao) UpdateProject(id int, update ProjectUpdate) error {
	query, args, err := buildUpdateQuery("projects", id, update)
	if err != nil {
		return err
	}

	result, err := dao.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", id)
	}

	return nil
}

// UpdateMilestone performs a partial update on a milestone
func (dao *PortfolioDao) UpdateMilestone(id int, update MilestoneUpdate) error {
	query, args, err := buildUpdateQuery("milestones", id, update)
	if err != nil {
		return err
	}

	result, err := dao.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone with id %d not found", id)
	}

	return nil
}

// CreateProject creates a new project and returns its ID
func (dao *PortfolioDao) CreateProject(name, description string) (int, error) {
	rows, err := dao.db.Query(createProject, name, description)
	if err != nil {
		return 0, fmt.Errorf("failed to create project: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("failed to get created project id")
	}

	var id int
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan project id: %w", err)
	}

	return id, nil
}

// CreateMilestone creates a new milestone and returns its ID
func (dao *PortfolioDao) CreateMilestone(m models.Milestone) (int, error) {
	rows, err := dao.db.Query(
		createMilestone,
		m.Title,
		m.Milestone_date,
		m.Description,
		m.Body_url,
		m.Github_url,
		m.Image_url,
		m.Milestone_type,
		m.Status,
		m.Project_id,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create milestone: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("failed to get created milestone id")
	}

	var id int
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan milestone id: %w", err)
	}

	return id, nil
}

// DeleteProject deletes a project by ID (cascades to milestones)
func (dao *PortfolioDao) DeleteProject(id int) error {
	result, err := dao.db.Exec(deleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", id)
	}

	return nil
}

// DeleteMilestone deletes a milestone by ID
func (dao *PortfolioDao) DeleteMilestone(id int) error {
	result, err := dao.db.Exec(deleteMilestone, id)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("milestone with id %d not found", id)
	}

	return nil
}

// GetMilestoneById returns a single milestone by ID
func (dao *PortfolioDao) GetMilestoneById(id int) (*models.Milestone, error) {
	rows, err := dao.db.Query(getMilestoneById, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query milestone: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("milestone with id %d not found", id)
	}

	var m models.Milestone
	var bodyURL, githubURL, imageURL sql.NullString
	var projectID sql.NullInt64

	err = rows.Scan(
		&m.ID, &m.Title, &m.Milestone_date, &m.Description,
		&bodyURL, &githubURL, &imageURL,
		&m.Milestone_type, &m.Status, &projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan milestone row: %w", err)
	}

	// Convert NullString to regular string
	m.Body_url = bodyURL.String
	m.Github_url = githubURL.String
	m.Image_url = imageURL.String
	if projectID.Valid {
		m.Project_id = int(projectID.Int64)
	}
	m.Tags = make([]string, 0)

	return &m, rows.Err()
}

// GetAllPublishedMilestones returns all milestones with 'published' status
func (dao *PortfolioDao) GetAllPublishedMilestones() ([]models.Milestone, error) {
	rows, err := dao.db.Query(getAllPublishedMilestones)
	if err != nil {
		return nil, fmt.Errorf("failed to query published milestones: %w", err)
	}
	defer rows.Close()

	var milestones []models.Milestone
	for rows.Next() {
		var m models.Milestone
		var bodyURL, githubURL, imageURL sql.NullString
		var projectID sql.NullInt64

		err := rows.Scan(
			&m.ID, &m.Title, &m.Milestone_date, &m.Description,
			&bodyURL, &githubURL, &imageURL,
			&m.Milestone_type, &m.Status, &projectID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone row: %w", err)
		}

		// Convert NullString to regular string
		m.Body_url = bodyURL.String
		m.Github_url = githubURL.String
		m.Image_url = imageURL.String
		if projectID.Valid {
			m.Project_id = int(projectID.Int64)
		}
		m.Tags = make([]string, 0)

		milestones = append(milestones, m)
	}

	return milestones, rows.Err()
}
