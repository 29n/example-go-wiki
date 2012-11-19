<div class="article">
$if Id:
    <h2>$Title</h2>
    <div class="content">
        $:markdown(Body)
    </div>
    <p class="actions">
        <a href="/edit/$Id" class="edit button">编辑本文</a>
    </p>
$else:
    <h2>维其示例页面</h2>
    <div class="content">
        <p>您现在所看到是本维基的示例页面，这说明您现在还没有添加任何内容。</p>
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
