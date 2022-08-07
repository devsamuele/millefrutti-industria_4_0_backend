package pasteurizer

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store struct {
	db  *sql.DB
	log *log.Logger
}

func (s Store) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func NewStore(db *sql.DB, log *log.Logger) Store {
	return Store{db: db, log: log}
}

func (s Store) CheckLottoAndAr(ctx context.Context, tx *sql.Tx, cd_lotto, cd_ar string) (bool, error) {
	row := tx.QueryRowContext(ctx, `select count(*) from ARLotto where cd_ARLotto = @p1 and cd_AR = @p2`, cd_lotto, cd_ar)
	if err := row.Err(); err != nil {
		return false, nil
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (s Store) CheckLottoAndArInDoc(ctx context.Context, tx *sql.Tx, cd_lotto, cd_ar string) (bool, error) {
	row := tx.QueryRowContext(ctx, `select count(*) from DoRig where cd_ARLotto = @p1 and cd_AR = @p2`, cd_lotto, cd_ar)
	if err := row.Err(); err != nil {
		return false, nil
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (s Store) CreateLottoArca(ctx context.Context, tx *sql.Tx, cd_lotto, cd_ar string, now time.Time) error {
	_, err := tx.ExecContext(ctx, `insert into ARLotto (Cd_ARLotto, Cd_AR, Descrizione, UserIns, UserUpd, TimeIns, TimeUpd) 
	values(@p1,@p2,@p3,@p4,@p5,@p6,@p7)`, cd_lotto, cd_ar, "", "opcua-service", "opcua-service", now, now)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteLottoArca(ctx context.Context, tx *sql.Tx, cd_lotto, cd_ar string) error {
	_, err := tx.ExecContext(ctx, `delete from ARLotto where cd_ARLotto = @p1 and cd_ar = @p2`, cd_lotto, cd_ar)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) QueryWork(ctx context.Context) ([]Work, error) {

	rows, err := s.db.QueryContext(ctx, `select top(50) id, cd_lotto, cd_ar, basil_amount, packages, date, document_created, status, created from xPastorizzatore order by date desc`)
	if err != nil {
		return make([]Work, 0), err
	}
	defer rows.Close()

	works := make([]Work, 0)
	for rows.Next() {
		var w Work
		if err := rows.Scan(&w.ID, &w.CdLotto, &w.CdAr, &w.BasilAmount, &w.Packages, &w.Date, &w.DocumentCreated, &w.Status, &w.Created); err != nil {
			return make([]Work, 0), err
		}
		works = append(works, w)
	}

	return works, nil
}

func (s Store) QueryWorkByID(ctx context.Context, id int) (Work, error) {
	row := s.db.QueryRowContext(ctx, `select top(1) id, cd_lotto, cd_ar, basil_amount, packages, date, document_created, status, created from xPastorizzatore where id = @p1`, id)
	if err := row.Err(); err != nil {
		return Work{}, err
	}

	var w Work
	if err := row.Scan(&w.ID, &w.CdLotto, &w.CdAr, &w.BasilAmount, &w.Packages, &w.Date, &w.DocumentCreated, &w.Status, &w.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Work{}, ErrNotFound
		}
		return Work{}, err
	}

	return w, nil
}

func (s Store) QueryActiveWork(ctx context.Context) (Work, error) {
	row := s.db.QueryRowContext(ctx, `select top(1) id, cd_lotto, cd_ar, basil_amount, packages, date, document_created, status, created from xPastorizzatore where status != 'done'`)
	if err := row.Err(); err != nil {
		return Work{}, err
	}

	var w Work
	if err := row.Scan(&w.ID, &w.CdLotto, &w.CdAr, &w.BasilAmount, &w.Packages, &w.Date, &w.DocumentCreated, &w.Status, &w.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Work{}, ErrNotFound
		}
		return Work{}, err
	}

	return w, nil
}

func (s Store) ExistActiveWork(ctx context.Context) (bool, error) {
	row := s.db.QueryRowContext(ctx, `select count(*) from xPastorizzatore where status != 'done'`)
	if err := row.Err(); err != nil {
		return false, err
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (s Store) DeleteWork(ctx context.Context, tx *sql.Tx, id int) error {
	_, err := tx.ExecContext(ctx, `delete from xPastorizzatore where id = @p1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) InsertWork(ctx context.Context, tx *sql.Tx, w Work) (int, error) {
	row := tx.QueryRowContext(ctx, `insert into xPastorizzatore (cd_lotto, cd_ar, basil_amount, packages, date, document_created, status, created) 
	values(@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8); select ID = convert(bigint, SCOPE_IDENTITY())`, w.CdLotto, w.CdAr, w.BasilAmount, w.Packages, w.Date, w.DocumentCreated, w.Status, w.Created)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Store) UpdateWork(ctx context.Context, tx *sql.Tx, w Work) error {
	_, err := tx.ExecContext(ctx, `update xPastorizzatore 
	set cd_lotto = @p1, cd_ar = @p2, basil_amount = @p3, packages = @p4, date = @p5, document_created = @p6, status = @p7, created = @p8 
	where id = @p9`, w.CdLotto, w.CdAr, w.BasilAmount, w.Packages, w.Date, w.DocumentCreated, w.Status, w.Created, w.ID)
	if err != nil {
		return err
	}

	return nil
}
