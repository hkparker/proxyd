package main

type TLSConfig map[string]string

type ServiceConfig struct {
	Front       string
	Back        string
	FrontConfig TLSConfig
	BackConfig  TLSConfig
}

type ServicePack []ServiceConfig

func (service_pack ServicePack) RunServices() {
	listener_failed := make(chan error)
	for _, service_config := range service_pack {
		go ListenAndProxy(service_config, listener_failed)
	}
	for _ = range service_pack {
		<-listener_failed
	}
}

func ListenAndProxy(config ServiceConfig, failed chan error) {
	listener, err := ListenEither(config.Front, config.FrontConfig)
	if err != nil {
		failed <- err
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			// call to accept failed
		} else {
			go ProxyBack(conn, config.Back, config.BackConfig)
		}
	}
}
