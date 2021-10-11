package service

import "github.com/SinKingCloud/sinking-go/sinking-consul/app/model"

// Clusters 集群列表
var Clusters = make(map[string]Cluster)

// RegisterClusters 注册集群列表
var RegisterClusters = make(map[string]Cluster)

// Cluster 集群信息结构
type Cluster struct {
	Hash          string         `json:"hash"`            //标识hash
	Ip            string         `json:"ip"`              //集群ip
	Port          string         `json:"port"`            //集群端口
	LastHeartTime model.DateTime `json:"last_heart_time"` //上次心跳时间
	Status        int            `json:"status"`          //集群状态(0:正常/1:异常)
}
