package Service

type IService interface {

	RegisterProtoHandler()
	BeforeRunThreadHook()
	AfterRunThreadHook()
	Run()
	ObserverPProf(addr string)
}
