# 1、basics

golang基础库，常用的功能

1. redis扩展功能，集群锁、资源控制器等
2. set工具包含：add、has、del，交、差、并计算等
3. uuid一系列功能
4. 携程池
5. 通用util功能集合
6. 简单易用的enum

## 2、enum

```go
// int类型
e := New[struct {
First  Elem[int] `code:"1" text:"123"`
Second Elem[int] `code:"2" text:"123"`
}]()

fmt.Println(Code2Text[int](e))
fmt.Println(e.First.Code, e.First.Text)

// string类型
e := New[struct {
First  Elem[string] `code:"1" text:"123"`
Second Elem[string] `code:"2" text:"123"`
}]()

fmt.Println(Code2Text[string](e))
fmt.Println(e.First.Code, e.First.Text)
```

## 3、redis

[示例代码](./redis/redis_simple.go)

```go
// 集群锁
opts := DefaultOptions()
opts.Addr = "127.0.0.1:6379"
opts.DB = 0

r := NewRedis(&opts)
defer r.Close()

lockKey := "test_key"
r.WithLock(lockKey, nil, 10, 3, 500, func () {
fmt.Println("ok")
})
```

[所有的扩展](./redis/ext.go)

## 4、set集合

```go
s1 := New[string]()
s2 := New[string]()

s1.Add("1", "2", "3")
s2.Add("3", "5", "6")

fmt.Println("Union: ", s1.Union(s2).ToList())
fmt.Println("Intersection: ", s1.Intersection(s2).ToList())
fmt.Println("Difference: ", s1.Difference(s2).ToList())
fmt.Println("SymmetricDifference: ", s1.SymmetricDifference(s2).ToList())
```

## 5、uuid

```go
uu, err := NewV4()
if err != nil {
fmt.Println(err.Error())
}

fmt.Println(uu.String())
```

## 6、workpool

```go
wp := NewWorkPool(100, "test", 0, 50)
for i := 0; i < 100; i++ {
wp.SubmitJob(&JobBag{
JobFunc: func (i ...interface{}) {
if i[0].(int)%2 == 0 {
// panic(fmt.Errorf("1111111111111"))
}
fmt.Println("------------", i[0])
// time.Sleep(time.Second * 1)

},
Params: []interface{}{i},
})
}

for {
if wp.IsFinished() {
break
}
time.Sleep(time.Millisecond * 10)
}
wp.ShutDownPool()
fmt.Println("done")
```

## 7、util

[各种常用的方法](./util/util.go)
