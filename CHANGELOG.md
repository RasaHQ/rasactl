
<a name="1.0.4"></a>
## [1.0.4](https://github.com/RasaHQ/rasactl/compare/1.0.3...1.0.4) (2022-02-07)

### Bug Fixes

* Return valid deployment status if helm release is different than the default ([#40](https://github.com/RasaHQ/rasactl/issues/40))
* Set default helm chart version to 4.3.3 ([#41](https://github.com/RasaHQ/rasactl/issues/41))

### Code Refactoring

* Use xerrors package ([#38](https://github.com/RasaHQ/rasactl/issues/38))


<a name="1.0.3"></a>
## [1.0.3](https://github.com/RasaHQ/rasactl/compare/1.0.2...1.0.3) (2022-01-27)

### Bug Fixes

* Support for using different helm release name and versions ([#31](https://github.com/RasaHQ/rasactl/issues/31))
* Return a URL if a web browser is not detected ([#29](https://github.com/RasaHQ/rasactl/issues/29))


<a name="1.0.2"></a>
## [1.0.2](https://github.com/RasaHQ/rasactl/compare/1.0.1...1.0.2) (2022-01-25)

### Bug Fixes

* Use hostAliases for rasa-x pod instead of extraHosts for docker ([#27](https://github.com/RasaHQ/rasactl/issues/27))
* Use local RASA X URL if it connects to local Rasa OSS ([#26](https://github.com/RasaHQ/rasactl/issues/26))


<a name="1.0.1"></a>
## [1.0.1](https://github.com/RasaHQ/rasactl/compare/1.0.0...1.0.1) (2022-01-22)

### Bug Fixes

* Use correct Helm value path for rabbitmq credentials ([#23](https://github.com/RasaHQ/rasactl/issues/23))


<a name="1.0.0"></a>
## [1.0.0](https://github.com/RasaHQ/rasactl/compare/0.5.0...1.0.0) (2022-01-19)

### Breaking changes

* Update default helm values ([#17](https://github.com/RasaHQ/rasactl/issues/17))

### Bug Fixes

* Update dependencies ([#18](https://github.com/RasaHQ/rasactl/issues/18))


<a name="0.5.0"></a>
## [0.5.0](https://github.com/RasaHQ/rasactl/compare/0.4.3...0.5.0) (2022-01-05)

### Features

* Add additional deployment statuses (upgrading & installing) ([#16](https://github.com/RasaHQ/rasactl/issues/16))


<a name="0.4.3"></a>
## [0.4.3](https://github.com/RasaHQ/rasactl/compare/0.4.2...0.4.3) (2021-12-02)

### Bug Fixes

* Check if a namespace is managed by rasactl on delete ([#15](https://github.com/RasaHQ/rasactl/issues/15))


<a name="0.4.2"></a>
## [0.4.2](https://github.com/RasaHQ/rasactl/compare/0.4.1...0.4.2) (2021-11-30)


<a name="0.4.1"></a>
## [0.4.1](https://github.com/RasaHQ/rasactl/compare/0.4.0...0.4.1) (2021-11-18)

### Bug Fixes

* Fix adding a label to namespace if the labels field is absent ([#11](https://github.com/RasaHQ/rasactl/issues/11))


<a name="0.4.0"></a>
## [0.4.0](https://github.com/RasaHQ/rasactl/compare/0.3.0...0.4.0) (2021-10-18)

### Features

* Enable template for the values file ([#9](https://github.com/RasaHQ/rasactl/issues/9))
* Set Rasa X URL via environment variables  ([#8](https://github.com/RasaHQ/rasactl/issues/8))


<a name="0.3.0"></a>
## [0.3.0](https://github.com/RasaHQ/rasactl/compare/0.2.0...0.3.0) (2021-10-15)

### Features

* Add an option to pass creds via env variables ([#7](https://github.com/RasaHQ/rasactl/issues/7))


<a name="0.2.0"></a>
## [0.2.0](https://github.com/RasaHQ/rasactl/compare/0.2.0-rc.2...0.2.0) (2021-10-14)


<a name="0.2.0-rc.2"></a>
## [0.2.0-rc.2](https://github.com/RasaHQ/rasactl/compare/0.2.0-rc.1...0.2.0-rc.2) (2021-10-12)


<a name="0.2.0-rc.1"></a>
## [0.2.0-rc.1](https://github.com/RasaHQ/rasactl/compare/0.1.0...0.2.0-rc.1) (2021-10-07)


<a name="0.1.0"></a>
## [0.1.0](https://github.com/RasaHQ/rasactl/compare/0.0.26...0.1.0) (2021-09-23)


<a name="0.0.26"></a>
## [0.0.26](https://github.com/RasaHQ/rasactl/compare/0.0.25...0.0.26) (2021-09-22)


<a name="0.0.25"></a>
## [0.0.25](https://github.com/RasaHQ/rasactl/compare/0.0.24...0.0.25) (2021-09-21)


<a name="0.0.24"></a>
## [0.0.24](https://github.com/RasaHQ/rasactl/compare/0.0.23...0.0.24) (2021-09-20)


<a name="0.0.23"></a>
## [0.0.23](https://github.com/RasaHQ/rasactl/compare/0.0.22...0.0.23) (2021-09-16)


<a name="0.0.22"></a>
## [0.0.22](https://github.com/RasaHQ/rasactl/compare/0.0.21...0.0.22) (2021-09-16)


<a name="0.0.21"></a>
## [0.0.21](https://github.com/RasaHQ/rasactl/compare/0.0.20...0.0.21) (2021-09-08)


<a name="0.0.20"></a>
## [0.0.20](https://github.com/RasaHQ/rasactl/compare/0.0.19...0.0.20) (2021-09-08)


<a name="0.0.19"></a>
## [0.0.19](https://github.com/RasaHQ/rasactl/compare/0.0.18...0.0.19) (2021-09-07)


<a name="0.0.18"></a>
## [0.0.18](https://github.com/RasaHQ/rasactl/compare/0.0.17...0.0.18) (2021-09-07)


<a name="0.0.17"></a>
## [0.0.17](https://github.com/RasaHQ/rasactl/compare/0.0.16...0.0.17) (2021-09-06)


<a name="0.0.16"></a>
## [0.0.16](https://github.com/RasaHQ/rasactl/compare/0.0.15...0.0.16) (2021-09-06)


<a name="0.0.15"></a>
## [0.0.15](https://github.com/RasaHQ/rasactl/compare/0.0.14...0.0.15) (2021-09-06)


<a name="0.0.14"></a>
## [0.0.14](https://github.com/RasaHQ/rasactl/compare/0.0.13...0.0.14) (2021-09-03)


<a name="0.0.13"></a>
## [0.0.13](https://github.com/RasaHQ/rasactl/compare/0.0.12...0.0.13) (2021-09-02)


<a name="0.0.12"></a>
## [0.0.12](https://github.com/RasaHQ/rasactl/compare/0.0.11...0.0.12) (2021-09-02)


<a name="0.0.11"></a>
## [0.0.11](https://github.com/RasaHQ/rasactl/compare/0.0.10...0.0.11) (2021-08-25)


<a name="0.0.10"></a>
## [0.0.10](https://github.com/RasaHQ/rasactl/compare/0.0.9...0.0.10) (2021-08-25)


<a name="0.0.9"></a>
## [0.0.9](https://github.com/RasaHQ/rasactl/compare/0.0.8...0.0.9) (2021-08-23)


<a name="0.0.8"></a>
## [0.0.8](https://github.com/RasaHQ/rasactl/compare/0.0.7...0.0.8) (2021-08-23)


<a name="0.0.7"></a>
## [0.0.7](https://github.com/RasaHQ/rasactl/compare/0.0.6...0.0.7) (2021-08-20)


<a name="0.0.6"></a>
## [0.0.6](https://github.com/RasaHQ/rasactl/compare/0.0.5...0.0.6) (2021-08-18)


<a name="0.0.5"></a>
## [0.0.5](https://github.com/RasaHQ/rasactl/compare/0.0.4...0.0.5) (2021-08-18)


<a name="0.0.4"></a>
## [0.0.4](https://github.com/RasaHQ/rasactl/compare/0.0.3...0.0.4) (2021-08-17)


<a name="0.0.3"></a>
## 0.0.3 (2021-08-16)

