package provider

import (
	"errors"
	"fmt"
	"github.com/dainis/zabbix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

var HOST_INTERFACE_TYPES = map[string]zabbix.InterfaceType{
	"agent": 1,
	"snmp":  2,
	"ipmi":  3,
	"jmx":   4,
}

var interfaceSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"dns": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"ip": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"main": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
			ForceNew: true,
		},
		"port": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "10050",
			ForceNew: true,
		},
		"type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "agent",
			ForceNew: true,
		},
		"interface_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	},
}

func resourceZabbixHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixHostCreate,
		Read:   resourceZabbixHostRead,
		Update: resourceZabbixHostUpdate,
		Delete: resourceZabbixHostDelete,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Technical name of the host.",
			},
			"host_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "(readonly) ID of the host",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: "Visible name of the host.",
			},
			"monitored": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			//any changes to interface will trigger recreate, zabbix api kinda doesn't
			//work nicely, interface can get linked to various things and replacement
			//simply doesn't work
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     interfaceSchema,
				Required: true,
				ForceNew: true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"templates": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func getInterfaces(d *schema.ResourceData) (zabbix.HostInterfaces, error) {
	interfaceCount := d.Get("interfaces.#").(int)

	interfaces := make(zabbix.HostInterfaces, interfaceCount)

	for i := 0; i < interfaceCount; i++ {
		prefix := fmt.Sprintf("interfaces.%d.", i)

		interfaceType := d.Get(prefix + "type").(string)

		typeId, ok := HOST_INTERFACE_TYPES[interfaceType]

		if !ok {
			return nil, errors.New(fmt.Sprintf("%s isnt valid interface type", interfaceType))
		}

		ip := d.Get(prefix + "ip").(string)
		dns := d.Get(prefix + "dns").(string)

		if ip == "" && dns == "" {
			return nil, errors.New("Atleast one of two dns or ip must be set")
		}

		useip := 1

		if ip == "" {
			useip = 0
		}

		main := 1

		if !d.Get(prefix + "main").(bool) {
			main = 1
		}

		interfaces[i] = zabbix.HostInterface{
			IP:    ip,
			DNS:   dns,
			Main:  main,
			Port:  d.Get(prefix + "port").(string),
			Type:  typeId,
			UseIP: useip,
		}
	}

	return interfaces, nil
}

func getHostGroups(d *schema.ResourceData, api *zabbix.API) (zabbix.HostGroupIds, error) {
	configGroups := d.Get("groups").(*schema.Set)
	setHostGroups := make([]string, configGroups.Len())

	for i, g := range configGroups.List() {
		setHostGroups[i] = g.(string)
	}

	log.Printf("Groups %v\n", setHostGroups)

	groupParams := zabbix.Params{
		"output": "extend",
		"filter": map[string]interface{}{
			"name": setHostGroups,
		},
	}

	groups, err := api.HostGroupsGet(groupParams)

	if err != nil {
		return nil, err
	}

	if len(groups) < configGroups.Len() {
		log.Printf("Not all of the specified groups were found on zabbix server")

		for _, n := range configGroups.List() {
			found := false

			for _, g := range groups {
				if n == g.Name {
					found = true
					break
				}
			}

			if !found {
				return nil, errors.New(fmt.Sprintf("Host group %s doesnt exist in zabbix server", n))
			} else {
				log.Printf("%s exists on zabbix server", n)
			}
		}
	}

	hostGroups := make(zabbix.HostGroupIds, len(groups))

	for i, g := range groups {
		hostGroups[i] = zabbix.HostGroupId{
			GroupId: g.GroupId,
		}
	}

	return hostGroups, nil
}

func getTemplates(d *schema.ResourceData, api *zabbix.API) (zabbix.TemplateIds, error) {
	configTemplates := d.Get("templates").(*schema.Set)
	templateNames := make([]string, configTemplates.Len())

	if configTemplates.Len() == 0 {
		return nil, nil
	}

	for i, g := range configTemplates.List() {
		templateNames[i] = g.(string)
	}

	log.Printf("Templates %v\n", templateNames)

	groupParams := zabbix.Params{
		"output": "extend",
		"filter": map[string]interface{}{
			"host": templateNames,
		},
	}

	templates, err := api.TemplatesGet(groupParams)

	if err != nil {
		return nil, err
	}

	if len(templates) < configTemplates.Len() {
		log.Printf("Not all of the specified templates were found on zabbix server")

		for _, n := range configTemplates.List() {
			found := false

			for _, g := range templates {
				if n == g.Name {
					found = true
					break
				}
			}

			if !found {
				return nil, errors.New(fmt.Sprintf("Template %s doesnt exist in zabbix server", n))
			} else {
				log.Printf("Template %s exists on zabbix server", n)
			}
		}
	}

	hostTemplates := make(zabbix.TemplateIds, len(templates))

	for i, t := range templates {
		hostTemplates[i] = zabbix.TemplateId{
			TemplateId: t.TemplateId,
		}
	}

	return hostTemplates, nil
}

func createHostObj(d *schema.ResourceData, api *zabbix.API) (*zabbix.Host, error) {
	host := zabbix.Host{
		Host:   d.Get("host").(string),
		Name:   d.Get("name").(string),
		Status: 0,
	}

	//0 is monitored, 1 - unmonitored host
	if !d.Get("monitored").(bool) {
		host.Status = 1
	}

	hostGroups, err := getHostGroups(d, api)

	if err != nil {
		return nil, err
	}

	host.GroupIds = hostGroups

	interfaces, err := getInterfaces(d)

	if err != nil {
		return nil, err
	}

	host.Interfaces = interfaces

	templates, err := getTemplates(d, api)

	if err != nil {
		return nil, err
	}

	host.TemplateIds = templates
	return &host, nil
}

func resourceZabbixHostCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	host, err := createHostObj(d, api)

	if err != nil {
		return err
	}

	hosts := zabbix.Hosts{*host}

	err = api.HostsCreate(hosts)

	if err != nil {
		return err
	}

	log.Printf("Created host id is %s", hosts[0].HostId)

	d.Set("host_id", hosts[0].HostId)
	d.SetId(hosts[0].HostId)

	return nil
}

func resourceZabbixHostRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	log.Printf("Will read host with id %s", d.Get("host_id").(string))

	host, err := api.HostGetById(d.Get("host_id").(string))

	if err != nil {
		return err
	}

	log.Printf("Host name is %s", host.Name)

	d.Set("host", host.Host)
	d.Set("name", host.Name)

	d.Set("monitored", host.Status == 0)

	params := zabbix.Params{
		"output": "extend",
		"hostids": []string{
			d.Id(),
		},
	}

	templates, err := api.TemplatesGet(params)

	if err != nil {
		return err
	}

	templateNames := make([]string, len(templates))

	for i, t := range templates {
		log.Printf("%d i", i)
		templateNames[i] = t.Host
	}

	d.Set("templates", templateNames)

	groups, err := api.HostGroupsGet(params)

	if err != nil {
		return err
	}

	groupNames := make([]string, len(groups))

	for i, g := range groups {
		groupNames[i] = g.Name
	}

	d.Set("groups", groupNames)

	return nil
}

func resourceZabbixHostUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	host, err := createHostObj(d, api)

	if err != nil {
		return err
	}

	host.HostId = d.Id()

	//interfaces can't be updated, changes will trigger recreate
	//sending previous values will also fail the update
	host.Interfaces = nil

	hosts := zabbix.Hosts{*host}

	err = api.HostsUpdate(hosts)

	if err != nil {
		return err
	}

	log.Printf("Created host id is %s", hosts[0].HostId)

	return nil
}

func resourceZabbixHostDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return api.HostsDeleteByIds([]string{d.Id()})
}
