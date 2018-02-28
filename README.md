# elastic-go-tool
一个简易的用go语言开发的elastic api 工具包,底部引用了(http://github.com/olivere/elastic)

**使用条件**
elasticsearch版本选择的是6.2.2,向下兼容，即可以使用6.x以前的版本.
下载地址:(http://https://www.elastic.co/cn/products/elasticsearch)
**注意**
1. elasticsearch要求jdk1.8以上的java版本，请确保安装正确的jdk，jdk配置自行百度.
2. elasticsearch版本6.x有几个改变的地方一定要注意，首先是同一个index不能有多type了，可以形象的理解为，存储User的库只能用来存放users，而不能同时存user和school，单index只能有单type
3. elasticsearch5.x以后就移除了string类型，也就是go里的string，映射到es里是text或者keyword,如:
```go
type User struct {
		Id   int `json:"id"`
		Name string `json:"name"` //下面的name字段对应的type应该用text而不是string，用了string会报类型转换的错误
	}
var mapping = `
	{
		"settings":{
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings":{
			"user":{
				"properties":{
					"name":{
						"type":"keyword"
					}
				}
			}
		}
	}`
```
4.使用InsertDocument方法，如果不存在index，会自动创建该index，mapping样式是自适应的,由es服务控制
5.使用GetDocument方法获取的对象，无法直接拿到具体的值，只能拿到它的索引信息，避免使用它，使用SearchDocuments作为替代方案
6.插入和修改使用一致，插入已经存在的id就会自动默认为update
7.使用SearchDocuments()时，Query的关键字段必须是keyword，而不能是text！！！！！！！！！！！不然就会报:
**go type=search_phase_execution_exception**

**版本**
v1版本还略微简陋，后续将更新以下:
1.将配置方式以toml的形式配置
2.添加自定义错误
3.添加事务模块
***

**Example**
```go
package main

import (
	t "ESAPI_Test/elastic-go"
	"fmt"
	"reflect"
)

type User struct {
	Name string `json:name`
}

func main() {
	//创建客户端
	client, ctx, err := t.GetClient()
	if err != nil {
		panic(err)
	}

	//创建索引
	indexName := "test_elastic"
	mapping := `
	{
		"settings":{
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings":{
			"user":{
				"properties":{
					"name":{
						"type":"keyword"
					}
				}
			}
		}
	}`
	index := t.Index{Index: indexName, Mapping: mapping}
	err = t.CreateIndex(client, ctx, index)
	if err != nil {
		panic(err)
	}

	//插入数据
	document1 := t.Document{Index: "test_elastic", Type: "user", Id: "1", Body: `{"name":"ft1"}`}
	document2 := t.Document{Index: "test_elastic", Type: "user", Id: "2", Body: `{"name":"ft1"}`}
	document3 := t.Document{Index: "test_elastic", Type: "user", Id: "3", Body: `{"name":"ft1"}`}
	document4 := t.Document{Index: "test_elastic", Type: "user", Id: "4", Body: `{"name":"ft1"}`}

	document5 := t.Document{Index: "test_elastic", Type: "user", Id: "5", Body: User{Name: "ft2"}}
	//插入记录
	err = t.InsertDocument(client, ctx, document1)
	if err != nil {
		panic(err)
	}
	err = t.InsertDocument(client, ctx, document2)
	if err != nil {
		panic(err)
	}
	err = t.InsertDocument(client, ctx, document3)
	if err != nil {
		panic(err)
	}
	err = t.InsertDocument(client, ctx, document4)
	if err != nil {
		panic(err)
	}
	err = t.InsertDocument(client, ctx, document5)
	if err != nil {
		panic(err)
	}

	//通过id获取数据
	document := t.Document{Index: "test_elastic", Type: "user", Id: "1"}
	result, err := t.GetDocument(client, ctx, document)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(result)
	}

	//删除一条document
	err = t.DeleteDocument(client, ctx, document)
	if err != nil {
		panic(err)
	}

	//查询documents
	termQuery := &t.TermSearch{
		ElemType:   reflect.TypeOf(User{}),
		Query:      t.QueryStruct{Key: "name", Value: "ft1"},
		Index:      "test_elastic",
		Type:       "user",
		SortField:  "name",
		Asc:        true,
		StartIndex: 0,
		QuerySize:  5,
	}
	//执行Search
	var results []interface{}
	results, err = t.SearchDocuments(client, ctx, termQuery)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(results, len(results))
	}

	//删除索引
	index = t.Index{Index: "test_elastic"}

	err = t.DeleteIndex(client, ctx, index)
	if err != nil {
		panic(err)
	}
}

```