package model

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Model[T any] struct {
	db *gorm.DB

	total uint
	mu    sync.Mutex
}

// New создает новый экземпляр Model
func New[T any](db *gorm.DB) (*Model[T], error) {
	if err := db.AutoMigrate(new(T)); err != nil {
		return nil, err
	}

	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	m := &Model[T]{
		db: db,
	}

	m.total = uint(total)

	return m, nil
}

// Total возвращает общее количество записей в таблице
func (m *Model[T]) Total() uint {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.total
}

// ByID возвращает запись по ID
func (m *Model[T]) ByID(id uint) (*T, error) {
	var result T
	if err := m.db.First(&result, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ByColumn возвращает запись в которой column = value
func (m *Model[T]) ByColumn(column string, value any) (*T, error) {
	var result T
	if err := m.db.First(&result, column, value).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ByColumnList возвращает список записей, где column = value
func (m *Model[T]) ByColumnList(column string, value any, offset, limit int) ([]T, error) {
	var results []T
	if err := m.db.Offset(offset).Limit(limit).Find(&results, column, value).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// ByColumnContains возвращает первую запись, где в указанном столбце содержится подстрока
func (m *Model[T]) ByColumnContains(column string, substrs []string) (*T, error) {
	queryString := strings.Join(slices.Repeat([]string{fmt.Sprintf("%s LIKE ?", column)}, len(substrs)), " AND ")
	querySubstrs := make([]any, len(substrs))
	for i, substr := range substrs {
		querySubstrs[i] = fmt.Sprintf("%%%s%%", substr)
	}
	var result T
	if err := m.db.Model(new(T)).Where(queryString, querySubstrs).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ByColumnContainsAll возвращает список записей, где в указанном столбце содержатся подстроки substrs
func (m *Model[T]) ByColumnContainsList(column string, substrs []string, offset int, limit int) ([]T, error) {
	queryString := strings.Join(slices.Repeat([]string{fmt.Sprintf("%s LIKE ?", column)}, len(substrs)), " AND ")
	querySubstrs := make([]any, len(substrs))
	for i, substr := range substrs {
		querySubstrs[i] = fmt.Sprintf("%%%s%%", substr)
	}
	var results []T
	if err := m.db.Model(new(T)).Where(queryString, querySubstrs...).Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// CountByColumnContains возвращает количество записей, где в указанном столбце содержатся подстроки substrs
func (m *Model[T]) CountByColumnContains(column string, substrs ...string) (int, error) {
	queryString := strings.Join(slices.Repeat([]string{fmt.Sprintf("%s LIKE ?", column)}, len(substrs)), " AND ")
	querySubstrs := make([]any, len(substrs))
	for i, substr := range substrs {
		querySubstrs[i] = fmt.Sprintf("%%%s%%", substr)
	}

	var count int64
	if err := m.db.Model(new(T)).Where(queryString, querySubstrs...).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// Create создает новую запись
func (m *Model[T]) Create(entity *T) error {
	err := m.db.Create(entity).Error
	if err == nil {
		m.total++
	}
	return err
}

// Update обновляет существующую запись
func (m *Model[T]) Update(entity *T) error {
	return m.db.Save(entity).Error
}

// List возвращает список записей с пагинацией
func (m *Model[T]) List(offset int, limit int) ([]T, error) {
	var results []T
	if err := m.db.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Delete "мягко" удаляет запись с возможностью восстановления
func (m *Model[T]) Delete(entity *T) error {
	err := m.db.Delete(entity).Error
	if err == nil {
		m.total--
	}
	return err
}
