package zabbix

import (
	"github.com/AlekSi/reflector"
)

// https://www.zabbix.com/documentation/2.2/manual/api/reference/template/object
type Template struct {
	TemplateId  string `json:"templateid,omitempty"`
	Host        string `json:"host"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

type Templates []Template

type TemplateId struct {
	TemplateId string `json:"templateid"`
}

type TemplateIds []TemplateId

// Wrapper for template.get: https://www.zabbix.com/documentation/2.2/manual/api/reference/template/get
func (api *API) TemplatesGet(params Params) (res Templates, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("template.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}
