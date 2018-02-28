package elastic_go

import "reflect"

type Index struct{
	Index string //索引名
	Mapping string //mapping样式

}

type Type struct{
	Type string
}

type Document struct{
	Index string //索引名
	Type string //表名
	Id string //id号
	Body interface{} //bodyString or bodyJson
}

type TermSearch struct{
	ElemType reflect.Type //使用Each()方法需要传递获取对象的类型，比如reflect.Typeof(User{})
	Query QueryStruct //精准查询term查询的参数比如 name:"ft",满足属性name值为ft的查询条件
	Index string //索引名
	Type string //表名
	SortField string //按照该字段排序
	Asc bool //升序true，降序false
	StartIndex int //起始位置0开始
	QuerySize int //查询多少条数据
}
type QueryStruct struct{
	Key string
	Value string
}