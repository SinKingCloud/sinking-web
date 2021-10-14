package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ClusterRegister 注册集群
func ClusterRegister(s *sinking_web.Context) {
	type register struct {
		Ip   string `form:"ip" json:"ip"`
		Port string `form:"port" json:"port"`
	}
	cluster := &register{}
	err := s.BindJson(&cluster)
	if err != nil || cluster.Ip == "" || cluster.Port == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	service.ClustersRegister(cluster.Ip, cluster.Port)
	response.Success(s, "注册集群成功", nil)
}
