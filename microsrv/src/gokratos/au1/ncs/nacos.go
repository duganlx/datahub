package ncs

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	nacosIp          = "192.168.15.42"
	nacosPort        = uint64(8848)
	nacosNamespaceId = "ad13b472-04a0-4cf5-a4ee-d8cfd4cf81f5"
)

func NewNacosConfigClient() (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacosIp, nacosPort),
	}

	cc := &constant.ClientConfig{
		NamespaceId: nacosNamespaceId,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  cc,
		},
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewNacosInamingClient() (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacosIp, nacosPort),
	}

	cc := &constant.ClientConfig{
		NamespaceId: nacosNamespaceId,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  cc,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
