## 作业
### 题目
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

### 思考
首先，sql.ErrNoRows 跟 io.EOF 有点相似，更像是一种功能性的“提醒”，而不是一个足以中断程序运行的“错误”。

这个问题其实跟业务功能的关联性也比较大。我认为在大部分的场景中，查询到空的结果集不应该被认为是一种错误，即空结果也是结果。

    例如：查找一个 id 编号为 x 的员工，查询不到该员工也应该是一个合理的结果。

从 dao 层开发者的角度看。既然空集是一种结果，那么就应该在 dao 的 handle 范围以内，因此在 dao 层中遇到该错误时，我更倾向于对该错误进行单独判断后返回`nil`或者空数组/struct，而不是 Wrap 这个 error，抛给上层。

在这里，到底是返回`nil`还是空数组/struct 也是一个值得讨论的话题，最起码是要在团队内部形成统一的规范

### 代码
```go
func GetEmployee(id int) (*string, error) {
    querySQL := fmt.Sprintf("select name from employees where id = %d", id)
    var name string
    err := db.QueryRow(querySQL).Scan(&name)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        } else {
            return nil, errors.Wrap(&QueryError{err, querySQL}, "fail to Query employee")
        }
    }
    return &name, nil
}
```
