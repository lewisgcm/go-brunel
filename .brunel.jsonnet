{
    version: "v1",
    description: "Brunel CI/CD Build Pipeline",
    stages: {
        runner_build: {
            steps: [
                {
                    image: "gcr.io/kaniko-project/executor:latest",
                    working_dir: "/workspace",
                    entrypoint: "sh",
                    args: [
                        "--dockerfile",
                        "./Dockerfile.runner",
                        "--context",
                        "dir://workspace",
                        "--destination",
                        "lewisgcm/go-brunel:runner"
                    ]
                }
            ]
        },
    }
}