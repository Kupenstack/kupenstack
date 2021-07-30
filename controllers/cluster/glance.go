package cluster

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getGlanceDesiredState(profile *unstructured.Unstructured) map[string]interface{} {

	conf := make(map[string]interface{})

	profileSpec := profile.Object["spec"].(map[string]interface{})

	if profileSpec["glanceApiConf"] != nil {
		conf["glance"] = profileSpec["glanceApiConf"]
	}
	if profileSpec["glanceLoggingConf"] != nil {
		conf["logging"] = profileSpec["glanceLoggingConf"]
	}
	if profileSpec["glanceApiPasteConf"] != nil {
		conf["paste"] = profileSpec["glanceApiPasteConf"]
	}
	if profileSpec["glanceRegistryConf"] != nil {
		conf["glance_registry"] = profileSpec["glanceRegistryConf"]
	}
	if profileSpec["glanceRegistryPasteConf"] != nil {
		conf["paste_registry"] = profileSpec["glanceRegistryPasteConf"]
	}
	if profileSpec["glancePolicyConf"] != nil {
		conf["policy"] = profileSpec["glancePolicyConf"]
	}
	if profileSpec["glanceApiAuditMapConf"] != nil {
		conf["api_audit_map"] = profileSpec["glanceApiAuditMapConf"]
	}
	if profileSpec["glanceSwiftStoreConf"] != nil {
		conf["swift_store"] = profileSpec["glanceSwiftStoreConf"]
	}

	config := make(map[string]interface{})
	config["conf"] = conf

	componentConfig := make(map[string]interface{})
	componentConfig["glance"] = config

	return componentConfig
}
