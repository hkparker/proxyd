package main

type TLSConfig map[string]string

type TTPDServiceConfig struct {
	Front		string
	Back		string
	FrontConfig	TLSConfig
	BackConfig	TLSConfig
}

type TTPDServicePack []TTPDServiceConfig

func (service_pack TTPDServicePack) RunServices() {
	listener_failed := make(chan error)
	for _, service_config := range(service_pack) {
		service_up[service_config.Front] = true
		service_up[service_config.Back] = true
		go ListenAndProxy(service_config, listener_failed)
	}
	LogServicesStarted()
	for _ = range(service_pack) {
		<-listener_failed
	}
}

func ListenAndProxy(config TTPDServiceConfig, failed chan error) {
		listener, err := ListenEither(config.Front, config.FrontConfig)
		if err != nil {
			// LogListenFailed()
			failed <- err
			return
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				// call to accept failed...
			} else {
				go ProxyBack(conn, config.Back, config.BackConfig)
			}
		}
}
