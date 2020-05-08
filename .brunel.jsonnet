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
                    image: "docker:dind-rootless",
                    workingDir: "/workspace/",
                    entryPoint: "sh",
                    environment: {
                        "DOCKER_HOST": "tcp://docker:2375",
                    },
                    args: [
                        "-c",
                        "--",
                        "docker build -f ./Dockerfile.runner .",
                    ]
                }
            ]
        },
    }
}