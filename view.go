/*
A simple wiki engine based on mysql database and golang.
*/
package main

// 导入 kview
import (
    "bufio"
    "bytes"
    "github.com/ziutek/kview"
    "github.com/knieriem/markdown"
)

// 声明我们的维基页面视图
var main_view, edit_view kview.View

func init() {
     // 加载 layout 模板
     layout := kview.New("layout.kt")

     // 加载用来展示文章列表的模板
     article_list := kview.New("list.kt")

     // 创建主页面
     main_view = layout.Copy()
     main_view.Div("left", article_list)
     main_view.Div("right", kview.New("show.kt", utils))

     // 创建编辑页面
     edit_view = layout.Copy()
     edit_view.Div("left", article_list)
     edit_view.Div("right", kview.New("edit.kt"))
}

var (
    mde = markdown.Extensions{
        Smart:        true,
        Dlists:       true,
        FilterHTML:   true,
        FilterStyles: true,
    }
    utils = map[string]interface{} {
        "markdown": func(txt string) []byte {
            p := markdown.NewParser(&mde)
            var buf bytes.Buffer
            w := bufio.NewWriter(&buf)
            r := bytes.NewBufferString(txt)
            p.Markdown(r, markdown.ToHTML(w))
            w.Flush()
            return buf.Bytes()
        },
    }
)
