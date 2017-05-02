package compute

import (
	. "git.intra.weibo.com/golang/openstack-sdk-for-go/http"
	. "git.intra.weibo.com/golang/openstack-sdk-for-go/identity"
	. "git.intra.weibo.com/golang/openstack-sdk-for-go/util"
)

type OskComputeClient struct {
	CommonParam RequestParam
	Instance    *OskComputeInstance
}

func NewComputeClient(paramsMap map[string]interface{}) (client *OskComputeClient) {
	client = &OskComputeClient{}
	client.CommonParam.Attr = paramsMap
	client.CommonParam.Params = ""
	client.CommonParam.Content_Type = JSON_CONTENT_TYPE

	identitycli := NewIdentityClient(paramsMap)
	identityresp := identitycli.Instance.GetToken()
	if identityresp == nil || identityresp.Err != nil {
		return nil
	}
	client.CommonParam.AccessToken = identityresp.Token

	client.Instance = &OskComputeInstance{client.CommonParam}
	return client
}

func (this OskComputeClient) SetCommonParamAttr(paramsMap map[string]interface{}) {
	this.CommonParam.Attr = paramsMap
	this.Instance.CommonParam.Attr = paramsMap
}
