# xdb-postgres
对github.com/lib/pq的封装，适配xdb的框架



# 解析说明


## 参数化支持

```sql
@{field} 

如：
select * from table t where t.name = @{name}  
select * from table t where t.name = @{t.name}  
解析结果：
select * from table t where t.name = $1

```



## & 符合链接

将参数进行and链接，如果参数值不存在或者为空将不会生成and条件
```sql
&{field} ， &{t.field}

如： 
select * from table t where t.id = @{id} &{name} 
select * from table t where t.id = @{id} &{t.name} 

解析结果：
select * from table t where t.id = $1 and name = $2 --参数存在
select * from table t where t.id = $1 and t.name = $2 --参数存在
或者
select * from table t where t.id = $1 --参数不存在或者为空,空字符

```

## | 符合链接

将参数进行or链接，如果参数值不存在或者为空将不会生成or条件

```sql
|{field} 

如：  
select * from table t where t.id = @{id} |{name} 
select * from table t where t.id = @{id} |{t.name} 

解析结果：
 select * from table t where t.id = $1 or name = $2 --参数存在
 select * from table t where t.id = $1 or t.name = $2 --参数存在
 或者
select * from table t where t.id = $1 --参数不存在或者为空,空字符

```


## like / not like 支持

```sql
&{like field} ，&{like %field}， &{like field%} ，&{like %field%}
&{notlike field} ，&{notlike %field}， &{notlike field%} ，&{notlike %field%}
&{t.field like property} ，&{t.field like %property}， &{t.field like property%} ，&{t.field like %property%}
&{t.field notlike property} ，&{t.field notlike %property}， &{t.field notlike property%} ，&{t.field notlike %property%}，&{t.field not like %property%}
----(|符号类似)

样例： 
select * from table t where t.id = @{id} &{like name} 
select * from table t where t.id = @{id} &{like %name}
select * from table t where t.id = @{id} &{like name%}
select * from table t where t.id = @{id} &{like %name%}

select * from table t where t.id = @{id} &{t.field like %newname%}
select * from table t where t.id = @{id} &{t.field not like %newname%}

解析结果：
select * from table t where t.id = $1 and name like $2 -- $2
select * from table t where t.id = $1 and name like $2 -- $2数据已完成%的构造
select * from table t where t.id = $1 and name like $2 -- $2数据已完成%的构造
select * from table t where t.id = $1 and name like $2 -- $2数据已完成%的构造

select * from table t where t.id = @p_id and t.field like $2 -- $2数据已完成%的构造
select * from table t where t.id = @p_id and t.field not like $2 -- $2数据已完成%的构造

```



## in/not in 表达式支持

```sql
---注意：in表达式只接受数组切片数据,其他类似直接返回空

&{in field} ,&{in t.field}
&{notin field} ,&{notin t.field}
&{t.field in property}
&{t.field notin property}
&{t.field not in property}
----(|符号类似)


样例： 
select * from table t where t.id = @{id} &{in t.name} 
select * from table t where t.id = @{id} &{t.name in myinputname} 
select * from table t where t.id = @{id} &{t.name not in myinputname} 

解析结果：
select * from table t where t.id = $1 and t.name = any($2)  --name:[1,2,3]
select * from table t where t.id = $1 and t.name = any($2)  --myinputname:["1","2","3"]
select * from table t where t.id = $1 and t.name != all($2) 

``` 



## cte 表达式支持

```sql
WITH indexed_ids AS (
	SELECT id
	FROM unnest(@{ids}::int[]) WITH ORDINALITY AS t(id, idx)
	ORDER BY id  
)
SELECT u.*
FROM demotest u
JOIN indexed_ids i ON u.id = i.id


解析结果：
WITH indexed_ids AS (
    SELECT id
    FROM unnest($1::int[]) WITH ORDINALITY AS t(id, idx)
    ORDER BY id  
)
SELECT u.* 
FROM demotest u
JOIN indexed_ids i ON u.id = i.id


注：
@{ids} = $1 是int数组类型 []int{1,2,3}

``` 
