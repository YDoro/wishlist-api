package postgresDB

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/ydoro/wishlist/internal/domain"
)

type wishlistRepo struct {
	DB *sql.DB
}

func NewWishlistRepository(db *sql.DB) *wishlistRepo {
	return &wishlistRepo{
		DB: db,
	}
}

func (r *wishlistRepo) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `INSERT INTO wishlists (id, customer_id, title, items) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, wishlist.ID, wishlist.CustomerId, wishlist.Title, pq.Array(wishlist.Items))
	return err
}

func (r *wishlistRepo) GetById(ctx context.Context, wishlistId string) (*domain.Wishlist, error) {
	query := `SELECT id, customer_id, title, items FROM wishlists WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, wishlistId)

	wishlist := &domain.Wishlist{}
	err := row.Scan(&wishlist.ID, &wishlist.CustomerId, &wishlist.Title, &wishlist.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return wishlist, nil
}

func (r *wishlistRepo) GetByTitle(ctx context.Context, customerId string, title string) (*domain.Wishlist, error) {
	query := `SELECT id, customer_id, title, items FROM wishlists WHERE customer_id = $1 AND title = $2`
	row := r.DB.QueryRowContext(ctx, query, customerId, title)

	wishlist := &domain.Wishlist{}
	err := row.Scan(&wishlist.ID, &wishlist.CustomerId, &wishlist.Title, &wishlist.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return wishlist, nil
}

func (r *wishlistRepo) GetByCustomerId(ctx context.Context, customerId string) ([]*domain.Wishlist, error) {
	query := `SELECT id, customer_id, title, items FROM wishlists WHERE customer_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []*domain.Wishlist
	for rows.Next() {
		wishlist := &domain.Wishlist{}
		err := rows.Scan(&wishlist.ID, &wishlist.CustomerId, &wishlist.Title, &wishlist.Items)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}

	return wishlists, nil
}

func (r *wishlistRepo) UpdateWishlistName(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `UPDATE wishlists SET title = $1 WHERE id = $2 AND customer_id = $3`
	result, err := r.DB.ExecContext(ctx, query, wishlist.Title, wishlist.ID, wishlist.CustomerId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *wishlistRepo) DeleteWishlist(ctx context.Context, wishlistId string) error {
	query := `DELETE FROM wishlists WHERE id = $1`
	result, err := r.DB.ExecContext(ctx, query, wishlistId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *wishlistRepo) AddItemToWishlist(ctx context.Context, wishlistId string, itemId string) error {
	query := `UPDATE wishlists SET items = array_append(items, $1) WHERE id = $2`
	result, err := r.DB.ExecContext(ctx, query, itemId, wishlistId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *wishlistRepo) RemoveItemFromWishlist(ctx context.Context, wishlistId string, itemId string) error {
	query := `UPDATE wishlists SET items = array_remove(items, $1) WHERE id = $2`
	result, err := r.DB.ExecContext(ctx, query, itemId, wishlistId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
