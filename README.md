# 如何基于Go创建数据库驱动的Web应用

- url_title: how-to-write-database-driven-web-application-using-go
- date: 2012-11-19 19:45:50
- update: 2012-11-19 19:46:50
- public: true
- tags: Wiki, Golang, Go, Web Application, Programming, Database, MySQL
- category: tech
- summary: 本文将尝试前着去介绍如何使用 kview/kasia.go 以及 MyMySQL 基于Go语言开发一个小型的简单的数据库驱动的Web应用，就像大家做的一样，我们来开发一个小型的维基系统。

----------------------------------------------------------------------

> 本文原文为[How to write database-driven Web application using Go](https://github.com/ziutek/simple_go_wiki)的 *README.md* 文件，如果您想查看本文的原文，请点击前面的英文原文标题，找到该项目的 *README.md* 文件即可，格式为 *MarkDown*，如果你需要HTML版本的，可能还需要自己安装MarkDown相关的工具。

-----------------------------------------------------------------------

本文将尝试前着去介绍如何使用 kview/kasia.go 以及 MyMySQL 基于Go语言开发一个小型的简单的数据库驱动的Web应用，就像大家做的一样，我们来开发一个小型的维基系统。

## 对您个人的要求

+ 一些程序开发经验；
+ 关于 HTML 与 HTTP 的基础知识；
+ 了解MySQL 以及 MySQL 命令行工具的使用；
+ 在MySQL中创建数据库；
+ 最新的 Go 编译器 － 移步 [Go 首页](http://weekly.golang.org/doc/install.html) 以了解更多详情。

## 数据库

让我们从创建应用所使用的数据库开始该项目。

如果你还没有安装MySQL，那么需要你首先安装它，接下来我们会用到它，如果你已经安装了，那么首先我们先来创建该应用数据库。

    cox@CoxStation:~$ mysql -u root -p
    Enter password: 
    Welcome to the MySQL monitor.  Commands end with ; or \g.
    Your MySQL connection id is 120
    Server version: 5.5.28-0ubuntu0.12.04.2 (Ubuntu)

    Copyright (c) 2000, 2012, Oracle and/or its affiliates. All rights reserved.

    Oracle is a registered trademark of Oracle Corporation and/or its
    affiliates. Other names may be trademarks of their respective
    owners.

    Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

    mysql> create database go_wiki;
    Query OK, 1 row affected (0.00 sec)

    mysql> use go_wiki;
    Database changed
    mysql> CREATE TABLE articles (
        ->     id INT AUTO_INCREMENT PRIMARY KEY,
        ->     title VARCHAR(80) NOT NULL,
        ->     body TEXT NOT NULL
        -> ) DEFAULT CHARSET=utf8;
    Query OK, 0 rows affected (0.09 sec)

    mysql> exit
    Bye

到现在为止你已经创建了一个可能使用的数据库，同时还添加了一张表 *articles*用来保存我们的维基中的文章数据，你可以在应用中直接使用 *root* 帐户连接数据库，当然，更好的办法是创建一个独立的用来测试我们应用用户：

    mysql> GRANT SELECT, INSERT, UPDATE, DELETE ON articles TO go_wiki@localhost;
    Query OK, 0 rows affected (0.00 sec)

    mysql> SET PASSWORD FOR go_wiki@localhost = PASSWORD('go_wiki');
    Query OK, 0 rows affected (0.00 sec)

现在我们记下刚才所获得的数据：

+ 数据库地址：*localhost*
+ 数据库名称：*go_wiki*
+ 数据库用户：*go_wiki*
+ 数据库密码：*go_wiki*

## 视图

现在我们开始写点儿 Go 代码了，像以前一样，创建一个工作空间（如果你还不知道如何创建工作空间，请阅读我以前的文章[《如何写 Go 程序》](/article/how-to-write-go-code.html)），我在这里就简单的将创建的流程所运行的命令复制到这里：

    cox@CoxStation:~$ mkdir go_wiki
    cox@CoxStation:~$ cd go_wiki
    cox@CoxStation:~/go_wiki$ export GOPATH=$HOME/go_wiki
    cox@CoxStation:~/go_wiki$ mkdir bin src pkg
    cox@CoxStation:~/go_wiki$ export PATH=$PATH:$HOME/go_wiki/bin

我所创建的工作空间是：

+ HOME: /home/cox
+ GOPATH: $HOME/go_wiki

定义应用的视图我使用 *kview* 以及 *kasia.go* 这两个包，你应该先安装它们：

    cox@CoxStation:~/go_wiki$ go get github.com/ziutek/kview

上面的命令会因你的网络不同而需要不同的时间，因为Go需要从网络上下载你所需要的一切东西，所以，请保证你的计算机是已经连网了的，下载完成之后，它会自动的为你安装 *kasia.go* 和 *kview*。

下一步，在 *$GOPATH/src* 中创建我们的项目：

在 *$GOPATH/src/go_wiki*目录中，创建 *view.go* 文件，它的内容如下：

    /*
    A simple wiki engine based on mysql database and golang.
    */
    package main

    // 导入 kview
    import "github.com/ziutek/kview"

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
        main_view.Div("right", kview.New("show.kt"))

        // 创建编辑页面
        edit_view = layout.Copy()
        edit_view.Div("left", article_list)
        edit_view.Div("right", kview.New("edit.kt"))
    }

如你所看到的，我们的应用将由两人个页面组成：

+ *main_view* － 用来向用户展示文章内容
+ *edit_view* － 向用户提供创建与编辑功能

任何一个页面又都有两个子栏组成：

+ *left* － 已存在的文章列表
+ *right* － 根据页面不同而展示不同的内容（主页为内容展示，编辑页为编辑表单）

接下来我们来创建第一个 *Kasia* 模板，它将定义我们网站的整个页面布局，我们需要在模板目录中创建一个  *layout.kt* 文件：

    cox@CoxStation:~/go_wiki/src/go_wiki$ mkdir templates
    cox@CoxStation:~/go_wiki/src/go_wiki$ emacs templates/layout.kt

它的内容为：

    
    <!DOCTYPE html>
    <html lang="zh-CN">
    <head>
    <meta charset="utf-8" />
    <title>MySQL数据库驱动的基于Go的维基系统</title>
    <link href="/style.css" type="text/css" rel="stylesheet" />
    </head>
    <body>
    <div class="container">
        <h1>MySQL数据库驱动的基于Go的维基系统</h1>
        <div class="columns">
            <div class="left">$left.Render(Left)</div>
            <div class="right">$right.Render(Right)</div>
        </div>
    </div>
    </body>
    </html>

上面这个简单布局的职责是：

+ 创建基础的HTML文档结构，包括 *doctype* 、 *head* 、 *body* 等
+ 渲染 *left* 与 *right* 两个分栏的具体内容，使用的数据是由 *Left* 与 *Right* 提供的

*Render* 方法是由 *kview* 包提供的，它能在特定的调用它的位置根据提供给它的数据渲染出子视图，这样的子视图可以有它自己的布局、div元素等等，并且子视图中还可以其自己的子视图（但是对于我们现在的这个小项目来说是完全没有必要的了）。

下一步，我们创建 *list.kt* ，它将被渲染到 *left* 中，它的内容是：

    <a href="/edit/" class="button">创建新文章</a>
    <hr />
    <h2>最近更新</h2>
    <ul class="list">
    $for _, article in Articles:
        <li><a href="$article[Id]">$article[Title]</a></li>
    $end
    </ul>

这个模板输出一个*创建新文章*的链接以及一个文章标题链接列表。它使用到了一个 *for* 声明用来遍历 *Articles* 列表（它是一个分片(slice))，之后的遍历得到的每一个项都使用 *article[Id}* 来创建与该文章对应的页面的URL地址，以及 *article[Title]* 来创建标题，*Articles*、*Id*、以及*Title*都是 *ArticleList* 类型中定义的成员（我们会在本文档的后面来定义它），*Id*与*Title*都将对应到*Row*片段中相应的记录，*article*将直接存储从数据库中取得的记录中的某一条。

*for* 声明在本地创建两个局部变量(*\_* 与 *article*），第一个是遍历次数，我们不需要它，所以直接将其丢掉，第二个就是从数据库取得的一条记录，对于那个遍历次数，虽然我们这里没有使用到它，但是它却很有用，比如下面这样的：

    $for nn+, article in Articles:
        <li class="$if even(nn):even$else:odd$end">
            <a href="$article[Id]">$article[Title]</a>
        </li>
    $end

这会在每一个记录上面添加一个类，如果是第奇数条记录，则添加 *odd*，否则添加 *even*，我们这里使用 nn+ 是因为我们希望第一条记录在输入时， nn 不会 0，而是1。

现在来创建 *show.kt* 模板，它将被用来渲染文章数据：

    <div class="article">
    $if Id:
        <h2>$Title</h2>
        <div class="content">
            $Body
        </div>
        <p class="actions">
            <a href="/edit/$Id" class="edit button">编辑本文</a>
        </p>
    $else:
        <h2>维其示例页面</h2>
        <div class="content">
            <p>您现在所看到是本维基的示例页面，这说明您现在还没有指定任何内容。</p>
            <p>本维基使用MySQl数据库存储数据，基于Go语言开发。</p>
            <p>点击下方或者左侧的"创建新文章"链接以创建您的第一篇文章。</p>
            <h3>本维其使用到的技术</h3>
            <ul>
                <li><a href='https://github.com/ziutek/kasia.go'>kasia.go</a></li>
                <li><a href='https://github.com/ziutek/kview'>kview</a></li>
                <li><a href='https://github.com/ziutek/mymysql'>MyMySQL</a></li>
            </ul>
        </div>
        <p class="actions">
            <a href="/edit/" class="button">创建新文章</a>
        </p>
    $end
    </div>

在这里你可以看到，我们使用了一个 *if-else*声明，如果我们指定了展示哪条文章数据了，那么就展示这些数据，如果没有的话，我们就展示一个默认的内容页（当然，这个内容页是无法让前端用户修改的）。

### 插一点有半 *上下文堆栈* 的知识

要使用 *kview* 包渲染某些模板，你需要用到两人个方法，在 Go 代码中的话，就是 ＊Exec*，在模板代码中的话，就是 *Render*，通常的你还需要给它们传递一些变量，比如像下面这样的：

    v.Exec(wr, a, b)
    v.Render(a, b)

与视图 *v* 关联的模板将以下面这样的上下文堆栈渲染：

    []interface{}{globals, a, b}

你可以看到，这里有一个 *globals* 变量，它是一个 *map*，包含了一些全局变量：

+ 子视图（或者子模板）通过 *Div* 方法添加到视图 *v* 中
+ *len* 以及 *fmt* 工具
+ 你传递给 *New* 函数的变量也会被动态的添加到这些全局变量中

如果你想了解得更详细一些，可以看看 [*kview文档*](https://github.com/ziutek/kview/blob/master/README.md)。

变量 *b* 处于这个堆栈的最底端，如果你在模板中这样写：

    $x $@[1].y

那么 *Exec* 或者 *Render* 方法将以下面这种方式去搜索 *x* 或者 *y* 属性：

+ 首先在 *b* 中搜索 *x*，如果没有找到，再在 *a* 中搜索，如果还是没有找到，那么再去公用的全局变量中搜索。
+ *y* 将只会在 *a* 中进行搜索，因为你已经直接指定了要在当前堆栈的那一个元素中去搜索（*@[1]y*）。

在上面的示例中，符号 *@* 表示堆栈自向，所以，你可以像下面这样的写：

+ *Go代码*: v.Exec(os.Stdout, "Hello", "World!", map[string]string{"cox": "Antu"})
+ *模板中*: $@[1] $@[2] $@[1] $@ $cox!
+ *输出*: Hello World! Hello Antu!

最后，你还可以像下面这样输出整个堆栈以查看它的所有内容：

    $for i, v in @:
        $i: $v<br />
    $end

想了解更多请移步[*Kasia.go文档*](https://github.com/ziutek/kasia.go/blob/master/README.md)。

插了这点小知识之后，我们该回到前面一直在进行中的项目了，让我们来创建最后一个模板 *edit.kt* 文件：

    <form action="/$Id" method="post">
    <h2 class="field">
        <input type="text" name="title" id="title" value="$Title" placeholder="标题：$Title"/>
    </h2>
    <div class="field">
        <textarea name="body" id="body" placehoder="内容：$Body">$Body</textarea>
    </div>
    <div class="actions">
        <input type="submit" value="退出" class="button" />
        <input type="submit" name="submit" value="保存" class="button" />
    </div>
    </form>

我们现在需要一个样式表来让这个维基好看一些，当然这不是必须的，你可以直接使用我下面的这一份样式表：

    body {
    font: 16px/1.62 "Xin Gothic", "Hiragino Sans GB", "Microsoft YaHei", "WenQuanYi Micro Hei", Arial, sans-serif;
    color: #333;
    margin: 0;
    }
    h1, h2, h3, h4, h5, h6, strong, em {
    font-family: "Xin Gothic", "Hiragino Sans GB", "Microsoft YaHei", "WenQuanYi Micro Hei", Arial, sans-serif;
    }
    h1 {
    margin: 0;
    padding: .5em;
    background: #ddd;
    border-bottom: .1em solid #aaa;
    }
    h1 a {
    text-decoration: none;
    color: #333;
    }
    h2 {
    margin: 0;
    padding: .5em;
    }
    h2 input {
    border: .1em solid #aaa;
    font-size: 1em;
    width: 80%;
    }
    div.field {padding: 1em}
    div.field label {
    font-size: 1.5em;
    display: block;
    }
    textarea {
    width: 90%;
    border: .2em solid #aaa;
    font-size: 1.2em;
    line-height: 1.5em;
    height: 10em;
    }
    .columns {
    letter-spacing: -.45em;
    }
    .columns .left, .columns .right {
    display: inline-block;
    letter-spacing: normal;
    min-height: 20em;
    vertical-align: top;
    }
    .columns .left {
    width: 30%;
    border-right: .2em solid #aaa;
    }
    .columns .right {
    width: 69%;
    border-left: .2em solid #aaa;
    }
    .content {
    padding: .5em 1em;
    }
    .button {
    display: inline-block;
    border: .1em solid #ddd;
    background: #eee;
    border-radius: .4em;
    padding: .5em 1em;
    margin: .5em;
    color: #333;
    text-decoration:none;
    }
    .button:hover {
    background: #999;
    }
    .left a.button {
    width: 70%;
    text-align: center;
    }
    .list {
    list-style: none;
    padding: 0 .8em;
    }
    .list a {
    display: inline-block;
    text-decoration: none;
    color: #333;
    }

## 连接到 MySQl 数据库服务器

我们使用 *MyMySQL* 包来连接 MySQl 数据库，先安装它：

    cox@CoxStation:~/go_wiki$ go get github.com/ziutek/mymysql/autorc

现在我们可以为我们的应用写 MySQl 连接器了，创建一个 *mysql.go* 文件，在该文件的第一部分我们会导入一些必须的包，定义一些常量和全局变量：

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

声明之后，MySQL 连接处理器已经可以连接到数据库了，但是我们并不会明显的去连接它。

在我们的应用中，我们将使用 MyMySQl 的 *autorecon* 接口，这是一些不需要连接数据库即可使用的函数集，更重要的是，使用它们，我们将不需要在因为网络原因或者数据库服务器重启导致与数据库连接中断之后重新手动再次连接数据库，它们会帮我们做好这些事情。

下一步我们定义一些 MySQL 错误处理程序：

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

接着再来定义初始化函数：

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

*Register* 方法所注册的命令将在数据库连接创建之后立马执行， *Prepare* 方法则准备好服务器端的声明，因为我们使用了 *mymysql/autorc* 包，所以，当我们在准备第一个服务器端声明时，数据库连接就会被创建。

我们使用 *预先声明* 来代替普通的数据库查询的原因是这样做更加的安全，现在我们不再需要任何的其它函数来过滤用户输入的数据，因为SQL的逻辑与数据已经被完全分离开了，如果不这样做，我们总是很容易被别人进行数据库注入攻击。

下面让我们来创建从数据库中获取数据提供给页面使用的代码：

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

再定义函数用来获取或者更新文章数据：

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

最后的那个函数使用了 MySQL *INSERT ... ON DUPLICATE KEY UPDATE* 查询，它的功能是 *如果ID存在则更新，否则就插入新数据。*

## 控制器

程序的最后一步就是创建一个控制器来与用户进行互动了，让我们创建 *controller.go* 文件：

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

在上面的代码中：

+ *show* 绑定了 *GET* 方法以及 */(.\*)* 这个URL结构，它被用来展示主页视图以及展示用户选择查看的文章。
+ *edit* 绑定了 *GET* 方法以及 */edit/(.\*)* URL 结构，它负责处理渲染编辑页面
+ *update* 绑定了 *POST* 方法以及 */(.\*)* URL结构，它负责更新文章数据到数据库，只有当用户点击了“ *保存* ”按钮才更新数据，如果点击的是 *取消* 那么直接不进行任何处理，当更新完成之后，将页面重定向至刚才更新的文章的展示页面。

## 运行我们的应用

    cox@CoxStation:~/go_wiki/src/go_wiki$ go run *.go

上面这行命令就可以运行我们的应用了，我们可以创建一个 Bash 脚本来启动我们的应用：

    Ccox@CoxStation:~/go_wiki/src/go_wiki$ emacs run.sh

其内容为：

    go run *.go

之后为其添加可执行权限，再运行它即可启动我们的应用：

    cox@CoxStation:~/go_wiki/src/go_wiki$ chmod +x run.sh
    cox@CoxStation:~/go_wiki/src/go_wiki$ ./run.sh 

从 *controlle.go* 代码中可以知道，我们的应用监听的是 *2223* 端口，打开浏览器，地址栏中输入：[http://127.0.0.1:2223](http://127.0.0.1:2223) 即可访问到我们的应用。：

## 获取该示例的代码

你可以通过下面的命令获取到该示例的英文原版：

    git clone git://github.com/ziutek/simple_go_wiki.git

或者直接将其安装到你当前环境下 $GOPATH 所指定的工作空间中：

    go get github.com/ziutek/simple_go_wiki

如果你对我的翻译版本（与原版有些话不同）感兴趣的话，可以使用下面的地址下载：

+ [http://dl.antusoft.com/examples/golang/go-wiki/go-wiki.tar.gz](http://dl.antusoft.com/examples/golang/go-wiki/go-wiki.tar.gz)
+ [直接从本站下载:/examples/go-wiki/go-wiki.tar.gz](/examples/go-wiki/go-wiki.tar.gz)

## 其它框架

本示例中使用的是 http 完成的 controller，你还可以使用如 [web.go](http://www.getwebgo.com/) 或者 [twister](https://github.com/garyburd/twister) 来重新实现。

## 使用 Markdown 格式化文章内容

我们的示例现在的文章内容是直接从数据库中读取的没有任何格式的纯文本，这在网页中显示出来之后就是一大段连续的没有分段分行的纯文本，很不符合我们的阅读，尤其是像本文这种需要清晰的格式化与条理的文章更加无法阅读了，所以，我们还可以选择一些文本格式化工具，比如我一直使用的 [Markdown](http://daringfireball.net/projects/markdown/syntax) ，这需要你安装 [markdown package](https://github.com/knieriem/markdown) 。

要在我们的项目中使用 Markdown，首先安装 Markdown 包：

    cox@CoxStation:~/go_wiki$ go get github.com/knieriem/markdown

然后还需要修改两人个文件 *view.go* 以及 *show.kt* ：

在 *view.go* 中，做下面这样的修改：

    import "github.com/ziutek/kview"

改为：

    import (
        "bufio"
        "bytes"
        "github.com/ziutek/kview"
        "github.com/knieriem/markdown"
    )

然后添加工具函数：

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

再将它添加到 *show.kt* 的全局变量中去：

    main_view.Div("right", kview.New("show.kt", utils))

最后我们需要在 *show.kt* 文件中将 *$Body* 修改为 *$:markdown(Body)* 即可。

## 应用运行截图

![Go Wiki 首页](/examples/go-wiki/front.png)

![Go Wiki 编辑页](/examples/go-wiki/edit.png)
