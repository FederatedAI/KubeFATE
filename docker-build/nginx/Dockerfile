FROM ubuntu:18.04

# Docker Build Arguments
ARG RESTY_VERSION="1.17.8.2"
ARG RESTY_LUAROCKS_VERSION="3.3.1"
ARG RESTY_OPENSSL_VERSION="1.1.1g"
ARG RESTY_OPENSSL_PATCH_VERSION="1.1.1f"
ARG RESTY_OPENSSL_URL_BASE="https://www.openssl.org/source"
ARG RESTY_PCRE_VERSION="8.44"
ARG RESTY_J="4"
ARG RESTY_PREFIX="/data/projects/fate/proxy"
ARG RESTY_CONFIG_OPTIONS="\
    --with-luajit \
    --with-http_ssl_module \
    --with-http_v2_module \
    --with-stream \
    --with-stream_ssl_module \
    "

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        gettext-base \
        libgd-dev \
        libgeoip-dev \
        libncurses5-dev \
        libperl-dev \
        libreadline-dev \
        libxslt1-dev \
        make \
        perl \
        unzip \
        zlib1g-dev \
        libssl-dev

RUN cd /tmp \
    && curl -fSL https://openresty.org/download/openresty-${RESTY_VERSION}.tar.gz -o openresty-${RESTY_VERSION}.tar.gz \
    && tar xzf openresty-${RESTY_VERSION}.tar.gz \
    && cd /tmp/openresty-${RESTY_VERSION} \
    && eval ./configure --prefix=${RESTY_PREFIX} -j${RESTY_J} ${RESTY_CONFIG_OPTIONS} \
    && make -j${RESTY_J} \
    && make -j${RESTY_J} install \
    && cd /tmp \
    && rm -rf \
        openssl-${RESTY_OPENSSL_VERSION}.tar.gz openssl-${RESTY_OPENSSL_VERSION} \
        pcre-${RESTY_PCRE_VERSION}.tar.gz pcre-${RESTY_PCRE_VERSION} \
        openresty-${RESTY_VERSION}.tar.gz openresty-${RESTY_VERSION} \
    && apt-get autoremove -y

WORKDIR /data/projects/fate/proxy/
COPY proxy/conf ./nginx/conf
COPY proxy/lua ./nginx/lua

ENV PATH=$PATH:/data/projects/fate/proxy/bin:/data/projects/fate/proxy/nginx/sbin

CMD ["openresty", "-g", "daemon off;"]

STOPSIGNAL SIGQUIT

