package database

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
    DBMS = "mysql"
    USER = "test_user" // DBのユーザ名
    PASS = "test_user" // DBのパスワード
    HOST = "localhost" // DBのホスト名もしくはIP
    PORT = "3306"      // DBのポート
    DBNAME = "exampledb" // DB名
    CONNECT = USER + ":" + PASS + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME + "?parseTime=true"
)

var db *gorm.DB

type Item struct {
    gorm.Model
    ItemCode string
    Text   string
}

// DB初期化
func Init() {
    var err error
    db, err = gorm.Open(DBMS, CONNECT)
    if err != nil {
        panic("データベースの接続に失敗しました")
    }
    db.AutoMigrate(&Item{})
}

// 登録
func Insert(itemCode string, text string) uint {
    var item Item
    item.ItemCode = itemCode
    item.Text = text
//    db.Create(&Item{ItemCode: itemCode, Text: text})
    db.Create(&item)
    return item.ID
}

// 更新
func Update(id int, itemCode string, text string) {
    var item Item
    db.First(&item, id)
    item.ItemCode = itemCode
    item.Text = text
    db.Save(&item)
}

// 削除
func Delete(id int) {
    var item Item
    db.First(&item, id)
    db.Delete(&item)
}

// 全取得
func SelectAll() []Item {
    var items []Item
    db.Order("created_at desc").Find(&items)
    return items
}

// 1行取得
func SelectRow(id int) Item {
    var item Item
    db.First(&item, id)
    return item
}

// DBクローズ
func Close() {
    db.Close()
}
