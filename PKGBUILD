pkgname=smtd
pkglongname="send-me-the-data"
pkgver=0.1
pkgrel=1
pkgdesc='Send me the Data'
arch=(x86_64)
url="https://github.com/foxpy/$pkglongname"
license=('MIT')
depends=(glibc git)
makedepends=(go)
source=(
    "git+$url.git"
)
sha256sums=(
    SKIP
)

prepare() {
  cd "$pkglongname"
  GOFLAGS="-mod=readonly" go mod vendor -v
}

build() {
  cd "$pkglongname"
  export CGO_LDFLAGS=${LDFLAGS}
  export CGO_CPPFLAGS=${CPPFLAGS}
  export CGO_CFLAGS=${CFLAGS}
  export CGO_CXXFLAGS=${CXXFLAGS}
  export GOFLAGS="-buildmode=pie -mod=vendor -modcacherw -buildvcs=false"
  export GOPATH="$srcdir"

  local ld_flags=" \
    -compressdwarf=false \
    -linkmode=external \
  "
  go build -v -ldflags "$ld_flags" -trimpath -o build/smtd ./cmd/server
}

check() {
  cd "$pkglongname"
  go test -short ./...
}

package() {
  cd "$pkglongname"
  install -vDm 755 -t "${pkgdir}/usr/bin" build/smtd
  install -vDm 600 -t "${pkgdir}/etc" install/smtd.conf
  install -vDm 644 -t "${pkgdir}/usr/lib/systemd/system" install/smtd.service
  install -vDm 644 install/smtd.tmpfiles "${pkgdir}/usr/lib/tmpfiles.d/smtd.conf"
  install -vDm 644 install/smtd.sysusers "${pkgdir}/usr/lib/sysusers.d/smtd.conf"
}
