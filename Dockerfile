# Helper image to run podman 4 daemon in rootless mode for testing
FROM archlinux

RUN pacman -Sy --noconfirm podman podman-dnsname netavark aardvark-dns fuse-overlayfs && \
    useradd -m podman && \
    echo podman:10000:5000 > /etc/subuid && \
    echo podman:10000:5000 > /etc/subgid

USER podman
ENTRYPOINT [ "/usr/bin/podman" ]
CMD [ "system", "service", "--time=0", "tcp://:10888" ]