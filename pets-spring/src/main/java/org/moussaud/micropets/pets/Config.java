package org.moussaud.micropets.pets;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import java.io.Serializable;
import java.util.List;

@Configuration
@ConfigurationProperties(prefix = "myconfig")
public class Config implements Serializable {

    public Config() {

    }

    public List<Backend> getBackends() {
        return backends;
    }

    public boolean isUseServiceDiscovery() {
        return useServiceDiscovery;
    }

    public void setBackends(List<Backend> backends) {
        this.backends = backends;
    }

    public void setUseServiceDiscovery(boolean useServiceDiscovery) {
        this.useServiceDiscovery = useServiceDiscovery;
    }

    private static final long serialVersionUID = -8900976432592584351L;

    private List<Backend> backends;
    private boolean useServiceDiscovery;

    @Override
    public String toString() {
        return "Config{" +
                "backends=" + backends +
                ", useServiceDiscovery=" + useServiceDiscovery +
                '}';
    }
}