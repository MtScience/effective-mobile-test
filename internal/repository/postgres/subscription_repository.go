package postgres

import (
	"context"
	"errors"
	"fmt"
	"subscriptions/internal/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "subscriptions/internal/domain/subscription"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	logging.Logger.Info(
		fmt.Sprintf(
			"creating database entry for new subscription with user_id=%s and service=%s",
			sub.UserID,
			sub.Service,
		),
	)

	const query = `
		INSERT INTO subscriptions (user_id, service, price, subscribed_on, unsubscribed_on)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, service, price, subscribed_on, unsubscribed_on
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		sub.UserID,
		sub.Service,
		sub.Price,
		sub.SubscribedOn,
		sub.UnsubscribedOn,
	)
	created, err := scanSubscription(row)
	if err != nil {
		logging.Logger.Error("creating database entry failed", "err", err)
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	logging.Logger.Info(
		fmt.Sprintf("retrieving database entry with id=%d", id),
	)

	const query = `
		SELECT id, user_id, service, price, subscribed_on, unsubscribed_on
		FROM subscriptions
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logging.Logger.Error("no entries found")
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	logging.Logger.Info(
		fmt.Sprintf("updating database entry with id=%d", sub.ID),
	)

	const query = `
		UPDATE subscriptions
		SET user_id = $1
		    service = $2
			price = $3
			subscribed_on = $4
			unsubscribed_on = $5
		WHERE id = $6
		RETURNING id, user_id, name, price, subscribed_on, unsubscribed_on
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		sub.UserID,
		sub.Service,
		sub.Price,
		sub.SubscribedOn,
		sub.UnsubscribedOn,
		sub.ID,
	)
	updated, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logging.Logger.Error("no entries found. Update failed")
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	logging.Logger.Info(
		fmt.Sprintf("deleting database entry with id=%d", id),
	)

	const query = `DELETE subscriptions WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		logging.Logger.Error("no entries found. Nothing to delete")
		return domain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Subscription, error) {
	logging.Logger.Info("retrieving all database entries")

	const query = `
		SELECT id, user_id, service, price, subscribed_on, unsubscribed_on
		FROM subscriptions
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]domain.Subscription, 0)
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			logging.Logger.Error("couldn't extract data for entry", "err", err)
			return nil, err
		}

		subs = append(subs, *sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *Repository) TotalPrice(ctx context.Context, request *domain.PriceRequest) (int64, error) {
	logging.Logger.Info("calculating total subscription price for given time interval")

	const queryTemplate = `
		SELECT sum(price * (extract(YEAR FROM sub_duration) * 12 + extract(MONTH FROM sub_duration))) AS total
		FROM (
		    SELECT
		        price,
		        age(least(unsubscribed_on, $2), greatest(subscribed_on, $1)) AS sub_duration
		    FROM (
		        SELECT
		            price,
		            subscribed_on,
		            coalesce(unsubscribed_on, CURRENT_DATE) AS unsubscribed_on
		        FROM subscriptions
		        WHERE (%s)
		    )
		    WHERE unsubscribed_on >= $1 AND subscribed_on <= $2
		)
	`

	var row pgx.Row
	switch {
	case request.UserID == nil:
		query := fmt.Sprintf(queryTemplate, "service = $3")
		row = r.pool.QueryRow(ctx, query, request.StartDate, request.EndDate, request.Service)
	case request.Service == nil:
		query := fmt.Sprintf(queryTemplate, "user_id = $3")
		row = r.pool.QueryRow(ctx, query, request.StartDate, request.EndDate, request.UserID)
	case request.UserID != nil && request.Service != nil:
		query := fmt.Sprintf(queryTemplate, "user_id = $3 AND service = $4")
		row = r.pool.QueryRow(ctx, query, request.StartDate, request.EndDate, request.UserID, request.Service)
	}

	var totalPrice int64
	if err := row.Scan(&totalPrice); err != nil {
		logging.Logger.Error("error during price calculation", "err", err)
		return -1, domain.ErrNotFound
	}

	return totalPrice, nil
}

type subscriptionScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner subscriptionScanner) (*domain.Subscription, error) {
	var sub domain.Subscription

	if err := scanner.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Service,
		&sub.Price,
		&sub.SubscribedOn,
		&sub.UnsubscribedOn,
	); err != nil {
		return nil, err
	}

	return &sub, nil
}
