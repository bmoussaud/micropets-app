package org.moussaud.micropets.pets;

public class Backend {

    private String name;
    private String url;


    public Backend() {
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    @Override
    public String toString() {
        return "Backend{" +
                "name='" + name + '\'' +
                ", url='" + url + '\'' +
                '}';
    }
}
