# Introduction

The guide help you to provide step by step instruction on setting up certificate exporter.

# Instructions

### 1. Setup the uaa user.

In order to connect and access the Ops Manager API for extracting all the certificates a minimum of readonly uaa username 
and password must be provided, i.e the username must have minimum `opsman.restricted_view` privilege.

Here is a quick step to create a read only API user.

For more information on creating users on Ops Manager checkout the [documentation](https://docs.pivotal.io/pivotalcf/2-6/customizing/opsman-users.html) and check 
this [documentation](https://docs.pivotal.io/pivotalcf/2-6/opsguide/config-rbac.html) for different role based access control available.

```
# Connect to ops manager UAA using the admin account
uaac target https://<OPSMAN URL>/uaa --skip-ssl-validation
uaac token owner get opsman admin -s "" -p <PASSWORD>

# Create a new read only user and assign a read permission
uaac user add prometheus-cert-exporter -p prometheus-cert-exporter-password --emails prometheus-cert-exporter@prometheus.com
uaac member add opsman.restricted_view prometheus-cert-exporter
```

### 2. Manifest and CF Push

Push the code to cloud foundry

```
# Clone the repository
git clone https://github.com/pivotal-gss/tanzu-certificate-exporter.git
cd tanzu-certificate-exporter

# Open and edit the manifest and provide the values of the env variable.
vi manifest

# push the app to cloud foundry
cf push
```

### 3. Register the route with prometheus 

###### Using [prometheus-boshrelease](https://github.com/bosh-prometheus/prometheus-boshrelease) release

if you are using prometheus which is part of the [prometheus-boshrelease](https://github.com/bosh-prometheus/prometheus-boshrelease) then in order to register the route.

```
# Open the manifests/prometheus.yml file
vi manifests/prometheus.yml

# Add in additional jobs under scrape config i.e under the section jobs > properties > prometheus > scrape_configs

# Say my cert exporter route is "cert-exporter.domain.com", my basic scrape config would be something like this

scrape_configs:
- file_sd_configs:
  - files:
    - /var/vcap/store/bosh_exporter/bosh_target_groups.json
  job_name: prometheus
  relabel_configs:
  - action: keep
    regex: prometheus\d?
    source_labels:
    - __meta_bosh_job_process_name
  - regex: (.*)
    replacement: ${1}:9090
    source_labels:
    - __address__
    target_label: __address__
- job_name: env10
  static_configs:
  - targets:
    - cert-exporter.domain.com
......

# Save the file and update the deployment
bosh -d prometheus deploy manifests/prometheus.yml --vars-store tmp/deployment-vars.yml
  
if you are using additional operator please don't forget to include them like eg.s below, 
check the prometheus-boshrelease for more information on it

bosh -d prometheus deploy manifests/prometheus.yml \
    --vars-store tmp/deployment-vars.yml \
    -o manifests/operators/monitor-bosh.yml \
    -v bosh_url= \
    -v bosh_username= \
    -v bosh_password= \
    --var-file bosh_ca_cert= \
    -v metrics_environment= \
    -o manifests/operators/monitor-cf.yml \
    -v metron_deployment_name= \
    -v system_domain= \
    -v uaa_clients_cf_exporter_secret= \
    -v uaa_clients_firehose_exporter_secret= \
    -v traffic_controller_external_port= \
    -v skip_ssl_verify=
    
# Ensure before confirming that the deployment is only updating the changes implemented above and nothing else, 
  eg.s below show the deployment is only going to publish our changes.

  instance_groups:
  - name: prometheus2
    jobs:
    - name: prometheus2
      properties:
        prometheus:
          scrape_configs:
+         - job_name: "<redacted>"
+           static_configs:
+           - targets:
+             - "<redacted>"

in case you find many variables being modified cancel the deployment during confirmation and ensure you have included 
all the bosh operator when you deployed this during the first time.

# If deployment had successfully completed, the connect to prometheus GUI and see if you can find in metrics from 
"cert_exporter_cert_expires_in_seconds"
```

###### Using your own prometheus

If you are managing your own prometheus in-house, then follow the below steps

```
# Open the prometheus.yml
vi prometheus.yml

# Add in additional jobs under scrape_configs

....
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['127.0.0.1:9090']

  - job_name: 'env10'
    static_configs:
    - targets: ['cert-exporter.domain.com']
    
# If deployment had successfully completed, the connect to prometheus GUI and see if you can find in metrics from 
"cert_exporter_cert_expires_in_seconds"
```


### 4. Repeat

Repeat the step 1 to 3 for additional foundation you would like to monitor. 

### 5. Register the Grafana dashboard

Once the prometheus scraping is setup, navigate to Grafana UI to setup the dashboard.

+ Click on the + sign on the left nav bar
+ Select import
+ Open and copy the content from [Grafana.Json](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/resources/Grafana.json)
+ Paste the Json onto the Grafana Import Page and click on Load
+ Correct any error if found, once satisfied click on import.