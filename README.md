## Log

```console
$ pgc -datadir db/001 initdb
$ pgc -datadir db/001 start
```

DB start to work on:

```
postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable
```

<http://overpass-api.de/query_form.html>

```xml
<query type="node">
  <!-- ボックスを指定 -->
  <bbox-query s="35.572" n="35.786" w="139.663" e="139.883" />
  <!-- コンビニエンスストアのみ抜き出す -->
  <has-kv k="shop" v="convenience" />
</query>
<print/>
```

Save as `.osm`
