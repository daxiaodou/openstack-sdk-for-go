package compute

import (
	"encoding/json"
	"errors"
	"fmt"

	. "git.intra.weibo.com/golang/openstack-sdk-for-go/http"
	"git.intra.weibo.com/golang/openstack-sdk-for-go/util"
	simplejson "github.com/bitly/go-simplejson"
)

type OskComputeInstance struct {
	CommonParam RequestParam
}

type OskCreateResponse struct {
	ServerId string
	Err      error
}

type OskDeleteResponse struct {
	Err error
}

type OskServerDetail struct {
	Ip     string
	Id     string
	Status string
	Err    error
}

/**
 * 获取单个服务器的信息
 */
func (this OskComputeInstance) GetServerDetail() (detail *OskServerDetail) {
	if this.CommonParam.Attr == nil {
		fmt.Print("[GetServerDetail] params is nil")
		return
	}

	this.CommonParam.HttpAction = GET
	fmt.Printf("getserverdetail \n")
	fmt.Printf("%v\n", this.CommonParam)
	fmt.Printf("%v\n", this.CommonParam.Attr)

	openstackConfigJson ,err := simplejson.NewJson(this.CommonParam.Attr["openstackConfig"].([]byte))
	this.CommonParam.Url,_ = openstackConfigJson.Get("OneServerDetailUrl").String()
	this.CommonParam.Url += util.ToString(this.CommonParam.Attr["serverid"])
	fmt.Println("[GetServerDetail] url = ", this.CommonParam.Url)
	detail = &OskServerDetail{}
	resp, err := this.CommonParam.DoGetRequest()
	if err != nil {
		fmt.Println("[GetServerDetail] http request error: ", err)
		detail.Err = err
		return
	}

	js, _ := simplejson.NewJson([]byte(resp))
	detail = this.parseServerDetail(js)
	return detail
}

/**
 * 创建一台机器
 */
func (this OskComputeInstance) CreateServer() (oskResp *OskCreateResponse) {
	oskResp = &OskCreateResponse{}

	if this.CommonParam.Attr == nil {
		fmt.Print("[CreateServer] params is nil")
		return
	}

	this.CommonParam.HttpAction = POST
	reqJson := map[string]interface{}{
		"server": map[string]interface{}{
			"name":      this.CommonParam.Attr["server_name"],
			"flavorRef": this.CommonParam.Attr["flavor_ref"],
			"imageRef":  this.CommonParam.Attr["os"],
			"security_groups": []interface{}{
				map[string]interface{}{
					"name": this.CommonParam.Attr["security_groups"],
				},
			},
			"networks": []interface{}{
				map[string]interface{}{
					"uuid": this.CommonParam.Attr["network"],
				},
			},
			"adminPass": this.CommonParam.Attr["adminPass"],
		},
	}

	if this.CommonParam.Attr["availability_zone"] != "" {
		if _, ok := reqJson["server"].(map[string]interface{}); ok {
			reqJson["server"].(map[string]interface{})["availability_zone"] = this.CommonParam.Attr["availability_zone"]
		}
	}

	str, err := json.Marshal(reqJson)
	if err != nil {
		fmt.Println("[CreateServer] Error encoding JSON")
		return
	}


	this.CommonParam.Params = string(str)

	openstackConfigJson ,err := simplejson.NewJson(this.CommonParam.Attr["openstackConfig"].([]byte))
	this.CommonParam.Url,_ = openstackConfigJson.Get("CreateServerUrl").String()

	fmt.Println("url: ", this.CommonParam.Url, "params: "+this.CommonParam.Params, " Token: ", this.CommonParam.AccessToken)

	resp, err := this.CommonParam.DoPostRequest()
	if err != nil {
		fmt.Println("[CreateServer] http request error: ", err)
		oskResp.Err = err
		return
	}
	fmt.Println("[CreateServer] request reponse= " + resp)
	json, err := simplejson.NewJson([]byte(resp))
	if err != nil {
		fmt.Println("[CreateServer] decode json error= " + resp)
		oskResp.Err = err
		return
	}
	serverId, _ := json.Get("server").Get("id").String()
	oskResp.ServerId = serverId
	fmt.Println("[CreateServer] server id = " + serverId)
	oskResp.Err = err
	if len(serverId) <= 0 {
		oskResp.Err = errors.New("server id empty")
	}
	return oskResp
}

/**
 * 删除一台机器
 */
func (this OskComputeInstance) DeleteServer() (oskResp *OskDeleteResponse) {
	oskResp = &OskDeleteResponse{}
	if this.CommonParam.Attr == nil {
		fmt.Print("[DeleteServer] params is nil")
		return
	}

	this.CommonParam.HttpAction = DELETE
	openstackConfigJson, err := simplejson.NewJson(this.CommonParam.Attr["openstackConfig"].([]byte))
	this.CommonParam.Url, _ = openstackConfigJson.Get("DeleteServerUrl").String()
	this.CommonParam.Url += util.ToString(this.CommonParam.Attr["serverid"])
	fmt.Println("[DeleteServer] url = ", this.CommonParam.Url)
	resp, err := this.CommonParam.DoDeleteRequest()
	if err != nil {
		fmt.Println("[DeleteServer] http request error: ", err)
		oskResp.Err = err
		return
	}

	fmt.Println("[DeleteServer] request reponse= " + resp)
	return oskResp
}

func (this OskComputeInstance) parseServerDetail(js *simplejson.Json) (resp *OskServerDetail) {
	resp = &OskServerDetail{}
	server := js.Get("server")
	resp.Id, _ = server.Get("id").String()
	resp.Status, _ = server.Get("status").String()
	addresses, ok := server.Get("addresses").Map()
	if ok == nil {
		for networkname, v := range addresses {
			fmt.Printf("networkname %v\n", networkname)
			newv := v.([]interface{})
			for _, vv := range newv {
				newvv := vv.(map[string]interface{})
				theip, ok := newvv["addr"].(string)
				if ok {
					resp.Ip = theip
				}
				fmt.Printf("resp ip\n")
				fmt.Printf("%v\n", resp.Ip)
				break
			}
			break
		}
	}
	return resp
}
