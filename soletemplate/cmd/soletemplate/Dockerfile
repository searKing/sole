FROM centos:7 as base
FROM base as base-amd64

# [must] maintainer label
LABEL "maintainer"="471030698@qq.com"

# [must] timezone must be Asia/Shanghai (CST)
ENV TZ=Asia/Shanghai
RUN yum -y install kde-l10n-Chinese \
    && yum -y reinstall glibc-common \
    && localedef -c -f UTF-8 -i zh_CN zh_CN.utf8 \
    && echo 'LANG="zh_CN.UTF-8"' >> /etc/locale.conf \
    && source /etc/locale.conf \
    && yum clean all
ENV LANG=zh_CN.UTF-8 \
    LC_ALL=zh_CN.UTF-8

RUN yum groupinstall -y "Development Tools" && yum clean all && rm -rf /var/cache

RUN yum install -y epel-release && yum clean all && rm -rf /var/cache

RUN yum install -y yum-utils \
    telnet tcpdump wget net-tools bind-utils  \
    sysstat lsof lrzsz unzip \
    && yum clean all && rm -rf /var/cache

RUN yum install -y centos-release-scl centos-release-scl-rh && yum clean all && rm -rf /var/cache

# [suggest] working dir is app
ADD . /app
WORKDIR /app

# [suggest] declare exposed port
EXPOSE 80
ENV LD_LIBRARY_PATH /app/lib:$LD_LIBRARY_PATH

# [suggest] use command array
ENTRYPOINT ["/bin/bash","-l","-c"]
CMD ["./bin/sole-template;ls"]