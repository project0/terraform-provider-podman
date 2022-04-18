# Helper image to run podman 4 daemon in rootless mode for testing
# Once github actions adds ubuntu 22.04 we may not need it anymore
FROM archlinux

RUN pacman -Sy --noconfirm podman podman-dnsname netavark aardvark-dns fuse-overlayfs && \
    useradd -m podman && \
    echo podman:10000:5000 > /etc/subuid && \
    echo podman:10000:5000 > /etc/subgid

USER podman
ENTRYPOINT [ "/usr/bin/podman" ]
CMD [ "system", "service", "--time=0", "tcp://:8888" ]