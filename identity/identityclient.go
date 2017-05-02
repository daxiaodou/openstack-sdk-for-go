package identity

import (
	. "git.intra.weibo.com/golang/openstack-sdk-for-go/http"
	. "git.intra.weibo.com/golang/openstack-sdk-for-go/util"
)

type OskIdentityClient struct {
	CommonParam RequestParam
	Instance    *OskIdentityInstance
}

func NewIdentityClient(paramsMap map[string]interface{}) (client *OskIdentityClient) {
	client = &OskIdentityClient{}
	client.CommonParam.Attr = paramsMap
	client.CommonParam.Params = ""
	client.CommonParam.Content_Type = JSON_CONTENT_TYPE

	client.Instance = &OskIdentityInstance{client.CommonParam}
	return client
}
