package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// pgUniqueViolation is Postgres's SQLSTATE for a unique constraint
// violation, used to recognize a race on the normalized_name unique
// index as a duplicate rather than a generic DB error.
const pgUniqueViolation = "23505"

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}

type IngredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,
	}
}

func (r IngredientRepository) FindAll(pagination m.PaginationRequest) ([]m.Ingredient, int64, error) {
	var ingredients []m.Ingredient
	var total int64

	if err := r.db.Model(&m.Ingredient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&ingredients).Error; err != nil {
		return nil, 0, err
	}

	return ingredients, total, nil
}

func (r IngredientRepository) FindSingle(ingredient m.Ingredient) (m.Ingredient, error) {

	result := r.db.First(&ingredient)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Ingredient{}, errors.New("not found")
		} else {
			return m.Ingredient{}, result.Error
		}
	}

	return ingredient, nil
}

func (r IngredientRepository) FindByName(name string) (m.Ingredient, error) {
	var ing m.Ingredient
	err := r.db.
		Where("name = ?", name).
		First(&ing).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Not found is signaled by a zero value + nil error, not an error.
			return m.Ingredient{}, nil
		}
		return m.Ingredient{}, err
	}

	return ing, nil
}

// normalize computes name's normalized form the same way Postgres's
// normalized_name generated column does, so the result can be compared
// against that column with a plain `=`.
func (r IngredientRepository) normalize(name string) (string, error) {
	var normalized string
	err := r.db.Raw(
		`SELECT lower(immutable_unaccent(trim(regexp_replace(?, '\s+', ' ', 'g'))))`,
		name,
	).Scan(&normalized).Error
	return normalized, err
}

// FindByNormalizedName looks up an ingredient whose normalized_name
// exactly matches name's normalized form (so "Tomaat " matches an
// existing "tomaat"). Mirrors FindByName's not-found-is-empty-not-error convention.
func (r IngredientRepository) FindByNormalizedName(name string) (m.Ingredient, error) {
	normalized, err := r.normalize(name)
	if err != nil {
		return m.Ingredient{}, err
	}

	var ing m.Ingredient
	err = r.db.Where("normalized_name = ?", normalized).First(&ing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m.Ingredient{}, nil
		}
		return m.Ingredient{}, err
	}

	return ing, nil
}

// FindSimilar returns up to limit ingredients whose normalized_name is
// at least 0.6 similar (pg_trgm) to name's normalized form, most similar first.
func (r IngredientRepository) FindSimilar(name string, limit int) ([]m.Ingredient, error) {
	normalized, err := r.normalize(name)
	if err != nil {
		return nil, err
	}

	var ings []m.Ingredient
	err = r.db.
		Where("similarity(normalized_name, ?) > 0.6", normalized).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "similarity(normalized_name, ?) DESC", Vars: []interface{}{normalized}},
		}).
		Limit(limit).
		Find(&ings).Error
	if err != nil {
		return nil, err
	}

	return ings, nil
}

// Create inserts ingredient. A unique_violation on normalized_name
// (the caller's own FindByNormalizedName check having raced with a
// near-simultaneous create) is reported as a "duplicate" sentinel
// rather than a raw DB error.
func (r IngredientRepository) Create(ingredient m.Ingredient) (m.Ingredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&ingredient).Error; err != nil {
			if isUniqueViolation(err) {
				return errors.New("duplicate")
			}
			return err
		}

		return nil
	}); err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) Update(ingredient m.Ingredient) (m.Ingredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Updates(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) Delete(ingredient m.Ingredient) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
