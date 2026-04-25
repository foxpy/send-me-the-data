pkgname=smtd
pkgver=0.1
pkgrel=1
pkgdesc='Send me the Data'
arch=(x86_64)
url='https://github.com/foxpy/send-me-the-data'
license=('MIT')
depends=(glibc git)
makedepends=(go)
source=(
    git+file://${PWD}
)
sha256sums=(
    SKIP
)

prepare() {
  cd "$pkgname"
  GOFLAGS="-mod=readonly" go mod vendor -v
}

build() {
  cd "$pkgname"
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
  go build -v -ldflags "$ld_flags" -o build/smtd ./cmd/server
}

check() {
  cd "$pkgname"
  go test -short ./...
}

package() {
  cd "$pkgname"
  install -vDm 755 -t "${pkgdir}/usr/bin" build/smtd
  install -vDm 600 -t "${pkgdir}/etc" install/smtd.conf
  install -vDm 644 -t "${pkgdir}/usr/lib/systemd/system" install/smtd.service
  install -vDm 644 install/smtd.tmpfiles "${pkgdir}/usr/lib/tmpfiles.d/smtd.conf"
  install -vDm 644 install/smtd.sysusers "${pkgdir}/usr/lib/sysusers.d/smtd.conf"
}
