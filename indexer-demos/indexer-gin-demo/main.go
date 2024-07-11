package main

import (
	"client-go-test/indexer-demos/indexer-gin-demo/basic"
	"client-go-test/indexer-demos/indexer-gin-demo/remote"
	"github.com/gin-gonic/gin"
)

/*
	nginx	c		web

	mysql	c		web

	tomcat	java	storage
*/

func main() {
	r := gin.Default()

	// kubernetes相关的初始化操作
	basic.DoInit()

	// 用于提供基本功能的路由组
	basicGroup := r.Group("/basic")

	// a. 查询指定语言的所有对象的key(演示2. IndexKeys方法)

	/*
			{
		  		"language": [
		  		  "indexer-tutorials/mysql-556b999fd8-22hqh",
		  		  "indexer-tutorials/nginx-deployment-696cc4bc86-2rqcg",
		          "indexer-tutorials/nginx-deployment-696cc4bc86-bkplx",
		  		  "indexer-tutorials/nginx-deployment-696cc4bc86-m7wwh"
		  		]
			}
	*/

	basicGroup.GET("get_obj_keys_by_language_name", basic.GetObjKeysByLanguageName)

	// b. 返回对象的key，返回对应的对象(演示Store.GetByKey方法)
	basicGroup.GET("get_obj_by_obj_key", basic.GetObjByObjKey)

	// c. 查询指定语言的所有对象(演示4. ByIndex方法)
	basicGroup.GET("get_obj_by_language_name", basic.GetObjByLanguageName)

	// d. 根据某个对象的key，获取同语言类型的所有对象(演示1. Index方法)
	basicGroup.GET("get_all_obj_by_one_name", basic.GetAllObjByOneName)

	// e. 返回所有语言类型(演示3. ListIndexFuncValues方法)
	basicGroup.GET("get_all_language", basic.GetAllLanguange)

	// f. 返回所有分类方式，这里应该是按服务类型和按语言类型两种(演示5. GetIndexers方法)
	basicGroup.GET("get_all_class_type", basic.GetAllClassType)

	remoteGroup := r.Group("/remote")
	// g. 使用clientset远程查询
	remoteGroup.GET("get_obj_by_obj_key_remote_query", remote.GetObjByObjKey)

	r.Run(":18080")
}
