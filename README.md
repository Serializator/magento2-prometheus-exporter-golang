# Magento 2 Prometheus Exporter
A Prometheus exporter, written in Golang, for Magento 2.

## Philosophy
It might be abnormal to start with the "philosophy" but I do believe it lays the groundwork for the "**why**" of another Prometheus exporter when one already exists or you might be considering writing your own.

Magento 2 in and of itself already is a beast of a **monolith**, even without additions made by merchants (e.g. modules). **Why** should introducing observability (using Prometheus) be another one of these additions which complicate the Magento 2 monolith even further?

"Magento 2 Prometheus Exporter" is  a Prometheus exporter which uses the Magento API to construct metrics which it then exports / exposes towards Prometheus, rather than extending Magento 2 itself.

## Contradiction
So... I explained the reasoning behind writing another Prometheus exporter for Magento 2 but what I didn't say is that there exists a **module** for Magento 2 which needs to be installed on the environment you are monitoring. Why talk all evil about "extending Magento 2" and "making a monolith even more complicated" and then still do it yourself?

- Magento 2 itself does not expose enough information for proper metrics
- APIs exposed by Magento 2 sometimes return and do more than we need

**Why not write the exporter in Magento 2 itself then?** Simple, this is because I believe strongly in the "Single Responsibility Principle" (whenever this is reasonable). **Reasonable?** Whenever you work with a third-party platform (e.g. Magento 2) it might not be reasonable to expect, require or criticize for the Single Responsibility Principle not being used.

- It is not possible to independently scale Magento 2 itself and the exporter when it is a module
- Magento 2 by itself (using the traditional `webapi.xml`) does not support output as text (only JSON and XML)
- If you were to extend it to allow for output as text you wouldn't be able to use integrations for access control

It would obviously be possible to extend Magento 2 in such a way to allow for output as text through the "Web API" and still make use of integrations (and thus access control) but then you would not be writing an exporter anymore. It would be another module trying to get things done within the boundaries and of Magento 2 with a negative impact on quality, stability, upgradability of the platform.

## Running It
"Magento 2 Prometheus Exporter" is the exporter itself, written in Golang, which can be scraped by Prometheus. It is meant to be ran as a process alongside your Magento 2 environment (whether it is on-premise or in the cloud).

Whatever Magento 2 environment the exporter is scraping metrics from needs to have the below module installed.
https://github.com/Serializator/magento2-module-prometheus-exporter

In a Kubernetes environment this might mean introducing th "Sidecar Pattern" or when running the exporter on a bare-metal server (on-premise) it might be a process managed by Supervisor.