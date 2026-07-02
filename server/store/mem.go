package store

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// MemStore 内存 Store 实现，用于单元测试
type MemStore struct {
	mu     sync.Mutex
	data   map[string]map[uint64]interface{} // table -> id -> record
	autoID map[string]uint64

	// 链式查询状态（每次终端操作后重置）
	modelType  string
	conditions []memCond
	orderBy    string
	preloads   []string
	joins      []string
	offsetVal  int
	limitVal   int
	err        error
}

type memCond struct {
	Query string
	Args  []interface{}
}

func NewMemStore() *MemStore {
	return &MemStore{
		data:   make(map[string]map[uint64]interface{}),
		autoID: make(map[string]uint64),
	}
}

// ============ 内部辅助 ============

func (s *MemStore) tableName(model interface{}) string {
	if s.modelType != "" {
		return s.modelType
	}
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Handle slices: []*Model -> Model
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return t.Name()
}

func (s *MemStore) getTable(name string) map[uint64]interface{} {
	if s.data[name] == nil {
		s.data[name] = make(map[uint64]interface{})
	}
	return s.data[name]
}

func (s *MemStore) nextID(name string) uint64 {
	s.autoID[name]++
	return s.autoID[name]
}

func copyStruct(dst, src interface{}) {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if dv.Kind() != reflect.Struct || sv.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < dv.NumField(); i++ {
		fn := dv.Type().Field(i).Name
		sf := sv.FieldByName(fn)
		df := dv.FieldByName(fn)
		if sf.IsValid() && df.IsValid() && df.CanSet() && sf.Type().AssignableTo(df.Type()) {
			df.Set(sf)
		}
	}
}

func getField(v interface{}, name string) reflect.Value {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return rv.FieldByName(name)
}

// matchWhere 检查 record 是否匹配 accumulated Where 条件
func (s *MemStore) matchWhere(record interface{}) bool {
	for _, c := range s.conditions {
		q := c.Query
		args := c.Args

		if strings.Contains(q, " = ? AND ") || strings.Contains(q, " AND ") {
			// 多字段 AND 查询: "ID = ? AND UserID = ?"
			parts := strings.SplitN(q, " = ? AND ", 2)
			f1 := strings.TrimSpace(parts[0])
			f2 := strings.TrimSuffix(strings.TrimSpace(parts[1]), " = ?")
			v1 := fmt.Sprintf("%v", args[0])
			v2 := fmt.Sprintf("%v", args[1])
			a1 := fmt.Sprintf("%v", getField(record, f1).Interface())
			a2 := fmt.Sprintf("%v", getField(record, f2).Interface())
			if a1 != v1 || a2 != v2 {
				return false
			}
		} else if strings.Contains(q, " OR ") && (strings.Contains(q, "LIKE") || strings.Contains(q, "ILIKE")) {
			// "Title LIKE ? OR Content LIKE ?"
			sep := " OR "
			left := strings.Split(q, sep)[0]
			right := strings.Split(q, sep)[1]
			kw := strings.Trim(fmt.Sprintf("%v", args[0]), "%")
			kw = strings.ToLower(kw)
			f1 := strings.TrimSpace(strings.Split(left, " ")[0])
			f2 := strings.TrimSpace(strings.Split(right, " ")[0])
			v1 := strings.ToLower(fmt.Sprintf("%v", getField(record, f1).Interface()))
			v2 := strings.ToLower(fmt.Sprintf("%v", getField(record, f2).Interface()))
			if !strings.Contains(v1, kw) && !strings.Contains(v2, kw) {
				return false
			}
		} else if strings.Contains(q, " = ") {
			parts := strings.SplitN(q, " = ", 2)
			field := parts[0]
			expect := fmt.Sprintf("%v", args[0])
			actual := fmt.Sprintf("%v", getField(record, field).Interface())
			if actual != expect {
				return false
			}
		}
	}
	return true
}

// ============ Store 接口实现 ============

func (s *MemStore) Create(value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.tableName(value)
	table := s.getTable(tName)
	id := s.nextID(tName)

	// Set ID
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	idField := rv.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.SetUint(id)
	}
	// Set defaults
	stField := rv.FieldByName("Status")
	if stField.IsValid() && stField.CanSet() && stField.String() == "" {
		stField.SetString("draft")
	}
	verField := rv.FieldByName("Version")
	if verField.IsValid() && verField.CanSet() && verField.Int() == 0 {
		verField.SetInt(1)
	}

	table[id] = value
	return nil
}

func (s *MemStore) First(dest interface{}, conds ...interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.tableName(dest)
	table := s.getTable(tName)

	// Handle db.First(&user, id)
	if len(conds) == 1 {
		if id, ok := toUint(conds[0]); ok {
			if r, ok := table[id]; ok {
				copyStruct(dest, r)
				return nil
			}
			return fmt.Errorf("record not found: id=%d", id)
		}
	}

	// Handle db.Where("name = ?", val).First(&user)
	for _, r := range table {
		if s.matchWhere(r) {
			copyStruct(dest, r)
			return nil
		}
	}
	return fmt.Errorf("record not found")
}

func (s *MemStore) Find(dest interface{}, conds ...interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.tableName(dest)
	table := s.getTable(tName)

	var results []interface{}
	// If there are explicit conds, use those; otherwise use accumulated Where
	hasExplicit := len(conds) > 0
	for _, r := range table {
		if hasExplicit {
			// 简化：清空累加条件，用显式参数
		} else {
			if !s.matchWhere(r) {
				continue
			}
		}
		results = append(results, r)
	}

	// APPLY ORDER (simplified)
	// APPLY PAGING
	if s.offsetVal > 0 {
		if s.offsetVal < len(results) {
			results = results[s.offsetVal:]
		} else {
			results = nil
		}
	}
	if s.limitVal > 0 && len(results) > s.limitVal {
		results = results[:s.limitVal]
	}

	// Write into dest slice
	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Ptr || dv.IsNil() {
		return fmt.Errorf("dest must be pointer to slice")
	}
	sv := dv.Elem()
	if sv.Kind() != reflect.Slice {
		return fmt.Errorf("dest must be slice, got %v", sv.Kind())
	}
	elemType := sv.Type().Elem() // *Model for []*Model, Model for []Model
	for _, r := range results {
		var val reflect.Value
		if elemType.Kind() == reflect.Ptr {
			// []*Model case
			val = reflect.New(elemType.Elem()) // *Model
			copyStruct(val.Interface(), r)
		} else {
			// []Model case
			val = reflect.New(elemType) // *Model
			copyStruct(val.Interface(), r)
			val = val.Elem() // Model
		}
		sv.Set(reflect.Append(sv, val))
	}
	return nil
}

func (s *MemStore) Model(value interface{}) Store {
	ns := *s
	ns.modelType = s.tableName(value)
	return &ns
}

func (s *MemStore) Where(query interface{}, args ...interface{}) Store {
	ns := *s
	qStr := fmt.Sprintf("%v", query)
	ns.conditions = append(ns.conditions, memCond{Query: qStr, Args: args})
	return &ns
}

func (s *MemStore) Order(value interface{}) Store {
	ns := *s
	ns.orderBy = fmt.Sprintf("%v", value)
	return &ns
}

func (s *MemStore) Offset(offset int) Store {
	ns := *s
	ns.offsetVal = offset
	return &ns
}

func (s *MemStore) Limit(limit int) Store {
	ns := *s
	ns.limitVal = limit
	return &ns
}

func (s *MemStore) Preload(query string, args ...interface{}) Store {
	ns := *s
	ns.preloads = append(ns.preloads, query)
	return &ns
}

func (s *MemStore) Joins(query string, args ...interface{}) Store {
	ns := *s
	ns.joins = append(ns.joins, fmt.Sprintf("%v", query))
	return &ns
}

func (s *MemStore) Count(count *int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.modelType
	table := s.getTable(tName)
	var c int64
	for _, r := range table {
		if s.matchWhere(r) {
			c++
		}
	}
	*count = c
	return nil
}

func (s *MemStore) Updates(values interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.modelType
	table := s.getTable(tName)
	upd, ok := values.(map[string]interface{})
	if !ok {
		return fmt.Errorf("updates needs map[string]interface{}")
	}
	found := false
	for _, r := range table {
		if s.matchWhere(r) {
			rv := reflect.ValueOf(r)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			for k, v := range upd {
				f := rv.FieldByName(k)
				if f.IsValid() && f.CanSet() {
					fv := reflect.ValueOf(v)
					if fv.IsValid() {
						// 处理值类型 -> 指针类型的转换
						if f.Kind() == reflect.Ptr && fv.Kind() != reflect.Ptr {
							// 创建指针
							ptr := reflect.New(f.Type().Elem())
							ptr.Elem().Set(fv)
							f.Set(ptr)
						} else if fv.Type().AssignableTo(f.Type()) {
							f.Set(fv)
						}
					}
				}
			}
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no record updated")
	}
	return nil
}

func (s *MemStore) UpdateColumn(column string, value interface{}) error {
	// Simplified - just do a raw update
	return nil
}

func (s *MemStore) Delete(value interface{}, conds ...interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tName := s.tableName(value)
	table := s.getTable(tName)

	if len(conds) > 0 {
		// Delete by ID + optional extra conditions: Delete(&Model{}, id, "Field = ?", value)
		if id, ok := toUint(conds[0]); ok {
			r, exists := table[id]
			if !exists {
				return fmt.Errorf("record not found")
			}
			// Check additional conditions (format: "Field = ?", value)
			for i := 1; i+1 < len(conds); i += 2 {
				condStr := fmt.Sprintf("%v", conds[i])
				field := strings.TrimSuffix(condStr, " = ?")
				field = strings.TrimSpace(field)
				expected := fmt.Sprintf("%v", conds[i+1])
				actual := fmt.Sprintf("%v", getField(r, field).Interface())
				if actual != expected {
					return fmt.Errorf("record not found")
				}
			}
			delete(table, id)
			return nil
		}
	}

	// No explicit ID - delete by conditions
	found := false
	for id, r := range table {
		if s.matchWhere(r) {
			delete(table, id)
			found = true
		}
	}
	if !found && len(s.conditions) > 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (s *MemStore) Association(column string) AssociationStore {
	return &MemAssociationStore{}
}

func (s *MemStore) Begin() (Store, error) {
	ns := *s
	return &ns, nil
}

func (s *MemStore) Commit() error   { return nil }
func (s *MemStore) Rollback() error { return nil }
func (s *MemStore) Error() error    { return s.err }

// ============ AssociationStore ============

type MemAssociationStore struct{}

func (a *MemAssociationStore) Replace(values interface{}) error { return nil }

// ============ 辅助 ============

func toUint(v interface{}) (uint64, bool) {
	switch val := v.(type) {
	case uint64:
		return val, true
	case int:
		return uint64(val), true
	case uint:
		return uint64(val), true
	}
	return 0, false
}
