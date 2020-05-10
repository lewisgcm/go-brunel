{
    version: "v1",
    description: "Brunel CI/CD Build Pipeline",
    stages: {
        "Build Runner": {
            services: [
                {
                    image: "docker:dind",
                    privileged: true,
                    wait: {
                        output: "Daemon has completed initialization",
                        timeout: 30
                    },
                    environment: {
                        "DOCKER_TLS_CERTDIR": "",
                    },
                    hostname: "docker"
                },
            ],
            steps: [
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                    },
                    args: [
                        "-c",
                        "--",
                        "docker build -t lewisgcm/go-brunel:runner -f ./Dockerfile.runner .",
                    ]
                },
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                    },
                    args: [
                        "-c",
                        "--",
                        "docker build -t lewisgcm/go-brunel:latest -f ./Dockerfile .",
                    ]
                },
                {
                    image: "docker:dind",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                        "DOCKER_HUB_ACCESS_TOKEN": brunel.secret('DOCKER_HUB_ACCESS_TOKEN')
                    },
                    args: [
                        "-c",
                        "--",
                        "docker login -u lewisgcm -p $DOCKER_HUB_ACCESS_TOKEN && docker push lewisgcm/go-brunel:runner && docker push lewisgcm/go-brunel:latest",
                    ]
                },
            ]
        },
    }
}