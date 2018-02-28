package elastic_go

import (
	"github.com/gin-gonic/gin"
)

type Response struct{
	Code int	    //返回的指示码，400表示可行，其它自定义
	Message string  //返回消息，错误则返回错误原因与改善建议
	Body interface{} //返回的数据对象，只有再code为400时并且需要回馈对象时才有值
}
func main(){
	router:=gin.Default()
	router.GET("v1/x-service/GET/",X_Func)
}
func X_Func (c *gin.Context){
	//TODO
}