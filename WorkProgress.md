# 工作进度

---
## 2021-9-8

### 计划任务

- 抽象出一个 packer 类用于解包
- 完成TCP收发包流程

---

## 2021-9-9

- 使用单独的包进行解包
- 暂定使用map当数据
- 暂定数据头格式:
- (2byte 总长) (1byte 总包数量) (1byte 序号) (数据)
- 暂定数据格式：
- (4byte Seq) (str Header) (str Data)
- str 格式为 (2byte 长度)([]byte 格式字符串)
- Seq 为包序号 uint32
- Header 为数据头部类 struct
- Data 为数据 map[string]interface{}

## 2022-3-17
- 准备完全重构
- 增加文档 doc.md


## 2022-3-19
- 完成 TCP Listen
- 优化 Unpacker
- 未来准备增加 UDP Listen

### 未完待续