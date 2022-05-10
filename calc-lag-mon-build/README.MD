Build calculation lag monitor & deploy it as a Docker container.

* Use script build-container.bat to create an image
* Use script redeploy-container.bat to update a running container
* Set environment variable DOCKER_HOST=ssh://... to build & deploy container to a remote host
    * Both build-container.bat and redeploy-container.bat should run with DOCKER_HOST=ssh://... for the new build to take effect on the remote host
