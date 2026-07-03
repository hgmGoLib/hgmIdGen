# hgmIdGen

一个零依赖的 Go id 生成库.通过 `时间 + 进程内递增计数 + 强随机` 组合,生成大概率递增、低碰撞、只含安全字符的 id.

```
go get github.com/hgmGoLib/hgmIdGen
```

```go
import "github.com/hgmGoLib/hgmIdGen"
```

## 能干啥

提供 3 种输出,按用途分别选用:

| 场景 | 调用 | 结果例子 | 长度 |
|------|------|----------|------|
| 一般数据库主键 (userId/orderId/postId) | `hgmIdGen.NewId()` | `17fatw9hnw1y8ge6emw3et3fmtu9d` | 固定 29 字符 |
| 高安全 token (sessionId/accessToken) | `hgmIdGen.NewSecureId()` | `17fatweptd1swfgg11ekgtx5j459wnf112v7jp81wn` | 固定 42 字符 |
| 二进制 id (省内存/自定义存储) | `hgmIdGen.NewIdBinary(dst)` | — | 固定 18 字节 |

两种字符串方案都保证:结果只含不易混淆的安全字符、同进程连续调用大概率递增、时间+递增+随机组合避免简单碰撞.

## 怎么干

### 生成 id

```go
id := hgmIdGen.NewId()             // 数据库主键
token := hgmIdGen.NewSecureId()    // 高安全 token

var buf []byte
buf = hgmIdGen.NewIdBinary(buf)    // 18 字节二进制,dst 传 nil 则新分配
```

### 二进制与字符串互转

`NewIdBinary` 与 `NewId` 是同一套方案,只差 base32 编码,可无损互转:

```go
str, err := hgmIdGen.IdBinaryToId(bin)   // 18 字节 -> 29 字符
bin, err := hgmIdGen.IdToIdBinary(str)   // 29 字符 -> 18 字节
```

### 校验外部传入的字符串

```go
hgmIdGen.IsValidId(s)        // 是否是合法的 NewId() 格式
hgmIdGen.IsValidSecureId(s)  // 是否是合法的 NewSecureId() 格式
```

### 从 id 反查生成时间

```go
hgmIdGen.ParseTimeFromId(id)         // NewId() 结果 -> 毫秒时间
hgmIdGen.ParseTimeFromSecureId(id)   // NewSecureId() 结果 -> 毫秒时间
解析失败统一返回 time.Time{}.
```

### 按时间范围扫描的下界

需要 `where id >= ?` 这种按时间范围扫描时,用对应的最小值构造:

```go
hgmIdGen.MinIdAtTime(t)         // 某毫秒时刻的最小 NewId()
hgmIdGen.MinSecureIdAtTime(t)   // 某毫秒时刻的最小 NewSecureId()
```

## 选型注意

* `NewId()` 适合数据库主键,没人故意 ddos 的情况下极难碰撞.
* `NewSecureId()` 随机段更长,即使有人故意撞库也很难碰撞,用于 sessionId / accessToken.

## 实现细节

* **`NewId()` (18 字节)**: 6 字节毫秒时间 (约可表示 8925 年) + 1\~4 字节变长递增计数 (前 2 位表示字节数,承载 6\~30 bit,最大 `2^30-1`,时间变化即归零) + 8\~11 字节强随机,最后整体 base32 编码为 29 字符.
* **`NewSecureId()` (26 字节)**: 结构同上,但随机段扩到 16\~19 字节,抗碰撞更强,base32 编码为 42 字符.
* **base32 字符表**: `123456789abcdefghjkmnpqrstuvwxyz`,去掉易混淆字符,且按 ascii 递增,保证编码结果可直接字符串排序.

## 设计须知与极限情况

选型、排错、或依赖某个"看起来成立"的性质(排序、递增、抗碰撞)之前,请先读这两份文档:

* [`doc/implDetail.md`](doc/implDetail.md) — 逐字段二进制布局:三段式结构(6 字节时间 + 1\~4 字节变长递增 + 剩余强随机)、变长编码规则、base32 长度换算、各函数长度用途速查.
* [`doc/edgeCase.md`](doc/edgeCase.md) — 设计前提与边界情况:递增计数器每毫秒 `2^30` 上限与溢出回绕、变长递增挤占随机段导致熵随负载下降、时钟回拨破坏递增/排序、排序保证的成立前提、`Min*AtTime` 只能当下界用、三函数共享同一全局递增计数器、randFast 随机源性质、时间上限、校验函数只做格式校验等.

* [`doc/faq.md`](doc/faq.md) — 对一次设计 review 的逐条回复:时间精度/counter 取舍、IsValidId 语义与 canonical 边界、命名风格、格式长期稳定性、纯随机 token 缺口等设计取舍的说明.

也可直接阅读对应源码文件 (`Id_*.go` / `SecureId_*.go` / `IdBinary_*.go`).

## License

[Unlicense](LICENSE) (public domain).
