package main

import (
    "log"
    "net/http"
    "strconv"
    "strings"
)

type ViewCtx struct {
    Left, Right interface{}
}

// 渲染主页
func show(wr http.ResponseWriter, art_num string) {
    id, _ := strconv.Atoi(art_num)
    main_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
}

// 渲染编辑页面
func edit(wr http.ResponseWriter, art_num string) {
    id, _ := strconv.Atoi(art_num)
    edit_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
}

// 更新数据库以及渲染主页
func update(wr http.ResponseWriter, req *http.Request, art_num string) {
    if req.FormValue("submit") == "保存" {
        id, _ := strconv.Atoi(art_num) // id == 0 表示创建新文章
        id = updateArticle(
            id, req.FormValue("title"), req.FormValue("body"),
        )
        // 如果我们插入一篇瓣文章，我们就修改 art_num 为新插入文章的 id
        // 这使得我们可以在成功插入新数据之后立马展示该条数据
        art_num = strconv.Itoa(id)
    }
    // 重定身至主页面并展示新插入的文章
    http.Redirect(wr, req, "/"+art_num, 303)
    // 我们可以直接使用 show(wr, art_num) 展示新插入的文章
    // 但是请查阅：http://en.wikipedia.org/wiki/Post/Redirect/Get
}

// 根据客户请求的方式以及URL地址来选择使用哪个控制器来处理
func router(wr http.ResponseWriter, req *http.Request) {
    root_path := "/"
    edit_path := "/edit/"

    switch req.Method {
    case "GET":
        switch {
        case req.URL.Path == "/style.css" || req.URL.Path == "/favicon.ico":
            http.ServeFile(wr, req, "static"+req.URL.Path)
        
        case strings.HasPrefix(req.URL.Path, edit_path):
            edit(wr, req.URL.Path[len(edit_path):])

        case strings.HasPrefix(req.URL.Path, root_path):
            show(wr, req.URL.Path[len(root_path):])
        }

    case "POST":
        switch {
        case strings.HasPrefix(req.URL.Path, root_path):
            update(wr, req, req.URL.Path[len(root_path):])
        }
    }
}

func main() {
    err := http.ListenAndServe(":2223", http.HandlerFunc(router))
    if err != nil {
        log.Fatalln("ListenAndServe:", err)
    }
}
