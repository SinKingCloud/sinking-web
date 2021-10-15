package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"time"
)

func checkCluster() {
	go func() {
		for {
			//检测集群状态
			clusterList := service.CopyClusters()
			for k := range clusterList {
				if clusterList[k].LastHeartTime+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
					service.ClustersLock.Lock()
					service.Clusters[k].Status = 1
					service.ClustersLock.Unlock()
				}
			}
			//检测服务状态
			serviceList := service.CopyService()
			for k, v := range serviceList {
				for k1, v1 := range v {
					for k2, v2 := range v1 {
						for k3, v3 := range v2 {
							for k4, v4 := range v3 {
								if v4.LastHeartTime+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
									service.ServicesLock.Lock()
									service.Services[k][k1][k2][k3][k4].Status = 1
									service.ServicesLock.Unlock()
									service.LocalServicesLock.Lock()
									if service.LocalServices[k][k1][k2][k3][k4] != nil {
										service.LocalServices[k][k1][k2][k3][k4].Status = 1
									}
									service.LocalServicesLock.Unlock()
								}
							}
						}
					}
				}
			}
			time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
		}
	}()
}
