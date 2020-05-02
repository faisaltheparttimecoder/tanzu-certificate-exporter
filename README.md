# Introduction

The exporter is highly inspired from [cert-exporter by joe-elliott](https://github.com/joe-elliott/cert-exporter), On this specific exporter we take a look at visualization of certificates with respect to cloud foundry.

This exporter uses the [vmware tanzu operations manager API](https://docs.pivotal.io/platform/2-8/security/pcf-infrastructure/managing-certificates.html) to get the cert information

# Installation / Usage

We have created a dedicated doc on how to setup and install the exporter, please check out the [set of instruction](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/Install.md) provided on the doc

# Dashboard

After running cert-exporter in your cluster it's easy to build a [custom dashboard](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/resources/Grafana.json) to expose information about the certs in your cluster. Follow the [guide](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/Install.md) on how to set it up.

-- Main Dashboard

![home](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/resources/Dash1.png)

-- Expanded Table

![table](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/resources/Dash2.png)

# Exported Metrics

cert-exporter exports the following metrics

```
# HELP cert_exporter_cert_expires_in_seconds Number of seconds til the cert expires.
# TYPE cert_exporter_cert_expires_in_seconds gauge
cert_exporter_cert_expires_in_seconds{configurable="false",env="env10",is_ca="false",issuer="",location="credhub",product_guid="cf-9c69d13d0df4b67292a9",property_reference="",valid_from="0001-01-01 00:00:00 +0000 UTC",valid_until="2021-04-29 10:38:04 +0000 UTC",variable_path="/p-bosh/cf-9c69d13d0df4b67292a9/diego-instance-identity-leaf-maestro"} 3.125483989441466e+07
# HELP cert_exporter_error_total Cert Exporter Errors
# TYPE cert_exporter_error_total counter
cert_exporter_error_total 0
```

where 
 
+ **cert_exporter_error_total**

  The total number of unexpected errors encountered by cert-exporter. A good metric to watch to feel comfortable certs are being exported properly.
  
+ **cert_exporter_cert_expires_in_seconds**
  
  The number of seconds until a certificate stored in the PEM format is expired. The property reference, path and issuer label indicates the exported cert.
  
# Customizing / Developing

If you wish to customize the code, follow the instruction as per the [doc](https://github.com/pivotal-gss/tanzu-certificate-exporter/blob/master/LocalSetup.md)