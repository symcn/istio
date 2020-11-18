package httpmx

import (
	wasm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/wasm/v3"
	http_conn "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	wasm_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/wasm/v3"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking"
	"istio.io/istio/pilot/pkg/networking/plugin"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/pkg/log"
)

var (
	httpmx_Log = log.RegisterScope("http_mx", "http metadata exchange debugging", 0)
)

// Plugin implements Istio Telemetry HTTP metadata exchange
type Plugin struct{}

// NewPlugin returns an instance of the httpmx plugin
func NewPlugin() plugin.Plugin {
	return Plugin{}
}

func (p Plugin) OnInboundListener(in *plugin.InputParams, mutable *networking.MutableObjects) error {
	if in.Node.Type != model.SidecarProxy {
		// Only care about sidecar.
		return nil
	}
	return buildFilter(in, mutable, false)
}

func (p Plugin) OnInboundFilterChains(in *plugin.InputParams) []networking.FilterChain {
	// todo: dont know what to implemnet here
	return nil
}

func (p Plugin) OnOutboundListener(in *plugin.InputParams, mutable *networking.MutableObjects) error {
	return nil
}

func (p Plugin) OnInboundPassthrough(in *plugin.InputParams, mutable *networking.MutableObjects) error {
	return nil
}

func (p Plugin) OnInboundPassthroughFilterChains(in *plugin.InputParams) []networking.FilterChain {
	return nil
}

func buildFilter(in *plugin.InputParams, mutable *networking.MutableObjects, isPassthrough bool) error {
	for i := range mutable.FilterChains {
		// todo: if wasm enabled or not?
		// switch tpc protocol?
		if in.ListenerProtocol == networking.ListenerProtocolHTTP || mutable.FilterChains[i].ListenerProtocol == networking.ListenerProtocolHTTP {
			// if wasm not enabled: build pure http mx
			if httpMxFilter := buildHttpMxFilter(); httpMxFilter != nil {
				mutable.FilterChains[i].HTTP = append(mutable.FilterChains[i].HTTP, httpMxFilter)
			}
		}
	}

	return nil
}

func buildHttpMxFilter() *http_conn.HttpFilter {
	vmConfig := &wasm_v3.VmConfig{Runtime: "envoy.wasm.runtime.null"}
	pluginis_vm := &wasm_v3.PluginConfig_VmConfig{
		VmConfig: vmConfig,
	}
	//pluginis_vm.isPluginConfig_Vm()

	wasmPluginConfig := &wasm_v3.PluginConfig{Vm: pluginis_vm}
	filterConfigProto := &wasm.Wasm{
		Config: wasmPluginConfig,
	}

	if filterConfigProto == nil {
		return nil
	}

	return &http_conn.HttpFilter{
		Name:       "istio.metadata_exchange",
		ConfigType: &http_conn.HttpFilter_TypedConfig{TypedConfig: util.MessageToAny(filterConfigProto)},
	}
	return nil
}

//type isPluginConfig_vm struct {}
//
//// NewPlugin returns an instance of the httpmx plugin
//func NewisPluginConfig_vm() wasm_v3.PluginConfig_VmConfig {
//	return wasm_v3.PluginConfig_VmConfig{}
//}
//func (vmp *isPluginConfig_vm) isPluginConfig_Vm() {
//}

//type isPluginConfig_Vm interface {
//	isPluginConfig_Vm()
//}
//
//type PluginConfig_VmConfig struct {
//	VmConfig *VmConfig `protobuf:"bytes,3,opt,name=vm_config,json=vmConfig,proto3,oneof"` // TODO: add referential VM configurations.
//}
//
//func (*PluginConfig_VmConfig) isPluginConfig_Vm() {}
