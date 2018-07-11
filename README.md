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

```sql
CREATE TABLE IF NOT EXISTS cs_geog (
  id SERIAL PRIMARY KEY,
  osm_id BIGINT,
  name TEXT,
  brand TEXT,
  geog GEOGRAPHY(POINT, 4326)
);

CREATE INDEX ON cs_geog (brand);

CREATE INDEX ON cs_geog USING GiST (geog);
```

```sql
INSERT INTO cs_geog(osm_id, name, brand, geog)
  SELECT id, name, brand,
    ST_GeographyFromText(format('SRID=4326;POINT(%s %s)', lon, lat))
  FROM osm_nodes;
```

### query

Dist version.

```sql
SELECT id, osm_id, name, brand, ST_Distance('SRID=4326;POINT(139.8265316 35.7140523)'::GEOGRAPHY, geog) FROM cs_geog ORDER BY ST_Distance('SRID=4326;POINT(139.8265316 35.7140523)'::GEOGRAPHY, geog) LIMIT 5;
```

KNN GiST version.

```sql
SELECT id, osm_id, name, brand, ST_Distance('SRID=4326;POINT(139.8265316 35.7140523)'::GEOGRAPHY, geog) FROM cs_geog ORDER BY 'SRID=4326;POINT(139.8265316 35.7140523)'::GEOGRAPHY <-> geog LIMIT 5;
```
