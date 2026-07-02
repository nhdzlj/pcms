package store

// Store 数据访问接口，生产环境用 GORM 实现，测试环境用内存 map 实现
type Store interface {
	Create(value interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Model(value interface{}) Store
	Where(query interface{}, args ...interface{}) Store
	Order(value interface{}) Store
	Offset(offset int) Store
	Limit(limit int) Store
	Preload(query string, args ...interface{}) Store
	Joins(query string, args ...interface{}) Store
	Count(count *int64) error
	Updates(values interface{}) error
	UpdateColumn(column string, value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	Association(column string) AssociationStore
	Begin() (Store, error)
	Commit() error
	Rollback() error
	Error() error
}

// AssociationStore 关联操作接口
type AssociationStore interface {
	Replace(values interface{}) error
}
