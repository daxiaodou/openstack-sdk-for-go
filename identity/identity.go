package identity

import (
	"encoding/json"
	"fmt"

	. "git.intra.weibo.com/golang/openstack-sdk-for-go/http"
	"github.com/bitly/go-simplejson"
)

type OskIdentityInstance struct {
	CommonParam RequestParam
}

type OskIdentityResponse struct {
	Token string
	Err   error
}

type OskIdentityInterface interface {
	GetToken() (resp *OskIdentityResponse)
}

func (this OskIdentityInstance) GetToken() (result *OskIdentityResponse) {
	result = &OskIdentityResponse{}
	if this.CommonParam.Attr == nil {
		return result
	}

	openstackConfigJson, err := simplejson.NewJson(this.CommonParam.Attr["openstackConfig"].([]byte))
	adminUser, _ := openstackConfigJson.Get("AdminUser").String()
	adminPass, _ := openstackConfigJson.Get("AdminPass").String()
	adminDomain, _ := openstackConfigJson.Get("AdminDomain").String()
	adminProject, _ := openstackConfigJson.Get("AdminProject").String()
	getTokenUrl, _ := openstackConfigJson.Get("GetTokenUrl").String()

	reqJson := map[string]interface{}{
		"auth": map[string]interface{}{
			"scope": map[string]interface{}{
				"project": map[string]interface{}{
					"domain": map[string]interface{}{
						"name": adminDomain,
					},
					"name": adminProject,
				},
			},
			"identity": map[string]interface{}{
				"methods": []string{
					"password",
				},
				"password": map[string]interface{}{
					"user": map[string]interface{}{
						"name": adminUser,
						"domain": map[string]interface{}{
							"name": adminDomain,
						},
						"password": adminPass,
					},
				},
			},
		},
	}

	str, err := json.Marshal(reqJson)
	this.CommonParam.Url = getTokenUrl
	this.CommonParam.Params = string(str)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}

	this.CommonParam.HttpAction = POST
	resp, err := this.CommonParam.DoPostIdentity()
	result.Err = err
	if err != nil {
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Println("OskIdentity get Token response state = ", resp.StatusCode)
		return
	}
	result.Token = resp.Header.Get("X-Subject-Token")
	return result
}
