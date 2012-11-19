<a href="/edit/" class="button">创建新文章</a>
<hr />
<h2>最近更新</h2>
<ul class="list">
$for _, article in Articles:
    <li><a href="$article[Id]">$article[Title]</a></li>
$end
</ul>