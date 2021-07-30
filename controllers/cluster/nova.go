package cluster

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getNovaDesiredState(profile *unstructured.Unstructured) map[string]interface{} {

	conf := make(map[string]interface{})

	profileSpec := profile.Object["spec"].(map[string]interface{})

	if profileSpec["novaApiPasteConf"] != nil {
		conf["paste"] = profileSpec["novaApiPasteConf"]
	}
	if profileSpec["novaPolicyConf"] != nil {
		conf["policy"] = profileSpec["novaPolicyConf"]
	}
	if profileSpec["novaConf"] != nil {
		conf["nova"] = profileSpec["novaConf"]
	}
	if profileSpec["novaLoggingConf"] != nil {
		conf["logging"] = profileSpec["novaLoggingConf"]
	}
	if profileSpec["novaApiAuditMapConf"] != nil {
		conf["api_audit_map"] = profileSpec["novaApiAuditMapConf"]
	}

	config := make(map[string]interface{})
	config["conf"] = conf

	componentConfig := make(map[string]interface{})
	componentConfig["nova"] = config

	return componentConfig
}
