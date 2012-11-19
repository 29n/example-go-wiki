/*
维基的MySQl操作器
*/
package main

import (
    "os"
    "log"
    "github.com/ziutek/mymysql/mysql"
    "github.com/ziutek/mymysql/autorc"
    _ "github.com/ziutek/mymysql/thrsafe"
)

const (
    db_proto = "tcp"
    db_addr  = "127.0.0.1:3306"
    db_user  = "go_wiki"
    db_pass  = "go_wiki"
    db_name  = "go_wiki"
)

var (
    // MySQl 连接处理器
    db = autorc.New(db_proto, "", db_addr, db_user, db_pass, db_name)

    // 预备声明
    artlist_stmt, article_stmt, update_stmt *autorc.Stmt
)

func mysqlError(err error) (ret bool) {
    ret = (err != nil)
    if ret {
        log.Println("MySQL error:", err)
    }
    return
}

func mysqlErrExit(err error) {
    if mysqlError(err) {
        os.Exit(1)
    }
}

func init() {
    var err error

    // 初始化命令
    db.Raw.Register("SET NAMES utf8")

    // 准备好服务器商的声明

    artlist_stmt, err = db.Prepare("SELECT id, title FROM articles")
    mysqlErrExit(err)

    article_stmt, err = db.Prepare(
        "SELECT title, body FROM articles WHERE id = ?",
    )
    mysqlErrExit(err)

    update_stmt, err = db.Prepare(
        "INSERT articles (id, title, body) VALUES (?, ?, ?)" +
        " ON DUPLICATE KEY UPDATE title=VALUES(title), body=VALUES(body)",
    )
    mysqlErrExit(err)
}

type ArticleList struct {
    Id, Title int
    Articles []mysql.Row
}


// 返回文章数据列表给 list.kt 模板使用，我们不使用 map 是因为那需要
// 做太多的事情，你或许在以后的项目中应该这么做，但是在这里，为了简单，
// 我们直接提供原始的 query 结果集以及索引给 id 和 title 字段。
func getArticleList() *ArticleList {
    rows, res, err := artlist_stmt.Exec()
    if mysqlError(err) {
        return nil
    }
    return &ArticleList{
        Id:       res.Map("id"),
        Title:    res.Map("title"),
        Articles: rows,
    }
}

type Article struct {
    Id int
    Title, Body string
}

// 获取一篇文章
func getArticle(id int) (article *Article) {
    rows, res, err := article_stmt.Exec(id)
    if mysqlError(err) {
        return
    }
    if len(rows) != 0 {
        article = &Article{
            Id:    id,
            Title: rows[0].Str(res.Map("title")),
            Body:  rows[0].Str(res.Map("body")),
        }
    }
    return
}

// 插入或者更新一篇文章，它返回被更新/新插入的文章记录的id
func updateArticle(id int, title, body string) int {
    _, res, err := update_stmt.Exec(id, title, body)
    if mysqlError(err) {
        return 0
    }
    return int(res.InsertId())
}
