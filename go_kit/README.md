# 注意事项
此工具库是甘来日常使用的工具库，部分源码来自github，由大龙主导，不涉及商业机密，日志格式使用beego

1、tasktimer中的文件仅是演示基于beego的定时任务的使用示例。  

2、生产环境项目编译打包时,要明确自己的项目运行的什么系统环境下,一般在linux环境,编译打包命令如下:
    cd 目录到自己项目下,执行命令: env GOOS=linux GOARCH=amd64 go build  

3、每个项目必须包含git忽略文件: .gitignore。  

4、API接口请求参数在多于3个以上,必须定义结构体。  

5、返回的所有error,必须是logiccode类型的,每层的业务错误必须通过logiccode封装,并在发生错误的地方用logkit输出真正的程序错误。  

# 工程目录结构 
* conf:           工程配置文件,必须以.conf文件结尾,文件名称用小写字母命名,不允许出现其他特殊字符。  


* controllers     控制层,文件名格式为XXXController,XXX首字母大写,驼峰式命名。  


* dao             DAO层,文件名格式为XXXDao,XXX首字母大写,驼峰式命名。  


* logic           业务逻辑层,文件名格式为XXXLogic,XXX首字母大写,驼峰式命名。  


* models          实体模型层,非特殊情况只允许定义一个model.go文件,所有的实体都定义在这个文件中。  


* routers         路由层,非特殊情况只允许定义一个router.go文件,所有的API一级目录映射都在该文件中定义。  


* test            自测文件,文件名格式XXX_test,XXX首字母大写,驼峰式命名。

* task            定时任务,文件名格式为XXXTask,XXX首字母大写,驼峰式命名。

* vendor          vendor文件,第三方依赖组件目录,由govendor自动生成。

* .gitignore      git忽略文件。

* xxx.go          main函数文件,文件名与项目名称一致,全小写。

# 工具库目录说明 
-  dbkit:      常规DB增删改查操作的工具库。  


- logiccode:  通过业务错误码,以及自定义错误码类。  


- logkit:     对beego的日志模块的封装。  


- strkit:     字符串的校验、类型转换工具库。  


- tasktimer:  定时任务DEMO。  


- timekit:    日期格式化工具库。  


- vendor:     依赖的第三方库。  

# 常规错误 
- Column count doesn't match value count at row   
    - SQL更新的字段与填充的"?"数量不一致时,会报此错误。