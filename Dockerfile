FROM ubuntu:22.04

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV NODE_VERSION 19.6.0
ENV GO_VERSION 1.19.6

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN ARCH= && dpkgArch="$(dpkg --print-architecture)" \
    && case "${dpkgArch##*-}" in \
      amd64) ARCH='x64';; \
      ppc64el) ARCH='ppc64le';; \
      s390x) ARCH='s390x';; \
      arm64) ARCH='arm64';; \
      armhf) ARCH='armv7l';; \
      i386) ARCH='x86';; \
      *) echo "unsupported architecture"; exit 1 ;; \
    esac \
    && set -ex \
    # libatomic1 for arm
    && apt-get update && apt-get install -y ca-certificates curl wget g++ gcc git libc6-dev make pkg-config gnupg dirmngr xz-utils libatomic1 --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && for key in \
      4ED778F539E3634C779C87C6D7062848A1AB005C \
      141F07595B7B3FFE74309A937405533BE57C7D57 \
      74F12602B6F1C4E913FAA37AD3A89613643B6201 \
      61FC681DFB92A079F1685E77973F295594EC4689 \
      8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
      C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
      890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4 \
      C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C \
      108F52B48DB57BB0CC439B2997B01419BD92F80A \
    ; do \
      gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys "$key" || \
      gpg --batch --keyserver keyserver.ubuntu.com --recv-keys "$key" ; \
    done \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-$ARCH.tar.xz" \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
    && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
    && grep " node-v$NODE_VERSION-linux-$ARCH.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
    && tar -xJf "node-v$NODE_VERSION-linux-$ARCH.tar.xz" -C /usr/local --strip-components=1 --no-same-owner \
    && rm "node-v$NODE_VERSION-linux-$ARCH.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt \
    && apt-mark auto '.*' > /dev/null \
    && find /usr/local -type f -executable -exec ldd '{}' ';' \
      | awk '/=>/ { print $(NF-1) }' \
      | sort -u \
      | xargs -r dpkg-query --search \
      | cut -d: -f1 \
      | sort -u \
      | xargs -r apt-mark manual \
    # && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
    && ln -s /usr/local/bin/node /usr/local/bin/nodejs \
    # smoke tests
    && node --version \
    && npm --version \
    && npx playwright install --with-deps chromium \
	&& url= \
	&& case "${dpkgArch##*-}" in \
		amd64) \
			url="https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz"; \
			sha256='e3410c676ced327aec928303fef11385702a5562fd19d9a1750d5a2979763c3d'; \
			;; \
		arm64) \
			url="https://dl.google.com/go/go$GO_VERSION.linux-arm64.tar.gz"; \
			sha256='e4d63c933a68e5fad07cab9d12c5c1610ce4810832d47c44314c3246f511ac4f'; \
			;; \
		i386) \
			url="https://dl.google.com/go/go$GO_VERSION.linux-386.tar.gz"; \
			sha256='da4546cc516ae88698e8643d6d22fe1465b44dd49fc36abb34d53a93c19581ad'; \
			;; \
		*) echo >&2 "error: unsupported architecture"; exit 1;; \
        esac \
    && wget -O /go.tgz "$url" --no-check-certificate \
    && rm -rf /usr/local/go \ 
    && tar -C /usr/local -xzf /go.tgz \
    && rm /go.tgz
    # && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false


WORKDIR /greypot
COPY . .

RUN cd /greypot/ui \
    && npm install \
    && npm run build \
    && cd /greypot/cmd/greypot-server \
    && go build -o /bin/greypot-server 


ENV PORT 7665
EXPOSE 7665

VOLUME /templates

CMD [ "/bin/greypot-server", "-templates", "/templates" ]