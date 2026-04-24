package repository

import (
	"context"
	"fmt"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"log"
)

type TemplateKPIRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.TemplateKPIResponse, int, error)
	GetPublishedByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error)
	Create(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	Delete(ctx context.Context, id string) error
}

type templateKPIRepository struct {
	client *db.PrismaClient
}

func NewTemplateKPIRepository(client *db.PrismaClient) TemplateKPIRepository {
	return &templateKPIRepository{client: client}
}

func (r *templateKPIRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.TemplateKPIResponse, int, error) {
	skip := (params.Page - 1) * params.Limit

	// Build filter
	var filters []db.TemplateKPIWhereParam
	if params.Search != "" {
		filters = append(filters, db.TemplateKPI.Or(
			db.TemplateKPI.Department.Contains(params.Search),
			db.TemplateKPI.TemplateName.Contains(params.Search),
			db.TemplateKPI.Description.Contains(params.Search),
		))
	}

	// Build sort
	sortDir := db.SortOrderDesc
	if params.SortDir == "asc" {
		sortDir = db.SortOrderAsc
	}
	var orderBy []db.TemplateKPIOrderByParam
	switch params.SortBy {
	case "department":
		orderBy = append(orderBy, db.TemplateKPI.Department.Order(sortDir))
	case "templateName":
		orderBy = append(orderBy, db.TemplateKPI.TemplateName.Order(sortDir))
	case "createdAt":
		orderBy = append(orderBy, db.TemplateKPI.CreatedAt.Order(sortDir))
	default:
		orderBy = append(orderBy, db.TemplateKPI.CreatedAt.Order(sortDir))
	}

	allKPI, err := r.client.TemplateKPI.FindMany(filters...).
		With(db.TemplateKPI.Items.Fetch()).
		OrderBy(orderBy...).
		Skip(skip).
		Take(params.Limit).
		Exec(ctx)

	if err != nil {
		log.Printf("[TemplateKPIRepository] Error in GetAll (FindMany): %v", err)
		return nil, 0, err
	}

	// Count total
	countRes, err := r.client.TemplateKPI.FindMany(filters...).Exec(ctx)
	if err != nil {
		log.Printf("[TemplateKPIRepository] Error in GetAll (Count): %v", err)
		return nil, 0, err
	}
	total := len(countRes)

	responses := make([]dto.TemplateKPIResponse, 0, len(allKPI))
	for _, du := range allKPI {
		items := make([]dto.KPIItemDTO, 0, len(du.Items()))
		for _, it := range du.Items() {
			items = append(items, dto.KPIItemDTO{
				ID:         it.ID,
				TemplateID: it.TemplateID,
				NameResult: it.NameResult,
				KpiResult:  it.KpiResult,
				Weight:     it.Weight,
				Target:     it.Target,
				Actual:     it.Actual,
				Score:      it.Score,
				FinalScore: it.FinalScore,
			})
		}

		responses = append(responses, dto.TemplateKPIResponse{
			InnerTemplateKPI: du.InnerTemplateKPI,
			Items:            items,
		})
	}

	return responses, total, nil
}

func (r *templateKPIRepository) GetPublishedByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error) {
	// Use case-insensitive matching for department
	var rawTemplates []struct {
		ID           string `json:"id"`
	}
	
	err := r.client.Prisma.QueryRaw(`
		SELECT id FROM "TemplateKPI" 
		WHERE LOWER(TRIM(department)) = LOWER(TRIM($1)) 
		AND "isPublished" = true
	`, department).Exec(ctx, &rawTemplates)

	if err != nil {
		log.Printf("[TemplateKPIRepository] Error fetching published KPI IDs: %v", err)
		return nil, err
	}

	if len(rawTemplates) == 0 {
		return []dto.TemplateKPIResponse{}, nil
	}

	ids := make([]string, len(rawTemplates))
	for i, t := range rawTemplates {
		ids[i] = t.ID
	}

	allKPI, err := r.client.TemplateKPI.FindMany(
		db.TemplateKPI.ID.In(ids),
	).With(db.TemplateKPI.Items.Fetch()).Exec(ctx)

	if err != nil {
		log.Printf("[TemplateKPIRepository] Error fetching full published KPIs: %v", err)
		return nil, err
	}

	responses := make([]dto.TemplateKPIResponse, 0, len(allKPI))
	for _, du := range allKPI {
		items := make([]dto.KPIItemDTO, 0, len(du.Items()))
		for _, it := range du.Items() {
			items = append(items, dto.KPIItemDTO{
				ID:         it.ID,
				TemplateID: it.TemplateID,
				NameResult: it.NameResult,
				KpiResult:  it.KpiResult,
				Weight:     it.Weight,
				Target:     it.Target,
				Actual:     it.Actual,
				Score:      it.Score,
				FinalScore: it.FinalScore,
			})
		}

		responses = append(responses, dto.TemplateKPIResponse{
			InnerTemplateKPI: du.InnerTemplateKPI,
			Items:            items,
		})
	}

	return responses, nil
}

func (r *templateKPIRepository) Create(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	// 1. Create Template
	du, err := r.client.TemplateKPI.CreateOne(
		db.TemplateKPI.Email.Set(req.Email),
		db.TemplateKPI.Department.Set(req.Department),
		db.TemplateKPI.TemplateName.Set(req.TemplateName),
		db.TemplateKPI.Description.Set(req.Description),
		db.TemplateKPI.Attachment.Set(req.Attachment),
		db.TemplateKPI.IsPublished.Set(req.IsPublished),
	).Exec(ctx)

	if err != nil {
		log.Printf("[TemplateKPIRepository] Error creating TemplateKPI: %v", err)
		return nil, err
	}

	// 2. Create Items if any
	for _, it := range req.Items {
		_, err := r.client.KPIItem.CreateOne(
			db.KPIItem.NameResult.Set(it.NameResult),
			db.KPIItem.KpiResult.Set(it.KpiResult),
			db.KPIItem.Weight.Set(it.Weight),
			db.KPIItem.Target.Set(it.Target),
			db.KPIItem.Actual.Set(it.Actual),
			db.KPIItem.Score.Set(it.Score),
			db.KPIItem.FinalScore.Set(it.FinalScore),
			db.KPIItem.Template.Link(db.TemplateKPI.ID.Equals(du.ID)),
		).Exec(ctx)
		if err != nil {
			log.Printf("[TemplateKPIRepository] Error creating KPIItem for template %s: %v", du.ID, err)
		}
	}

	// 3. Fetch full template with items
	full, err := r.client.TemplateKPI.FindUnique(db.TemplateKPI.ID.Equals(du.ID)).With(db.TemplateKPI.Items.Fetch()).Exec(ctx)
	if err != nil {
		return &dto.TemplateKPIResponse{InnerTemplateKPI: du.InnerTemplateKPI}, nil
	}

	items := make([]dto.KPIItemDTO, 0, len(full.Items()))
	for _, it := range full.Items() {
		items = append(items, dto.KPIItemDTO{
			ID:         it.ID,
			TemplateID: it.TemplateID,
			NameResult: it.NameResult,
			KpiResult:  it.KpiResult,
			Weight:     it.Weight,
			Target:     it.Target,
			Actual:     it.Actual,
			Score:      it.Score,
			FinalScore: it.FinalScore,
		})
	}

	return &dto.TemplateKPIResponse{
		InnerTemplateKPI: full.InnerTemplateKPI,
		Items:            items,
	}, nil
}

func (r *templateKPIRepository) Update(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	// 1. Fetch current data to preserve fields
	current, err := r.client.TemplateKPI.FindUnique(db.TemplateKPI.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	// 2. Prepare update data (merge logic)
	department := req.Department
	if department == "" {
		department = current.Department
	}
	templateName := req.TemplateName
	if templateName == "" {
		templateName = current.TemplateName
	}
	description := req.Description
	if description == "" {
		description = current.Description
	}
	
	isPublished := current.IsPublished
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	updateParams := []db.TemplateKPISetParam{
		db.TemplateKPI.Department.Set(department),
		db.TemplateKPI.TemplateName.Set(templateName),
		db.TemplateKPI.Description.Set(description),
		db.TemplateKPI.IsPublished.Set(isPublished),
	}
	
	if req.Attachment != nil {
		updateParams = append(updateParams, db.TemplateKPI.Attachment.SetOptional(req.Attachment))
	} else {
		updateParams = append(updateParams, db.TemplateKPI.Attachment.SetOptional(current.InnerTemplateKPI.Attachment))
	}

	// 3. Prepare Transactions
	var ops []db.PrismaTransaction

	// Update Template op
	ops = append(ops, r.client.TemplateKPI.FindUnique(
		db.TemplateKPI.ID.Equals(id),
	).Update(updateParams...).Tx())

	// Update Items if provided
	if len(req.Items) > 0 {
		// Delete existing op
		ops = append(ops, r.client.KPIItem.FindMany(db.KPIItem.TemplateID.Equals(id)).Delete().Tx())

		// Create new ops
		for _, it := range req.Items {
			ops = append(ops, r.client.KPIItem.CreateOne(
				db.KPIItem.NameResult.Set(it.NameResult),
				db.KPIItem.KpiResult.Set(it.KpiResult),
				db.KPIItem.Weight.Set(it.Weight),
				db.KPIItem.Target.Set(it.Target),
				db.KPIItem.Actual.Set(it.Actual),
				db.KPIItem.Score.Set(it.Score),
				db.KPIItem.FinalScore.Set(it.FinalScore),
				db.KPIItem.Template.Link(db.TemplateKPI.ID.Equals(id)),
			).Tx())
		}
	}

	// Execute Transaction
	if err := r.client.Prisma.Transaction(ops...).Exec(ctx); err != nil {
		log.Printf("[TemplateKPIRepository] Update transaction failed for %s: %v", id, err)
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	// 4. Fetch full template with items for response
	full, err := r.client.TemplateKPI.FindUnique(db.TemplateKPI.ID.Equals(id)).With(db.TemplateKPI.Items.Fetch()).Exec(ctx)
	if err != nil {
		// Even if fetch fails, the update happened. But we should try to return something.
		return nil, fmt.Errorf("transaction succeeded but failed to fetch updated data: %w", err)
	}

	items := make([]dto.KPIItemDTO, 0, len(full.Items()))
	for _, it := range full.Items() {
		items = append(items, dto.KPIItemDTO{
			ID:         it.ID,
			TemplateID: it.TemplateID,
			NameResult: it.NameResult,
			KpiResult:  it.KpiResult,
			Weight:     it.Weight,
			Target:     it.Target,
			Actual:     it.Actual,
			Score:      it.Score,
			FinalScore: it.FinalScore,
		})
	}

	return &dto.TemplateKPIResponse{
		InnerTemplateKPI: full.InnerTemplateKPI,
		Items:            items,
	}, nil
}

func (r *templateKPIRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.TemplateKPI.FindUnique(
		db.TemplateKPI.ID.Equals(id),
	).Delete().Exec(ctx)
	if err != nil {
		log.Printf("[TemplateKPIRepository] Error deleting TemplateKPI %s: %v", id, err)
	}
	return err
}
