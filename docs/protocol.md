# protocol

自定义与嵌入式设备沟通的协议

## 1. Header格式

```
+------16bit-----+

+----------------+
|   version(16)  |
+----------------+
|     cmd(16)    |
+----------------+
|  payload size  |
|      (32)      |
+----------------+
```

- version 固定为 `1`
- cmd 具体见 payload格式
- payload size 表示payload的大小

**注意：网络字节序为大端排序！！！**

## 2. payload格式

通过不同的cmd代表不同种类的数据包。对应的payload不同。

### 2.1 客户端发送心跳（cmd = 1）

此类packet无需payload

### 2.2 客户端发送设备信息（cmd = 2）

- payload size 固定值为 4
- 发送设备ID（设备ID为32bit无符号整数）

### 2.3 客户端发送歌曲信息（cmd = 3）

循环发送：
- 歌曲ID(32 bit)
- 歌曲名长度(16 bit)
- 歌手名长度(16 bit)
- 歌曲名
- 歌手名

### 2.4 客户端发送退出信息（cmd = 4）

收到后，服务端将主动断开连接
此类packet无需payload

### 2.5 客户端发送歌曲播放完毕信息（cmd = 5）

收到后，服务端将发送下一首歌（参考2.6）
此类packet无需payload

### 2.6 服务端发送 播放/切换 歌曲命令（cmd = 1000）

- payload size 固定值为 4
- payload为歌曲ID（歌曲ID为32bit无符号整数）
