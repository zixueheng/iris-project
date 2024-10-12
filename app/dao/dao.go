/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-05 11:06:19
 * @LastEditTime: 2024-10-08 09:58:02
 */
package dao

import (
	"errors"
	"iris-project/app"
	"iris-project/global"

	"gorm.io/gorm"
)

// Model 基本模型
type Model interface {
	// 获取模型ID
	GetID() uint32
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return global.Db
}

// GetByID 根据ID查询，txs至多一个
func GetByID(m Model, txs ...*gorm.DB) bool {
	if m.GetID() == 0 {
		return false
	}
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = GetDB()
	}
	if err := tx.First(m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// UpdateByID 根据ID更新指定数据，txs至多一个
func UpdateByID(m Model, data map[string]interface{}, txs ...*gorm.DB) error {
	if m.GetID() == 0 {
		return errors.New("需指定ID")
	}
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = GetDB()
	}
	if err := tx.Model(m).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByID 根据ID删除，txs至多一个
func DeleteByID(m Model, txs ...*gorm.DB) error {
	if m.GetID() == 0 {
		return errors.New("需指定ID")
	}
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = GetDB()
	}
	if err := tx.Unscoped().Delete(m).Error; err != nil {
		return err
	}
	return nil
}

// CreateUpdate 创建或更新，注意不检测唯一索引，txs至多一个
func CreateUpdate(m Model, txs ...*gorm.DB) error {
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = GetDB()
	}
	if m.GetID() == 0 {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(m).Omit("created_at").Save(m).Error; err != nil {
			return err
		}
	}
	return nil
}

// CreateUpdateOmit 创建或更新（更新忽略一些列），注意不检测唯一索引，txs至多一个
func CreateUpdateOmit(m Model, omitColumns []string, txs ...*gorm.DB) error {
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	} else {
		tx = GetDB()
	}
	if m.GetID() == 0 {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(m).Omit(omitColumns...).Save(m).Error; err != nil {
			return err
		}
	}
	return nil
}

// SaveOne 保存一个，tx主要支持事务，一般传nil
func SaveOne(tx *gorm.DB, model interface{}) error {
	if tx == nil {
		tx = GetDB()
	}
	if err := tx.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// SaveAll 保存多个，tx主要支持事务，一般传nil，slice要求是指针切片
func SaveAll(tx *gorm.DB, slice interface{}) error {
	if tx == nil {
		tx = GetDB()
	}
	if err := tx.Create(slice).Error; err != nil {
		return err
	}
	return nil
}

// Count 按指定条件返回数量；tx主要支持事务，一般传nil
func Count(tx *gorm.DB, model interface{}, where map[string]interface{}) (count int64, err error) {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return 0, err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}
	tx.Model(model).Count(&count)
	return
}

// QueryOpts 查询选项
type QueryOpts struct {
	Select  []string
	OrderBy []string
	Preload []string
	Joins   []string // 如：[]string{"left join emails on emails.user_id = users.id"}
	Group   string
	// Limit, Skip int
}

func addQueryOpts(tx *gorm.DB, opts *QueryOpts) {
	if tx == nil || opts == nil {
		return
	}
	if len(opts.Select) > 0 {
		tx.Select(opts.Select)
	}
	for _, orderby := range opts.OrderBy {
		tx.Order(orderby)
	}
	for _, preload := range opts.Preload {
		tx.Preload(preload)
	}
	for _, join := range opts.Joins {
		tx.Joins(join)
	}
	if opts.Group != "" {
		tx.Group(opts.Group)
	}
}

// Scan 将查询结果扫描至结果result；tx主要支持事务，一般传nil；
// options：fields, orderBys, preloads,group []string（注意group只能传一个字符串）
func Scan(tx *gorm.DB, model interface{}, where map[string]interface{}, result interface{}, options ...[]string) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	switch len(options) {
	case 0:
		break
	case 1:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
	case 2:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
	case 3:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
		for _, preload := range options[2] {
			tx = tx.Preload(preload)
		}
	case 4:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
		for _, preload := range options[2] {
			tx = tx.Preload(preload)
		}
		for _, group := range options[3] {
			tx = tx.Group(group)
		}
	default:
		return errors.New("options参数错误，如：fields, orderBys, preloads []string")
	}
	tx.Model(model).Scan(result)
	return nil
}

// ScanOpts 将查询结果扫描至结果result；tx主要支持事务，一般传nil；
func ScanOpts(tx *gorm.DB, model interface{}, where map[string]interface{}, result interface{}, opts *QueryOpts) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	addQueryOpts(tx, opts)

	tx.Model(model).Scan(result)
	return nil
}

// Pluck 查询单个列column，并将结果扫描到切片result；tx主要支持事务，一般传nil
func Pluck(tx *gorm.DB, model interface{}, where map[string]interface{}, column string, result interface{}) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}
	tx.Model(model).Pluck(column, result)
	return nil
}

// Sum 求column字段合计值存放到指针amount
func Sum(tx *gorm.DB, model interface{}, where map[string]interface{}, column string, amount interface{}) error {
	return Pluck(tx, model, where, "COALESCE(SUM("+column+"), 0) as amount", amount)
}

// FindOne 按指定条件查找一个；tx主要支持事务，一般传nil；
// options：fields, orderBys, preloads []string
func FindOne(tx *gorm.DB, model interface{}, where map[string]interface{}, options ...[]string) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	switch len(options) {
	case 0:
		break
	case 1:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
	case 2:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
	case 3:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
		for _, preload := range options[2] {
			tx = tx.Preload(preload)
		}
	default:
		return errors.New("options参数错误，如：fields, orderBys, preloads []string")
	}

	tx.Take(model)
	return nil
}

// FindOneOpts 按指定条件查找一个；tx主要支持事务，一般传nil；
func FindOneOpts(tx *gorm.DB, model interface{}, where map[string]interface{}, opts *QueryOpts) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	addQueryOpts(tx, opts)

	tx.Take(model)
	return nil
}

// FindAll 按指定条件查询所有结果；tx主要支持事务，一般传nil，
// 参数options：fields, orderBys, preloads []string
func FindAll(tx *gorm.DB, dest interface{}, where map[string]interface{}, options ...[]string) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	switch len(options) {
	case 0:
		break
	case 1:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
	case 2:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
	case 3:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
		for _, preload := range options[2] {
			tx = tx.Preload(preload)
		}
	default:
		return errors.New("options参数错误，如：fields, orderBys, preloads []string")
	}

	tx.Find(dest)
	return nil
}

// FindAllOpts 按指定条件查询所有结果；tx主要支持事务，一般传nil，
func FindAllOpts(tx *gorm.DB, dest interface{}, where map[string]interface{}, opts *QueryOpts) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	addQueryOpts(tx, opts)

	tx.Find(dest)
	return nil
}

// FindLimit 按指定条件查询返回有限数量结果；tx主要支持事务，一般传nil，limit返回条数
// 参数options：fields, orderBys, preloads []string
func FindLimit(tx *gorm.DB, dest interface{}, where map[string]interface{}, limit int, options ...[]string) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	switch len(options) {
	case 0:
		break
	case 1:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
	case 2:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
	case 3:
		if len(options[0]) > 0 {
			tx = tx.Select(options[0])
		}
		for _, orderby := range options[1] {
			tx = tx.Order(orderby)
		}
		for _, preload := range options[2] {
			tx = tx.Preload(preload)
		}
	default:
		return errors.New("options参数错误，如：fields, orderBys, preloads []string")
	}

	tx.Limit(limit).Find(dest)
	return nil
}

// FindLimitOpts 按指定条件查询返回有限数量结果；tx主要支持事务，一般传nil，limit返回条数
func FindLimitOpts(tx *gorm.DB, dest interface{}, where map[string]interface{}, limit int, opts *QueryOpts) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	}

	addQueryOpts(tx, opts)

	tx.Limit(limit).Find(dest)
	return nil
}

// UpdateAll 按指定条件更新所有；tx主要支持事务，一般传nil
func UpdateAll(tx *gorm.DB, model interface{}, where map[string]interface{}, data map[string]interface{}) error {
	if tx == nil {
		tx = GetDB()
	}
	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	} else {
		return errors.New("请指定更新条件")
	}
	if err := tx.Model(model).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAll 按指定条件删除所有；tx主要支持事务，一般传nil
func DeleteAll(tx *gorm.DB, model interface{}, where map[string]interface{}) error {
	if tx == nil {
		tx = GetDB()
	}

	if where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(where); err != nil {
			return err
		} else {
			tx = tx.Where(conditionString, conditionValues...)
		}
	} else {
		return errors.New("请指定删除条件")
	}
	if err := tx.Unscoped().Delete(model).Error; err != nil {
		return err
	}
	return nil
}

// Transaction 执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	if err := GetDB().Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return err
		}
		return nil // 返回 nil 提交事务
	}); err != nil {
		return err
	}

	return nil
}

// SearchListData 通用列表查询条件
type SearchListData struct {
	Tx                         *gorm.DB
	Where                      map[string]interface{}
	Not                        map[string]interface{} // map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}: `name NOT IN ("jinzhu", "jinzhu 2")`
	Fields, Preloads, OrderBys []string
	Joins                      []string // 注意字段要加表名，如：[]string{"left join emails on emails.user_id = users.id"}
	Group                      string
	Page, Size                 int
	FindType                   string      // 查询方式：一般不传，其他 Scan
	Model                      interface{} // 用Scan方式查询要传model
}

// GetList 通用列表查询
func (s *SearchListData) GetList(dest interface{}, total *int64) error {
	if s.Tx == nil {
		s.Tx = GetDB()
	}

	if s.Where != nil {
		if conditionString, conditionValues, err := app.BuildCondition(s.Where); err != nil {
			return err
		} else {
			s.Tx = s.Tx.Where(conditionString, conditionValues...)
		}
	}
	if s.Not != nil {
		s.Tx = s.Tx.Not(s.Not)
	}

	if len(s.Fields) > 0 {
		s.Tx = s.Tx.Select(s.Fields)
	}

	if s.Group != "" {
		s.Tx = s.Tx.Group(s.Group)
	}

	for _, preload := range s.Preloads {
		s.Tx = s.Tx.Preload(preload)
	}

	for _, orderby := range s.OrderBys {
		s.Tx = s.Tx.Order(orderby)
	}
	for _, join := range s.Joins {
		s.Tx.Joins(join)
	}

	if s.FindType == "Scan" {
		if s.Model == nil {
			return errors.New("Scan方式查询需指定Model")
		}
		s.Tx.Model(s.Model).Offset((s.Page - 1) * s.Size).Limit(s.Size).Scan(dest).Offset(-1).Limit(-1).Count(total)
	} else {
		s.Tx.Offset((s.Page - 1) * s.Size).Limit(s.Size).Find(dest).Offset(-1).Limit(-1).Count(total)
	}

	return nil
}
