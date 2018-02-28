package elastic_go

import (
	"testing"
	"reflect"
)

//1.创建客户端
func TestGetClient(t *testing.T) {
	client, ctx, err := GetClient()

	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(client, ctx)
	}

}

//2.检查服务是否开启
func TestPing(t *testing.T) {
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	ok, err := Ping(client, ctx)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ok)
	}
}

//3.获取版本信息
func TestGetVersion(t *testing.T) {
	client, _, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	version, err := GetVersion(client)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(version)
	}
}

//4.创建索引
func TestCreateIndex(t *testing.T) {
	//准备数据
	type User struct {
		Name string `json:"name"`
	}
	var indexName = "ft_test2"
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
	index := Index{Index: indexName, Mapping: mapping}
	//2.拿到客户端实例
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	//3.创建索引
	err = CreateIndex(client, ctx, index)
	if err != nil {
		t.Fatal(err)
	}
}

//4.2 删除索引
func TestDeleteIndex(t *testing.T) {
	//1.准备数据
	index := Index{Index: "ft_test2"}
	//2.拿到客户端实例
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	//删除索引
	err = DeleteIndex(client, ctx, index)
	if err != nil {
		t.Fatal(err)
	}
}

/*
	5.插入数据document 数据格式 `{"name":"ft"}`
	5.1插入一条document
*/
func TestInsertDocument(t *testing.T) {
	//准备数据
	type User struct {
		Name string `json:"name"`
	}
	document1 := Document{Index: "ft_test2", Type: "user", Id: "1", Body: `{"name":"ft1"}`}
	document2 := Document{Index: "ft_test2", Type: "user", Id: "2", Body: `{"name":"ft1"}`}
	document3 := Document{Index: "ft_test2", Type: "user", Id: "3", Body: `{"name":"ft1"}`}
	document4 := Document{Index: "ft_test2", Type: "user", Id: "4", Body: `{"name":"ft1"}`}

	document5 := Document{Index: "ft_test2", Type: "user", Id: "5", Body: User{Name: "ft2"}}
	//获取客户端
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	//插入记录
	err = InsertDocument(client, ctx, document1)
	if err != nil {
		t.Fatal(err)
	}
	err = InsertDocument(client, ctx, document2)
	if err != nil {
		t.Fatal(err)
	}
	err = InsertDocument(client, ctx, document3)
	if err != nil {
		t.Fatal(err)
	}
	err = InsertDocument(client, ctx, document4)
	if err != nil {
		t.Fatal(err)
	}
	err = InsertDocument(client, ctx, document5)
	if err != nil {
		t.Fatal(err)
	}
}
//5.2 获取Document via Id
func TestGetDocument(t *testing.T) {
	document:=Document{Index:"ft_test",Type:"user",Id:"1"}
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	result,err:=GetDocument(client,ctx,document)
	if err!=nil {
		t.Fatal(err)
	}else{
		t.Log(result)
	}
}

//5.3删除Document
func TestDeleteDocument(t *testing.T) {
	document:=Document{Index:"ft_test",Type:"user",Id:"1"}
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	err =DeleteDocument(client,ctx,document)
	if err !=nil {
		t.Fatal(err)
	}
}

//5.4 查询Documents
func TestSearchDocuments(t *testing.T) {
	//准备数据
	type User struct {
		Name string `json:"name"`
	}
	//
	termQuery := &TermSearch{
		ElemType:reflect.TypeOf(User{}),
		Query: QueryStruct{Key:"name",Value:"ft1"},
		Index:"ft_test2",
		Type:"user",
		SortField:"name",
		Asc:true,
		StartIndex:0,
		QuerySize:5,
	}
	t.Log(termQuery.ElemType)
	//创建客户端
	client, ctx, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}
	//执行Search
	var results []interface{}
	results,err = SearchDocuments(client,ctx,termQuery)
	if err!=nil {
		t.Fatal(err)
	}else{
		t.Log(results,len(results))
	}
}