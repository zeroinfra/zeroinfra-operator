# ZeroInfra Operator


## ZeroInfra Operator CRDs
ZeroInfra Operator CRDs 清单（使用 `zeroinfra.github.com` 域名）

| 管理组             | API Group                              | CRD 名称               | 功能说明                                                                                     |
|--------------------|----------------------------------------|------------------------|--------------------------------------------------------------------------------------------|
| **存储管理**       | `storage.zeroinfra.github.com/v1`      | `ObjectStorage`        | 管理云原生对象存储（如 S3 兼容存储），支持存储桶策略、跨区域复制和生命周期管理。                              |
|                    |                                        | `FileStorage`          | 定义分布式文件存储（如 NFS/CIFS），支持动态容量扩展、访问控制列表和性能分级配置。                             |
|                    |                                        | `BackupPolicy`         | 统一存储备份策略，覆盖对象存储/文件存储/块存储，支持定时快照、增量备份和保留周期管理。                          |
| **数据库管理**     | `database.zeroinfra.github.com/v1`     | `MySQL`                | 管理单实例 MySQL 数据库，支持基础配置、备份恢复和连接池优化。                                                |
|                    |                                        | `MySQLCluster`         | 部署 MySQL 高可用集群（基于 Group Replication），支持自动故障转移、读写分离和滚动升级。                       |
|                    |                                        | `Redis`                | 单节点 Redis 实例管理，支持持久化配置、内存限制和基础监控。                                                  |
|                    |                                        | `RedisCluster`         | 管理 Redis 集群模式，支持分片扩缩容、副本分布策略和故障自愈。                                                |
|                    |                                        | `MongoDB`              | 单实例 MongoDB 部署，支持存储引擎选择、Oplog 大小配置和基础认证。                                            |
|                    |                                        | `MongoReplicaSet`      | 管理 MongoDB 副本集，支持自动选举、延迟节点配置和跨可用区部署。                                              |
|                    |                                        | `PostgreSQL`           | PostgreSQL 单节点或流复制集群，支持 WAL 归档、连接池管理和扩展插件安装。                                     |
| **消息中间件**     | `message.zeroinfra.github.com/v1`      | `KafkaCluster`         | 部署和管理 Kafka 集群，支持 Broker 动态扩缩容、ZooKeeper 集成和性能调优。                                    |
|                    |                                        | `KafkaTopic`           | 定义 Topic 的分区数、副本因子、保留策略，支持 ACL 权限和流量配额控制。                                        |
|                    |                                        | `RocketMQCluster`      | 管理 RocketMQ 集群（Broker/NameServer），支持主从模式、消息轨迹和存储清理策略。                              |
|                    |                                        | `RocketMQTopic`        | 配置 RocketMQ 主题的消息类型（顺序/广播）、消费组管理和延迟消息级别。                                         |
| **AI 管理**        | `ai.zeroinfra.github.com/v1`           | `Model`                | 管理 AI 模型元数据，支持版本控制、模型注册表和部署依赖声明。                                                  |
|                    |                                        | `Dataset`              | 定义训练数据集的生命周期，支持版本管理、数据预处理流水线和存储卷自动挂载。                                     |
|                    |                                        | `Inference`            | 部署模型推理服务，支持自动扩缩容（HPA）、GPU 资源分配和 Prometheus 监控集成。                                |
|                    |                                        | `Training`             | 编排分布式训练任务，支持多框架（PyTorch/TF）、GPU 资源调度和训练过程可视化。                                  |
|                    |                                        | `FineTuning`           | 管理模型微调任务，支持参数冻结策略、增量检查点保存和超参搜索（Hyperopt 集成）。                               |
| **通用资源**       | `common.zeroinfra.github.com/v1`       | `InfraMonitor`         | 统一基础设施监控配置，支持 Prometheus Exporter 自动发现、自定义告警规则和 Grafana 看板模板。                  |
|                    |                                        | `ResourceQuota`        | 定义多维度资源配额（CPU/GPU/存储/网络带宽），支持命名空间级配额和租户资源隔离。                                |
| **批处理任务**     | `batch.zeroinfra.github.com/v1`        | `CronJob`              | 增强型定时任务，支持 Cron 表达式、任务依赖关系、历史记录保留和失败重试策略（支持指数退避）。                     |
|                    |                                        | `Workflow`             | 基于 DAG 的工作流引擎，支持异构任务编排（容器/Serverless）、条件分支和人工审批节点。                           |


## ZeroInfra Operator Projects

```shell
mkdir storage-operator && cd storage-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra/storage-operator
kubebuilder create api --group storage --version v1 --kind ObjectStorage --force
kubebuilder create api --group storage --version v1 --kind FileStorage
kubebuilder create api --group storage --version v1 --kind BackupPolicy

mkdir database-operator && cd database-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra/database-operator
kubebuilder create api --group database --version v1 --kind MySQL
kubebuilder create api --group database --version v1 --kind MySQLCluster
kubebuilder create api --group database --version v1 --kind Redis
kubebuilder create api --group database --version v1 --kind RedisCluster
kubebuilder create api --group database --version v1 --kind MongoDB
kubebuilder create api --group database --version v1 --kind MongoReplicaSet
kubebuilder create api --group database --version v1 --kind PostgreSQL

mkdir message-operator && cd message-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra/message-operator
kubebuilder create api --group message --version v1 --kind KafkaCluster
kubebuilder create api --group message --version v1 --kind KafkaTopic
kubebuilder create api --group message --version v1 --kind RocketMQCluster
kubebuilder create api --group message --version v1 --kind RocketMQTopic

mkdir model-operator && cd model-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra model-operator
kubebuilder create api --group model --version v1 --kind Model
kubebuilder create api --group model --version v1 --kind Dataset
kubebuilder create api --group model --version v1 --kind Inference
kubebuilder create api --group model --version v1 --kind Training
kubebuilder create api --group model --version v1 --kind FineTuning

mkdir common-operator && cd common-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra/common-operator
kubebuilder create api --group common --version v1 --kind InfraMonitor
kubebuilder create api --group common --version v1 --kind ResourceQuota

mkdir batch-operator && cd batch-operator
kubebuilder init --domain zeroinfra.github.com --owner "zeroinfra" --repo github.com/zeroinfra/batch-operator
kubebuilder create api --group batch --version v1 --kind CronJob
kubebuilder create api --group batch --version v1 --kind Workflow
```