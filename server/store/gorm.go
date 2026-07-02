package store

import "gorm.io/gorm"

// GormStore 基于 GORM 的 Store 实现（生产环境）
type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) DB() *gorm.DB {
	return s.db
}

func (s *GormStore) Create(value interface{}) error {
	return s.db.Create(value).Error
}

func (s *GormStore) First(dest interface{}, conds ...interface{}) error {
	return s.db.First(dest, conds...).Error
}

func (s *GormStore) Find(dest interface{}, conds ...interface{}) error {
	return s.db.Find(dest, conds...).Error
}

func (s *GormStore) Model(value interface{}) Store {
	return &GormStore{db: s.db.Model(value)}
}

func (s *GormStore) Where(query interface{}, args ...interface{}) Store {
	return &GormStore{db: s.db.Where(query, args...)}
}

func (s *GormStore) Order(value interface{}) Store {
	return &GormStore{db: s.db.Order(value)}
}

func (s *GormStore) Offset(offset int) Store {
	return &GormStore{db: s.db.Offset(offset)}
}

func (s *GormStore) Limit(limit int) Store {
	return &GormStore{db: s.db.Limit(limit)}
}

func (s *GormStore) Preload(query string, args ...interface{}) Store {
	return &GormStore{db: s.db.Preload(query, args...)}
}

func (s *GormStore) Joins(query string, args ...interface{}) Store {
	return &GormStore{db: s.db.Joins(query, args...)}
}

func (s *GormStore) Count(count *int64) error {
	return s.db.Count(count).Error
}

func (s *GormStore) Updates(values interface{}) error {
	return s.db.Updates(values).Error
}

func (s *GormStore) UpdateColumn(column string, value interface{}) error {
	return s.db.UpdateColumn(column, value).Error
}

func (s *GormStore) Delete(value interface{}, conds ...interface{}) error {
	return s.db.Delete(value, conds...).Error
}

func (s *GormStore) Association(column string) AssociationStore {
	return &GormAssociationStore{assoc: s.db.Association(column)}
}

func (s *GormStore) Begin() (Store, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormStore{db: tx}, nil
}

func (s *GormStore) Commit() error {
	return s.db.Commit().Error
}

func (s *GormStore) Rollback() error {
	return s.db.Rollback().Error
}

func (s *GormStore) Error() error {
	return s.db.Error
}

// GormAssociationStore GORM 关联操作实现
type GormAssociationStore struct {
	assoc *gorm.Association
}

func (a *GormAssociationStore) Replace(values interface{}) error {
	return a.assoc.Replace(values)
}
