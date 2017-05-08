FROM fedora:rawhide
MAINTAINER Theodore Cowan
EXPOSE 8080

RUN rpm -i https://github.com/deuscapturus/tism/releases/download/0.0/tism-0.0-1.fc25.x86_64.rpm
RUN systemctl enable tism
RUN sed -i -e 's/!locked//g' /etc/shadow

ENTRYPOINT ["/usr/bin/tism"]
