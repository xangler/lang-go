package daos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/learn-go/web/module/models"
	"github.com/learn-go/web/pkg/dbutils"
)

type Locker interface {
	GetLockerById(ctx context.Context, id uint64, lockType dbutils.LockType) (*models.Locker, error)
	CountLockers(ctx context.Context, searchArgs *LockerSearchArgs, lockType dbutils.LockType) (int, error)
	ListLockers(ctx context.Context, searchArgs *LockerSearchArgs, offset, limit int, lockType dbutils.LockType) ([]*models.Locker, error)
	AddLocker(ctx context.Context, LockerObj *models.Locker) error
	UpdateLocker(ctx context.Context, LockerObj *models.Locker) error
	UpdateLockerWithVersion(ctx context.Context, LockerObj *models.Locker, version int64) error
	DeleteLocker(ctx context.Context, id uint64) error
	GetLockerByName(ctx context.Context, name string, lockType dbutils.LockType) (*models.Locker, error)
}

const (
	LockerTable      = "locker"
	LockerAdd        = "name, version, master"
	LockerAddPHolder = "?,?,?"
	LockerParams     = "id, " + LockerAdd + ", create_at, update_at"
	LockerSetStr     = "version=?,master=?"
)

type daoLocker struct {
	Q dbutils.Querier
}

type LockerSearchArgs struct {
	Name *string
}

func NewDaoLocker(q dbutils.Querier) (Locker, error) {
	return &daoLocker{
		Q: q,
	}, nil
}

func (d *daoLocker) GetLockerById(ctx context.Context, id uint64, lockType dbutils.LockType) (*models.Locker, error) {

	sb := strings.Builder{}
	args := []interface{}{}
	sb.WriteString("SELECT ")
	sb.WriteString(LockerParams)
	sb.WriteString(fmt.Sprintf(" FROM %s WHERE id = ?", LockerTable))
	args = append(args, id)

	if lockType == dbutils.WriteLock {
		sb.WriteString(" FOR UPDATE")
	}

	row := d.Q.QueryRowContext(
		ctx,
		sb.String(),
		args...,
	)
	if obj, err := d.scanRow(row); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err.Error())
		return nil, err
	} else {
		return obj, nil
	}
}

func (d *daoLocker) GetLockerByName(ctx context.Context, name string, lockType dbutils.LockType) (*models.Locker, error) {
	in := &LockerSearchArgs{
		Name: &name,
	}
	objs, err := d.ListLockers(ctx, in, 0, 0, lockType)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return nil, nil
	} else if len(objs) == 1 {
		return objs[0], nil
	} else {
		return nil, fmt.Errorf("GetLockerByName len is %d", len(objs))
	}
}

func (d *daoLocker) CountLockers(ctx context.Context, searchArgs *LockerSearchArgs, lockType dbutils.LockType) (int, error) {
	err := searchArgs.validate()
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	var count int
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("SELECT COUNT(*) FROM %s", LockerTable))
	whereStr, sqlArgs := searchArgs.toSQL()
	if whereStr != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereStr)
	}
	if sqlArgs == nil {
		sqlArgs = make([]interface{}, 0)
	}
	if lockType == dbutils.WriteLock {
		sb.WriteString(" FOR UPDATE")
	}

	row := d.Q.QueryRowContext(ctx, sb.String(), sqlArgs...)
	if err := row.Scan(&count); err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return count, nil
}

func (d *daoLocker) ListLockers(ctx context.Context, searchArgs *LockerSearchArgs, offset, limit int, lockType dbutils.LockType) ([]*models.Locker, error) {
	sb := strings.Builder{}
	sb.WriteString("SELECT ")
	sb.WriteString(LockerParams)
	sb.WriteString(" FROM ")
	sb.WriteString(LockerTable)

	err := searchArgs.validate()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if limit > 1024 {
		return nil, errors.New("the limit for list Lockers more than 1024")
	}
	whereStr, sqlArgs := searchArgs.toSQL()
	if whereStr != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereStr)
	}
	if sqlArgs == nil {
		sqlArgs = make([]interface{}, 0)
	}

	if limit > 0 {
		sb.WriteString(" LIMIT ? OFFSET ?")
		sqlArgs = append(sqlArgs, limit, offset)
	}
	if lockType == dbutils.WriteLock {
		sb.WriteString(" FOR UPDATE")
	}

	rows, err := d.Q.QueryContext(ctx, sb.String(), sqlArgs...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	objList := make([]*models.Locker, 0, limit)
	for rows.Next() {
		p, err := d.scanRow(rows)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		objList = append(objList, p)
	}
	err = rows.Err()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return objList, nil
}

func (d *daoLocker) AddLocker(ctx context.Context, LockerObj *models.Locker) error {
	_, err := d.Q.ExecContext(
		ctx,
		"INSERT INTO "+LockerTable+" ( "+LockerAdd+" ) "+
			"VALUES ("+LockerAddPHolder+")",
		LockerObj.Name, LockerObj.Version,
		LockerObj.Master,
	)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func (d *daoLocker) UpdateLocker(ctx context.Context, LockerObj *models.Locker) error {
	sqlStr := "UPDATE " + LockerTable + " SET " + LockerSetStr + " WHERE id = ?"
	if rs, err := d.Q.ExecContext(
		ctx, sqlStr,
		LockerObj.Version, LockerObj.Master,
		LockerObj.Id,
	); err != nil {
		log.Error(err.Error())
		return err
	} else if n, err := rs.RowsAffected(); err != nil {
		log.Error(err.Error())
		return err
	} else if n == 0 {
		log.Warnf("update Locker %d change nothing", LockerObj.Id)
	}
	return nil
}

func (d *daoLocker) UpdateLockerWithVersion(ctx context.Context, LockerObj *models.Locker, version int64) error {
	sqlStr := "UPDATE " + LockerTable + " SET " + LockerSetStr + " WHERE id = ? AND version = ?"
	if rs, err := d.Q.ExecContext(
		ctx, sqlStr,
		LockerObj.Version, LockerObj.Master,
		LockerObj.Id, version,
	); err != nil {
		log.Error(err.Error())
		return err
	} else if n, err := rs.RowsAffected(); err != nil {
		log.Error(err.Error())
		return err
	} else if n == 0 {
		log.Warnf("update Locker %d change nothing", LockerObj.Id)
	}
	return nil
}

func (d *daoLocker) DeleteLocker(ctx context.Context, id uint64) error {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE id = ?", LockerTable)
	if rs, err := d.Q.ExecContext(ctx, sqlStr, id); err != nil {
		log.Error(err.Error())
		return err
	} else if n, err := rs.RowsAffected(); err != nil {
		log.Error(err.Error())
		return err
	} else if n == 0 {
		log.Warnf("delete cannot found Lockers %d error", id)
	}
	return nil
}

func (args *LockerSearchArgs) validate() error {
	return nil
}

func (args *LockerSearchArgs) toSQL() (string, []interface{}) {
	if args == nil {
		return "", []interface{}{}
	}
	sqlArgs := []interface{}{}
	sqlStrList := []string{}
	if args.Name != nil {
		sqlStrList = append(sqlStrList, "name=?")
		sqlArgs = append(sqlArgs, *args.Name)
	}
	return strings.Join(sqlStrList[:], " AND "), sqlArgs
}

func (d *daoLocker) scanRow(s dbutils.Scanner) (*models.Locker, error) {
	l := &models.Locker{}
	err := s.Scan(
		&l.Id,
		&l.Name,
		&l.Version,
		&l.Master,
		&l.CreateAt,
		&l.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}
