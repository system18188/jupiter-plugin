# binding

golang 提交参数验证绑定包

```
type form struct {
    Page     binding.Int64 `json:"page" form:"page,default=1" binding:"required,gte=1,number" description:"页" default:"1"`
    PerPage  binding.Int64 `json:"perPage" form:"perPage,default=60" required:"omitempty,number,gte=10,lt=120" default:"60" description:"每个页面条数"`
    Sorter  binding.Int32 `json:"sorter" form:"sorter" binding:"omitempty,oneof=0 1" description:"排列 0 时间排序；1 其他排序" default:"0"`
}

valid := binding.Default(req.Request.Method, req.HeaderParameter(restful.HEADER_ContentType))
if err := valid.Bind(req.Request, form); err != nil {
    ...
    return
}
```


## 标签说明
```
nohtml ：验证不包括HTML代码

column : 验证数据库字段名

isDate ： 验证年月日格式：0000-00-00

isTime ： 验证时间格式：0000-00-00 00:00:00

required ：验证唯一

number ： 验证是否是数字

gte=1 ：验证数据大于等于一

lt=100 ： 验证数字不小于100

lt=100 ： 验证数字不小于100

omitempty ： 验证数据可为空

oneof=0 1 ： 验证数据只能由 0 或 1 组成


```
