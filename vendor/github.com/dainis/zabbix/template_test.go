package zabbix_test

import (
	. "."
	"testing"
)

func TestTemplates(t *testing.T) {
	api := getAPI(t)

	templates, err := api.TemplatesGet(Params{})
	if err != nil {
		t.Fatal(err)
	}

	if len(templates) == 0 {
		t.Fatal("No templates were obtained")
	}
}
