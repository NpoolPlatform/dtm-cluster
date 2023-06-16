# [DTM](https://www.dtm.pub/)

[Distributed transaction deply](https://www.dtm.pub/deploy/base.html)

## 优化点

+ 库表结构初始化
+ 配置文件使用 **configmap**

## 依赖
- 当mysql服务重启后，需要重启dtm服务
- 当dtm服务重启后，需要重启appuser，basal，review等依赖dtm的服务
