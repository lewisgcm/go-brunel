# go-brunel cloud native CI/CD
Brunel is a cloud native CI/CD system built with first class docker and kubernetes support.
The primary focus of brunel was to create an enterprise ready open source build system with the following properties:

* Build configuration as code
* Re-usable and shareable build configuration
* One command to run builds locally
* Limit resources for builds when running in the cloud
* Sidecar service support with readyness detection
* De-centralized and highly available


## Build Syntax
Brunel uses jsonnet as its underlying build format. See below for an example build configuration file.

Example file `.brunel.jsonnet`:

```jsonnet
// We support comments, and dynamic library loading!!
local shared = brunel.shared({
    repository: "/Users/lewis/Documents/Projects/test_repo",
    branch: "master",
    file: "file.jsonnet"
});

{
    version: "v1",
    description: "My-Project Description",
    stages: {
        test: {
            services: [
                {
                    image: "mysql:latest",
                    wait: {
                        output: "Ready for connections",
                        timeout: 10
                    },
                    hostname: "mysql"
                },
                {
                    image: "nginx:latest",
                    hostname: "nginx"
                }
            ],
            steps: [
                {
                    image: "alpine",
                    resources: {
                        limits: {
                            CPU: 0.5,
                            memory: "200m"
                        }
                    },
                    // Notice how we can curl our sidecars using the supplied hostname
                    // This works in both docker, and kubernetes
                    args: [ "-c", "--", "curl http://nginx" ]
                },
                shared.myCompaniesStandardGoTest(),
            ],
        },
        build: shared.myCompaniesStandardGoBuild(),
    }
}
```
