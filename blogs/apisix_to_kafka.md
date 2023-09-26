# APISIX 配置 kafka 连接

## 实验操作

### Offset Explorer 配置

使用的 Kafka Topics 为下面几个

- rms_heartbeat: 时间字符串 格式为 `{时:分:秒}` 如 `{11:37:12}`
- rms_order: json 对象，包含的字段有 create_date(年月日组成的数字, 如 20230922), order_id, order_time(11 位时间戳), rsp_time(11 位时间戳), fk_time(11 位时间戳), account, market, code, biz, price, vol, matchprice, matchvol, eam_id, status, msg
- rms_reject_cancel_order: 同 rms_order
- rms_reject_order: 同 rms_order
