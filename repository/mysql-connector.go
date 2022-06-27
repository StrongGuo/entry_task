package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"item-admin/model"
	"item-admin/proto"
	"time"
)

type MySQL struct {
	dbCon *sql.DB
}

const dsn = "root:chayeDAN11@tcp(127.0.0.1:3306)/item_activity_info?charset=utf8mb4&parseTime=true&loc=Local"

func NewMySQL() *MySQL {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error())
	}

	return &MySQL{dbCon: db}

}

func (m *MySQL) CreateItem(ctx context.Context, i model.Items) (int64, error) {
	queryText := fmt.Sprintf(`INSERT INTO items (item_id, item_name, item_desc, item_price, item_stock, status, creator_id, last_modifier_id) 
	VALUES ('%v','%v','%v','%v','%v','%v','%v','%v')`, i.ItemID, i.ItemName, i.ItemDesc, i.ItemPrice, i.ItemStock, i.Status, i.CreatorID, i.LastModifierID)
	insert, err := m.dbCon.ExecContext(ctx, queryText)

	if err != nil {
		return 0, err
	}

	return insert.LastInsertId()
}

func (m *MySQL) UpdateItem(ctx context.Context, i model.Items) (int64, error) {
	queryText := fmt.Sprintf(`UPDATE items SET item_name='%v', item_desc='%v', item_price=%v, item_stock=%v, status=%v, last_modifier_id=%v 
	WHERE item_id=%v`, i.ItemName, i.ItemDesc, i.ItemPrice, i.ItemStock, i.Status, i.LastModifierID, i.ItemID)
	insert, err := m.dbCon.ExecContext(ctx, queryText)
	n, _ := insert.RowsAffected()

	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 1, err
	}

	return insert.LastInsertId()
}

//RemoveItem 硬删除，删除item的同时，删除item_promotion_rel中的关联信息
func (m *MySQL) RemoveItem(ctx context.Context, i model.Items) (int64, error) {
	querytext1 := fmt.Sprintf("DELETE FROM items WHERE item_id=%v;", i.ItemID)
	querytext2 := fmt.Sprintf("DELETE FROM item_promotion_rel WHERE item_id=%v;", i.ItemID)
	remove, err := m.dbCon.ExecContext(ctx, querytext1)
	_, err = m.dbCon.ExecContext(ctx, querytext2)
	n, _ := remove.RowsAffected()

	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 1, err
	}

	return n, err
}

func (m *MySQL) GetItem(ctx context.Context, a model.ActivityItems) (ra model.ActivityItems, err error) {
	queryText := fmt.Sprintf(`SELECT i.item_id, i.item_name,i.item_desc, i.item_price, r.item_special_price, p.promotion_id 
	FROM items i LEFT JOIN item_promotion_rel r ON i.item_id = r.item_id
	LEFT JOIN promotions p ON r.promotion_id = p.promotion_id 
	AND p.start_time <= CURRENT_TIMESTAMP
	AND p.end_time >= CURRENT_TIMESTAMP
	WHERE i.item_id = %d
	ORDER BY p.ctime DESC
	LIMIT 1;`, a.ItemID)

	rows, err := m.dbCon.QueryContext(ctx, queryText)

	defer rows.Close()

	if err != nil {
		return ra, err
	}

	var activityitem model.ActivityItems
	for rows.Next() {
		rows.Scan(&activityitem.ItemID, &activityitem.ItemName, &activityitem.ItemDesc, &activityitem.ItemPrice, &activityitem.ItemSpecialPrice, &activityitem.PromotionID)
	}
	return activityitem, nil
}

func (m *MySQL) GetPromotions(ctx context.Context, i model.Items) (p []*proto.Promotion, err error) {
	queryText := fmt.Sprintf(`SELECT DISTINCT p.promotion_id, p.promotion_desc, r.item_special_price, p.creator_id, p.last_modifier_id, p.start_time , p.end_time, r.item_id 
	FROM promotions p LEFT JOIN item_promotion_rel r
	ON r.promotion_id = p.promotion_id
	LEFT JOIN items i ON r.item_id = r.item_id
	WHERE r.item_id = %d
	ORDER BY p.promotion_id;`, i.ItemID)

	rows, err := m.dbCon.QueryContext(ctx, queryText)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var promotions []*proto.Promotion
	for rows.Next() {
		var promotion proto.Promotion
		rows.Scan(&promotion.PromotionID, &promotion.PromotionDesc, &promotion.ItemSpecialPrice, &promotion.CreatorID, &promotion.LastModifierID, &promotion.StartTime, &promotion.EndTime, &promotion.ItemID)
		promotions = append(promotions, &promotion)
	}
	return promotions, nil
}

//CreatePromotion 为item新增活动，需同时插入两张表
func (m *MySQL) CreatePromotion(ctx context.Context, p model.Promotions) (int64, error) {
	timeUnix := time.Now().Unix() //生成时间戳作为promotion_id插入
	queryText1 := fmt.Sprintf(`INSERT INTO promotions (promotion_id, promotion_desc, start_time, end_time, creator_id, last_modifier_id) 
	VALUES ('%v','%v','%v','%v','%v','%v')`, timeUnix, p.PromotionDesc, p.StartTime, p.EndTime, p.CreatorID, p.LastModifierID)
	queryText2 := fmt.Sprintf(`INSERT INTO item_promotion_rel (item_id, promotion_id, item_special_price) 
	VALUES ('%v','%v','%v')`, p.ItemID, timeUnix, p.ItemSpecialPrice)
	insert, err := m.dbCon.ExecContext(ctx, queryText1)
	_, err2 := m.dbCon.ExecContext(ctx, queryText2)

	if err != nil || err2 != nil {
		return 0, err
	}

	return insert.LastInsertId()
}

//RemovePromotion 硬删除，删除promotion的同时，删除item_promotion_rel中的关联信息
func (m *MySQL) RemovePromotion(ctx context.Context, p model.Promotions) (int64, error) {
	queryText_1 := fmt.Sprintf("DELETE FROM promotions WHERE promotion_id=%v;", p.PromotionID)
	queryText_2 := fmt.Sprintf("DELETE FROM item_promotion_rel WHERE promotion_id=%v;", p.PromotionID)
	remove, err := m.dbCon.ExecContext(ctx, queryText_1)
	_, err = m.dbCon.ExecContext(ctx, queryText_2)
	n, _ := remove.RowsAffected()

	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 1, err
	}

	return n, err
}
