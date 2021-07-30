package cluster

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getNeutronDesiredState(profile *unstructured.Unstructured) map[string]interface{} {

	conf := make(map[string]interface{})

	profileSpec := profile.Object["spec"].(map[string]interface{})

	if profileSpec["neutronApiPasteConf"] != nil {
		conf["paste"] = profileSpec["neutronApiPasteConf"]
	}
	if profileSpec["neutronPolicyConf"] != nil {
		conf["policy"] = profileSpec["neutronPolicyConf"]
	}
	if profileSpec["neutronConf"] != nil {
		conf["neutron"] = profileSpec["neutronConf"]
	}
	if profileSpec["neutronLoggingConf"] != nil {
		conf["logging"] = profileSpec["neutronLoggingConf"]
	}
	if profileSpec["neutronApiAuditMapConf"] != nil {
		conf["api_audit_map"] = profileSpec["neutronApiAuditMapConf"]
	}
	if profileSpec["neutronDhcpAgentConf"] != nil {
		conf["dhcp_agent"] = profileSpec["neutronDhcpAgentConf"]
	}
	if profileSpec["neutronL3AgentConf"] != nil {
		conf["l3_agent"] = profileSpec["neutronL3AgentConf"]
	}
	if profileSpec["neutronMetadataAgent"] != nil {
		conf["metadata_agent"] = profileSpec["neutronMetadataAgent"]
	}

	plugins := make(map[string]interface{})
	if profileSpec["neutronMl2Conf"] != nil {
		plugins["ml2_conf"] = profileSpec["neutronMl2Conf"]
	}
	if profileSpec["neutronOpenvswitchAgentConf"] != nil {
		plugins["openvswitch_agent"] = profileSpec["neutronOpenvswitchAgentConf"]
	}
	if profileSpec["neutronSriovAgentConf"] != nil {
		plugins["sriov_agent"] = profileSpec["neutronSriovAgentConf"]
	}
	conf["plugins"] = plugins


	config := make(map[string]interface{})
	config["conf"] = conf

	componentConfig := make(map[string]interface{})
	componentConfig["neutron"] = config

	return componentConfig
}
