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