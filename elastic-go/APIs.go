package elastic_go

import (
	"github.com/olivere/elastic"
	"context"
	"github.com/pkg/errors"
)

/*
	1.获取一个客户端，成功获取后会返回一个上下文context和client实例，该client是默认开启状态
 */
func GetClient() (*elastic.Client, context.Context, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(URL),
		elastic.SetScheme(Scheme),
		elastic.SetHealthcheck(HeathCheck),
		elastic.SetSendGetBodyAs(Method),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read failed")
	}
	return client, context.Background(), nil
}

/*
	2.检查该客户端实例能否ping通服务，判断es服务有没有正常开启
*/
func Ping(client *elastic.Client, ctx context.Context) (bool, error) {
	_, _, err := client.Ping(URL).Do(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

/*
	3.获取elasticsearch的版本信息
 */
func GetVersion(client *elastic.Client) (string, error) {
	version, err := client.ElasticsearchVersion(URL)
	if err != nil {
		return "", err
	} else {
		return version, nil
	}
}

/*
	4.index索引操作,Index类型定义在DTO.go里，由Mapping和indexName构成
	4.1 增加索引
*/
func CreateIndex(client *elastic.Client, ctx context.Context, index Index) error {
	exists, err := client.IndexExists(index.Index).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err := client.CreateIndex(index.Index).BodyString(index.Mapping).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		return errors.New("该index已存在")
	}
	return nil
}

/*
	4.2删除索引
 */
func DeleteIndex(client *elastic.Client, ctx context.Context, index Index) error {
	_, err := client.DeleteIndex(index.Index).Do(ctx)
	if err != nil {
		return err
	} else {
		return nil
	}

}

/*
	5.document操作
	5.1插入Insert
*/
func InsertDocument(client *elastic.Client, ctx context.Context, document Document) error {

	if v, ok := document.Body.(string); ok {
		_, err := client.Index().Index(document.Index).Type(document.Type).Id(document.Id).BodyString(v).Do(ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		_, err := client.Index().Index(document.Index).Type(document.Type).Id(document.Id).BodyJson(document.Body).Do(ctx)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

//5.2获取数据
//Get()方法的具体数据还没实现，无法取出
func GetDocument(client *elastic.Client, ctx context.Context, document Document) (*elastic.GetResult, error) {
	result, err := client.Get().Index(document.Index).Type(document.Type).Id(document.Id).Do(ctx)
	if err != nil {
		return nil, err
	} else {
		_, err = client.Flush().Index(document.Index).Do(ctx)
		if err != nil {
			panic(err)
		}
		return result, nil
	}

}

//5.3获取数据组
func SearchDocuments(client *elastic.Client, ctx context.Context, termQuery *TermSearch) ([]interface{},error){

	tq := elastic.NewTermQuery(termQuery.Query.Key, termQuery.Query.Value)
	searchResult, err := client.Search().
		Index(termQuery.Index). // 指定index，返回一个*SearchService对象
		Query(tq). // 设置查询体，返回同一个*SearchService对象
		Sort(termQuery.SortField, termQuery.Asc). // 按照user升序排列
		From(termQuery.StartIndex).Size(termQuery.QuerySize). // 从第一条数据，找十条，即0-9
		Pretty(true). // 使查询request和返回的结果格式美观
		Do(ctx)
	if err != nil {
		return nil,err
	}
	return searchResult.Each(termQuery.ElemType),nil
}

//5.4删除数据
func DeleteDocument(client *elastic.Client, ctx context.Context, document Document) error {
	_, err := client.Delete().Index(document.Index).Type(document.Type).Id(document.Id).Do(ctx)
	if err != nil {
		return err
	} else {
		return nil
	}
}
