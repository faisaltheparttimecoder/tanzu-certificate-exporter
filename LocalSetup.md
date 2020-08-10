# Introduction

This guide helps with set of instruction on how you run this code locally on your laptop and customize the tool

# Instructions

+ Setup your GOPATH
  ```
  export GOPATH=<PATH>
  ```
+ Clone the repository
  ```
  git clone https://github.com/pivotal-gss/tanzu-certificate-exporter.git
  cd tanzu-certificate-exporter
  ```
+ Start your local web server of the code
  ```
  go run *.go -a https://<OPSMAN-URL>/ -u prometheus-cert-exporter -w prometheus-cert-exporter-password -e env10 -i 30 -k -p 8080
  ```
+ [Install docker](https://docs.docker.com/get-docker/) if you haven't done so
+ Create a `prometheus.yml` file which content similar to below
    ```
    # my global config
    global:
      scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
      evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
      # scrape_timeout is set to the global default (10s).
    
    
    # A scrape configuration containing exactly one endpoint to scrape:
    # Here it's Prometheus itself.
    scrape_configs:
      # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
      - job_name: 'prometheus'
        # metrics_path defaults to '/metrics'
        # scheme defaults to 'http'.
        static_configs:
        - targets: ['127.0.0.1:9090']
    
      - job_name: 'vmware-tanzu-cert-exporter'
        metrics_path: '/metrics'
        static_configs:
        - targets: ['<HOST-IP>:8080']
    ```
    **NOTE:** the `<HOST-IP>` should be physical IP of you laptop not `localhost` or `127.0.0.1`
+ Download and run prometheus / grafana docker images

    ```
    docker run -d --name prometheus -p 9090:9090 -v `pwd`/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
    docker run -d --name grafana -p 3000:3000 grafana/grafana
    ```
 
+ Open browser and navigate to link `localhost:9090` to access prometheus UI or `localhost:3000` to access Grafana UI