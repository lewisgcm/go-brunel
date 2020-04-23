# go-brunel cloud native CI/CD
Brunel is a cloud native CI/CD system built with first class docker and kubernetes support.
The primary focus of brunel was to create an enterprise ready open source build system with the following properties:

* Build configuration as code
* Re-usable and shareable build configuration
* One command to run builds locally
* Limit resources for builds when running in the cloud
* Sidecar service support with readyness detection
* De-centralized and highly available

## Getting Started

Before running brunel you first need to run the following steps:

### 1. Configuring OAuth
* gitlab
    1. Go to the page [https://gitlab.com/profile/applications](https://gitlab.com/profile/applications)
    2. Add a new application called 'brunel' with the redirect url 'http://localhost:8081/api/user/callback?provider=gitlab' where 'http://localhost:8081' is the external
URL for your brunel instance.

### 2. Generating certificates
If you are using docker compose the certs can be generated using the following command:
```bash
docker-compose exec server /opt/brunel/cert
```

Which will produce two certificates `cert` and `key`. Both variables should be set for both the server and each runner.
You can do this via the environment variables `BRUNEL_REMOTE_CREDENTIALS_CERT` and `BRUNEL_REMOTE_CREDENTIALS_KEY`.
The `docker-compose.yaml` file included in this project has a default key and cert, *this should be changed for production*.

### 3. Running
You can use the docker-compose file to quickly spin up a brunel server for local testing.
To do this put your oauth token and secret update the 'BRUNEL_OAUTH_GITLAB_KEY' and 'BRUNEL_OAUTH_GITLAB_SECRET' variables in the 'docker-compose.yaml' file.
You can now run brunel using the following:
```bash
docker-compose up -d
```

Then visit the URL http://localhost:8081.




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
