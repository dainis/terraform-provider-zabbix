# Terraform zabbix provider

Allows to manage zabbix hosts

```
provider "zabbix" {
  user = "admin"
  password = "zabbix"
  server_url = "http://localhost/zabbix/api_jsonrpc.php"
}

resource "zabbix_host" "zabbix1" {
  host = "127.0.0.1"
  interfaces = [{
    ip = "127.0.0.1"
    main = true
  }]
  groups = ["Linux servers", "${zabbix_host_group.zabbix.name}"]
  templates = ["Template ICMP Ping"]
}

resource "zabbix_host_group" "zabbix" {
  name = "something"
}
```