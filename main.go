package main

import (
    "sns-message-sender/database"
    "sns-message-sender/message"
    "strconv"
    "net/http"
    "github.com/gin-gonic/gin"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "os"
    "os/signal"
    "log"
    "fmt"
)

const (
    PORT = "80"
)

func main() {

    go func(){
        router := gin.Default()
        router.LoadHTMLGlob("templates/*.html")
        database.Init()

        // 一覧画面表示
        router.GET("/", func(ctx *gin.Context) {
            items := database.SelectAll()
            ctx.HTML(http.StatusOK, "index.html", gin.H{
                "items": items,
            })
        })
        // 登録
        router.POST("/insert", func(ctx *gin.Context) {
            id := database.Insert(ctx.PostForm("itemCode"), ctx.PostForm("text"))
            message.SendMessage("INSERT", fmt.Sprint(id), ctx.PostForm("itemCode"), ctx.PostForm("text"), ctx.PostForm("filterValue"))
            ctx.Redirect(http.StatusFound, "/")
        })
        // 編集画面表示
        router.GET("/detail/:id", func(ctx *gin.Context) {
            item := database.SelectRow(convertId(ctx.Param("id")))
            ctx.HTML(http.StatusOK, "detail.html", gin.H{"item": item})
        })
        // 更新
        router.POST("/update/:id", func(ctx *gin.Context) {
            database.Update(convertId(ctx.Param("id")), ctx.PostForm("itemCode"), ctx.PostForm("text"))
            message.SendMessage("UPDATE", ctx.Param("id"), ctx.PostForm("itemCode"), ctx.PostForm("text"), ctx.PostForm("filterValue"))
            ctx.Redirect(http.StatusFound, "/")
        })
        // 削除
        router.POST("/delete/:id", func(ctx *gin.Context) {
            database.Delete(convertId(ctx.Param("id")))
            message.SendMessage("DELETE", ctx.Param("id"), ctx.PostForm("itemCode"), ctx.PostForm("text"), ctx.PostForm("filterValue"))
            ctx.Redirect(http.StatusFound, "/")
        })

        router.Run(":" + PORT)
    }()

    // Ctrl + c で停止した際にDBをクローズする
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    database.Close()
}

// IDを数値に変換する
func convertId(idStr string) int {
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Fatal("IDの変換に失敗しました")
    }
    return id
}


