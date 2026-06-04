package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"

	"github.com/google/uuid"
)

type Repositories struct{ db *sql.DB }

func New(db *sql.DB) *Repositories { return &Repositories{db: db} }

func (r *Repositories) CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `INSERT INTO users (email, password_hash) VALUES ($1,$2) RETURNING id,email,password_hash,created_at,updated_at`, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *Repositories) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `SELECT id,email,password_hash,created_at,updated_at FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, response.ErrNotFound
	}
	return u, err
}

func (r *Repositories) ListAccountTypes(ctx context.Context) ([]domain.AccountType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name,description FROM account_types ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountTypes []domain.AccountType
	for rows.Next() {
		var accountType domain.AccountType
		if err := rows.Scan(&accountType.ID, &accountType.Name, &accountType.Description); err != nil {
			return nil, err
		}
		accountTypes = append(accountTypes, accountType)
	}
	return accountTypes, rows.Err()
}

func (r *Repositories) ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT a.id,a.user_id,a.account_type_id,t.name,a.name,a.balance,a.currency,a.gold_grams,a.gold_price_per_gram,a.is_active,a.created_at,a.updated_at FROM accounts a JOIN account_types t ON t.id=a.account_type_id WHERE a.user_id=$1 ORDER BY a.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.AccountTypeID, &a.Type, &a.Name, &a.Balance, &a.Currency, &a.GoldGrams, &a.GoldPrice, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *Repositories) FindAccount(ctx context.Context, userID, accountID uuid.UUID) (domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx, `SELECT a.id,a.user_id,a.account_type_id,t.name,a.name,a.balance,a.currency,a.gold_grams,a.gold_price_per_gram,a.is_active,a.created_at,a.updated_at FROM accounts a JOIN account_types t ON t.id=a.account_type_id WHERE a.user_id=$1 AND a.id=$2`, userID, accountID).Scan(&a.ID, &a.UserID, &a.AccountTypeID, &a.Type, &a.Name, &a.Balance, &a.Currency, &a.GoldGrams, &a.GoldPrice, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, response.ErrNotFound
	}
	return a, err
}

func (r *Repositories) AccountTypeName(ctx context.Context, accountTypeID int) (string, error) {
	var name string
	err := r.db.QueryRowContext(ctx, `SELECT name FROM account_types WHERE id=$1`, accountTypeID).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", response.ErrNotFound
	}
	return name, err
}

func (r *Repositories) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64) (domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx, `INSERT INTO accounts (user_id,account_type_id,name,balance,gold_grams,gold_price_per_gram) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id,user_id,account_type_id,name,balance,currency,gold_grams,gold_price_per_gram,is_active,created_at,updated_at`, userID, accountTypeID, name, balance, goldGrams, goldPrice).Scan(&a.ID, &a.UserID, &a.AccountTypeID, &a.Name, &a.Balance, &a.Currency, &a.GoldGrams, &a.GoldPrice, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return domain.Account{}, err
	}
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM account_types WHERE id=$1`, a.AccountTypeID).Scan(&a.Type)
	return a, nil
}

func (r *Repositories) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error) {
	var current domain.Account
	err := r.db.QueryRowContext(ctx, `SELECT id,user_id,account_type_id,name,balance,currency,gold_grams,gold_price_per_gram,is_active,created_at,updated_at FROM accounts WHERE user_id=$1 AND id=$2`, userID, accountID).Scan(&current.ID, &current.UserID, &current.AccountTypeID, &current.Name, &current.Balance, &current.Currency, &current.GoldGrams, &current.GoldPrice, &current.IsActive, &current.CreatedAt, &current.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, response.ErrNotFound
	}
	if err != nil {
		return domain.Account{}, err
	}
	if name != nil {
		current.Name = *name
	}
	if isActive != nil {
		current.IsActive = *isActive
	}
	err = r.db.QueryRowContext(ctx, `UPDATE accounts SET name=$3,is_active=$4,updated_at=NOW() WHERE user_id=$1 AND id=$2 RETURNING id,user_id,account_type_id,name,balance,currency,gold_grams,gold_price_per_gram,is_active,created_at,updated_at`, userID, accountID, current.Name, current.IsActive).Scan(&current.ID, &current.UserID, &current.AccountTypeID, &current.Name, &current.Balance, &current.Currency, &current.GoldGrams, &current.GoldPrice, &current.IsActive, &current.CreatedAt, &current.UpdatedAt)
	if err != nil {
		return domain.Account{}, err
	}
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM account_types WHERE id=$1`, current.AccountTypeID).Scan(&current.Type)
	return current, nil
}

func (r *Repositories) SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE accounts SET is_active=false,updated_at=NOW() WHERE user_id=$1 AND id=$2`, userID, accountID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repositories) HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM accounts WHERE user_id=$1 AND id=$2`, userID, accountID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repositories) ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error) {
	query := `SELECT id,user_id,name,type,is_default,created_at,updated_at FROM categories WHERE (user_id=$1 OR is_default=true)`
	args := []any{userID}
	if typ != "" {
		query += ` AND type=$2`
		args = append(args, typ)
	}
	query += ` ORDER BY is_default DESC,name ASC`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *Repositories) CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRowContext(ctx, `INSERT INTO categories (user_id,name,type,is_default) VALUES ($1,$2,$3,false) RETURNING id,user_id,name,type,is_default,created_at,updated_at`, userID, name, typ).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *Repositories) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRowContext(ctx, `UPDATE categories SET name=$3,updated_at=NOW() WHERE user_id=$1 AND id=$2 AND is_default=false RETURNING id,user_id,name,type,is_default,created_at,updated_at`, userID, categoryID, name).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Category{}, response.ErrNotFound
	}
	return c, err
}

func (r *Repositories) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE user_id=$1 AND id=$2 AND is_default=false`, userID, categoryID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repositories) CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	dbtx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Transaction{}, err
	}
	defer dbtx.Rollback()

	var delta float64
	switch tx.Type {
	case "income":
		delta = tx.Amount
	case "expense":
		delta = -tx.Amount
	case "transfer":
		delta = -tx.Amount
	default:
		return domain.Transaction{}, response.ErrBadRequest
	}

	if delta < 0 {
		var currentBalance float64
		err := dbtx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE user_id=$1 AND id=$2 AND is_active=true FOR UPDATE`, tx.UserID, tx.AccountID).Scan(&currentBalance)
		if err != nil {
			return domain.Transaction{}, response.ErrNotFound
		}
		if currentBalance+delta < 0 {
			return domain.Transaction{}, response.ErrInsufficient
		}
	}

	var sourceGoldDelta *float64
	if tx.GoldGrams != nil && delta < 0 {
		value := -*tx.GoldGrams
		sourceGoldDelta = &value
	} else if tx.GoldGrams != nil && delta > 0 {
		sourceGoldDelta = tx.GoldGrams
	}
	if err := updateAccountBalance(ctx, dbtx, tx.UserID, tx.AccountID, delta, sourceGoldDelta); err != nil {
		return domain.Transaction{}, err
	}
	if tx.Type == "transfer" {
		if tx.TransferAccountID == nil || *tx.TransferAccountID == tx.AccountID {
			return domain.Transaction{}, response.ErrBadRequest
		}
		var destinationGoldDelta *float64
		if tx.GoldGrams != nil {
			destinationGoldDelta = tx.GoldGrams
		}
		if err := updateAccountBalance(ctx, dbtx, tx.UserID, *tx.TransferAccountID, tx.Amount, destinationGoldDelta); err != nil {
			return domain.Transaction{}, err
		}
	}

	var created domain.Transaction
	err = dbtx.QueryRowContext(ctx, `INSERT INTO transactions (user_id,account_id,category_id,type,amount,gold_grams,description,date,transfer_account_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id,user_id,account_id,category_id,type,amount,gold_grams,description,date,transfer_account_id,created_at,updated_at`, tx.UserID, tx.AccountID, tx.CategoryID, tx.Type, tx.Amount, tx.GoldGrams, tx.Description, tx.Date, tx.TransferAccountID).Scan(&created.ID, &created.UserID, &created.AccountID, &created.CategoryID, &created.Type, &created.Amount, &created.GoldGrams, &created.Description, &created.Date, &created.TransferAccountID, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return domain.Transaction{}, err
	}
	if err := dbtx.Commit(); err != nil {
		return domain.Transaction{}, err
	}
	return created, nil
}

func updateAccountBalance(ctx context.Context, tx *sql.Tx, userID, accountID uuid.UUID, delta float64, goldGramDelta *float64) error {
	var goldGramDeltaArg any
	if goldGramDelta != nil {
		goldGramDeltaArg = *goldGramDelta
	}
	res, err := tx.ExecContext(ctx, `
		UPDATE accounts a
		SET
			balance = balance + $3,
			gold_grams = CASE
				WHEN t.name = 'gold'
					AND a.gold_grams IS NOT NULL
					AND $4::numeric IS NOT NULL
				THEN GREATEST(0, a.gold_grams + $4::numeric)
				WHEN t.name = 'gold'
					AND a.gold_grams IS NOT NULL
					AND COALESCE(NULLIF(a.gold_price_per_gram, 0), gp.price_per_gram) > 0
				THEN GREATEST(0, a.gold_grams + ($3 / COALESCE(NULLIF(a.gold_price_per_gram, 0), gp.price_per_gram)))
				ELSE a.gold_grams
			END,
			gold_price_per_gram = CASE
				WHEN t.name = 'gold'
					AND a.gold_grams IS NOT NULL
					AND COALESCE(NULLIF(a.gold_price_per_gram, 0), gp.price_per_gram) > 0
				THEN COALESCE(NULLIF(a.gold_price_per_gram, 0), gp.price_per_gram)
				ELSE a.gold_price_per_gram
			END,
			updated_at = NOW()
		FROM account_types t
		LEFT JOIN gold_prices gp ON gp.id = 1
		WHERE t.id = a.account_type_id
			AND a.user_id = $1
			AND a.id = $2
			AND a.is_active = true`, userID, accountID, delta, goldGramDeltaArg)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repositories) ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
	query := `SELECT tr.id,tr.user_id,tr.account_id,tr.category_id,tr.type,tr.amount,tr.gold_grams,tr.description,tr.date,tr.transfer_account_id,tr.created_at,tr.updated_at,a.account_type_id,a.name,a.balance,a.currency,a.gold_grams,a.gold_price_per_gram,a.is_active,a.created_at,a.updated_at,at.name,c.name,c.type,c.is_default,c.created_at,c.updated_at FROM transactions tr JOIN accounts a ON a.id=tr.account_id JOIN account_types at ON at.id=a.account_type_id LEFT JOIN categories c ON c.id=tr.category_id WHERE tr.user_id=$1`
	args := []any{userID}
	i := 2
	if start != "" {
		query += ` AND tr.date >= $` + strconv.Itoa(i)
		args = append(args, start)
		i++
	}
	if end != "" {
		query += ` AND tr.date < ($` + strconv.Itoa(i) + `::date + INTERVAL '1 day')`
		args = append(args, end)
		i++
	}
	if accountID != "" {
		query += ` AND tr.account_id=$` + strconv.Itoa(i)
		args = append(args, accountID)
		i++
	}
	if categoryID != "" {
		query += ` AND tr.category_id=$` + strconv.Itoa(i)
		args = append(args, categoryID)
		i++
	}
	if typ != "" {
		query += ` AND tr.type=$` + strconv.Itoa(i)
		args = append(args, typ)
		i++
	}
	query += ` ORDER BY tr.date DESC,tr.created_at DESC LIMIT $` + strconv.Itoa(i) + ` OFFSET $` + strconv.Itoa(i+1)
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Transaction
	for rows.Next() {
		var tr domain.Transaction
		var catName, catType sql.NullString
		var catDefault sql.NullBool
		var catCreatedAt, catUpdatedAt sql.NullTime
		tr.Account = &domain.Account{}
		tr.Category = &domain.Category{}
		if err := rows.Scan(&tr.ID, &tr.UserID, &tr.AccountID, &tr.CategoryID, &tr.Type, &tr.Amount, &tr.GoldGrams, &tr.Description, &tr.Date, &tr.TransferAccountID, &tr.CreatedAt, &tr.UpdatedAt, &tr.Account.AccountTypeID, &tr.Account.Name, &tr.Account.Balance, &tr.Account.Currency, &tr.Account.GoldGrams, &tr.Account.GoldPrice, &tr.Account.IsActive, &tr.Account.CreatedAt, &tr.Account.UpdatedAt, &tr.Account.Type, &catName, &catType, &catDefault, &catCreatedAt, &catUpdatedAt); err != nil {
			return nil, err
		}
		tr.Account.ID = tr.AccountID
		if tr.CategoryID != nil && catName.Valid {
			tr.Category.ID = *tr.CategoryID
			tr.Category.Name = catName.String
			tr.Category.Type = catType.String
			tr.Category.IsDefault = catDefault.Bool
			if catCreatedAt.Valid {
				tr.Category.CreatedAt = catCreatedAt.Time
			}
			if catUpdatedAt.Valid {
				tr.Category.UpdatedAt = catUpdatedAt.Time
			}
		} else {
			tr.Category = nil
		}
		result = append(result, tr)
	}
	return result, rows.Err()
}

func (r *Repositories) NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error) {
	accounts, err := r.ListAccounts(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	var total float64
	active := make([]domain.Account, 0, len(accounts))
	for _, a := range accounts {
		if a.IsActive {
			total += a.Balance
			active = append(active, a)
		}
	}
	return total, active, nil
}

func (r *Repositories) LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error) {
	var price domain.GoldPrice
	err := r.db.QueryRowContext(ctx, `SELECT price_per_gram,source,fetched_at,updated_at FROM gold_prices WHERE id=1`).Scan(&price.PricePerGram, &price.Source, &price.FetchedAt, &price.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.GoldPrice{}, response.ErrNotFound
	}
	return price, err
}

func (r *Repositories) SaveGoldPrice(ctx context.Context, price domain.GoldPrice) (domain.GoldPrice, error) {
	if price.FetchedAt.IsZero() {
		price.FetchedAt = time.Now().UTC()
	}
	if price.Source == "" {
		price.Source = "unknown"
	}
	err := r.db.QueryRowContext(ctx, `INSERT INTO gold_prices (id,price_per_gram,source,fetched_at,updated_at) VALUES (1,$1,$2,$3,NOW()) ON CONFLICT (id) DO UPDATE SET price_per_gram=EXCLUDED.price_per_gram,source=EXCLUDED.source,fetched_at=EXCLUDED.fetched_at,updated_at=NOW() RETURNING price_per_gram,source,fetched_at,updated_at`, price.PricePerGram, price.Source, price.FetchedAt).Scan(&price.PricePerGram, &price.Source, &price.FetchedAt, &price.UpdatedAt)
	if err != nil {
		return domain.GoldPrice{}, err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO gold_price_history (price_date,price_per_gram,source,fetched_at) VALUES ($1::date,$2,$3,$4) ON CONFLICT (price_date) DO UPDATE SET price_per_gram=EXCLUDED.price_per_gram,source=EXCLUDED.source,fetched_at=EXCLUDED.fetched_at WHERE gold_price_history.fetched_at <= EXCLUDED.fetched_at`, price.FetchedAt.Format("2006-01-02"), price.PricePerGram, price.Source, price.FetchedAt)
	return price, err
}

func (r *Repositories) ListGoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT price_date::text,price_per_gram,source FROM gold_price_history WHERE price_date >= CURRENT_DATE - (($1::int - 1) * INTERVAL '1 day') ORDER BY price_date ASC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.GoldPriceHistoryPoint
	for rows.Next() {
		var point domain.GoldPriceHistoryPoint
		if err := rows.Scan(&point.Date, &point.PricePerGram, &point.Source); err != nil {
			return nil, err
		}
		history = append(history, point)
	}
	return history, rows.Err()
}

func (r *Repositories) RefreshGoldAccountBalances(ctx context.Context, price domain.GoldPrice) error {
	_, err := r.db.ExecContext(ctx, `UPDATE accounts a SET balance=ROUND((a.gold_grams * $1)::numeric, 2), gold_price_per_gram=$1, updated_at=NOW() FROM account_types t WHERE t.id=a.account_type_id AND t.name='gold' AND a.gold_grams IS NOT NULL`, price.PricePerGram)
	return err
}

func (r *Repositories) CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.Budget, error) {
	var b domain.Budget
	err := r.db.QueryRowContext(ctx, `INSERT INTO budgets (user_id,category_id,month,year,amount) VALUES ($1,$2,$3,$4,$5) RETURNING id,user_id,category_id,month,year,amount,created_at,updated_at`, userID, categoryID, month, year, amount).Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Month, &b.Year, &b.Amount, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return domain.Budget{}, err
	}
	cat, catErr := r.fetchCategory(ctx, categoryID)
	if catErr == nil {
		b.Category = &cat
	}
	return b, nil
}

func (r *Repositories) ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.Budget, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT b.id,b.user_id,b.category_id,c.name,b.month,b.year,b.amount,b.created_at,b.updated_at FROM budgets b JOIN categories c ON c.id=b.category_id WHERE b.user_id=$1 AND b.month=$2 AND b.year=$3 ORDER BY c.name ASC`, userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var budgets []domain.Budget
	for rows.Next() {
		var b domain.Budget
		var cat domain.Category
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &cat.Name, &b.Month, &b.Year, &b.Amount, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		cat.ID = b.CategoryID
		b.Category = &cat
		budgets = append(budgets, b)
	}
	return budgets, rows.Err()
}

func (r *Repositories) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.Budget, error) {
	var b domain.Budget
	err := r.db.QueryRowContext(ctx, `UPDATE budgets SET amount=$3,updated_at=NOW() WHERE user_id=$1 AND id=$2 RETURNING id,user_id,category_id,month,year,amount,created_at,updated_at`, userID, budgetID, amount).Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Month, &b.Year, &b.Amount, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Budget{}, response.ErrNotFound
	}
	if err != nil {
		return domain.Budget{}, err
	}
	cat, catErr := r.fetchCategory(ctx, b.CategoryID)
	if catErr == nil {
		b.Category = &cat
	}
	return b, nil
}

func (r *Repositories) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM budgets WHERE user_id=$1 AND id=$2`, userID, budgetID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repositories) SpendingByCategoryInRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]domain.SpendingCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id as category_id, c.name, COALESCE(SUM(tr.amount),0) as amount FROM transactions tr JOIN categories c ON c.id=tr.category_id WHERE tr.user_id=$1 AND tr.type='expense' AND tr.date >= $2 AND tr.date < $3 GROUP BY c.id,c.name ORDER BY c.name ASC`, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.SpendingCategory
	for rows.Next() {
		var item domain.SpendingCategory
		if err := rows.Scan(&item.CategoryID, &item.Name, &item.Amount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repositories) SpendingByCategory(ctx context.Context, userID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error) {
	var totalExpense float64
	rows, err := r.db.QueryContext(ctx, `SELECT COALESCE(c.name,'Uncategorized'), SUM(tr.amount) FROM transactions tr LEFT JOIN categories c ON c.id=tr.category_id WHERE tr.user_id=$1 AND tr.type='expense' AND tr.date >= $2 AND tr.date < $3 GROUP BY COALESCE(c.name,'Uncategorized') ORDER BY SUM(tr.amount) DESC`, userID, start, end)
	if err != nil {
		return 0, nil, 0, err
	}
	defer rows.Close()
	var items []domain.SpendingCategory
	for rows.Next() {
		var item domain.SpendingCategory
		if err := rows.Scan(&item.Name, &item.Amount); err != nil {
			return 0, nil, 0, err
		}
		totalExpense += item.Amount
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return 0, nil, 0, err
	}
	for i := range items {
		if totalExpense > 0 {
			items[i].Percentage = items[i].Amount / totalExpense * 100
		}
	}
	var totalIncome float64
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(tr.amount), 0) FROM transactions tr WHERE tr.user_id=$1 AND tr.type='income' AND tr.date >= $2 AND tr.date < $3`, userID, start, end).Scan(&totalIncome)
	if err != nil {
		return 0, nil, 0, err
	}
	return totalExpense, items, totalIncome, nil
}

func (r *Repositories) fetchCategory(ctx context.Context, categoryID uuid.UUID) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRowContext(ctx, `SELECT id,user_id,name,type,is_default,created_at,updated_at FROM categories WHERE id=$1`, categoryID).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Category{}, response.ErrNotFound
	}
	return c, err
}
