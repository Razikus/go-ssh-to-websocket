FROM alpine:3.19

RUN apk add --no-cache openssh-server

RUN sed -i 's/#LogLevel INFO/LogLevel VERBOSE/' /etc/ssh/sshd_config
RUN echo 'PermitUserEnvironment yes' >> /etc/ssh/sshd_config
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
RUN sed -i 's/#LogLevel INFO/LogLevel VERBOSE/' /etc/ssh/sshd_config
RUN adduser -D admin
RUN echo 'admin:admin123123' | chpasswd
RUN cp /etc/ssh/sshd_config /config

CMD ssh-keygen -A && /usr/sbin/sshd -f /config -D -e