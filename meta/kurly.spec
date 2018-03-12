Name:           kurly
Version:        1.2.1
Release:        0
Summary:        alternative to the widely popular curl program
License:        Apache-2.0
Group:          Applications/Internet
Url:            https://github.com/davidjpeacock/kurly
Source:         https://github.com/davidjpeacock/kurly/archive/v1.2.1.tar.gz
BuildRequires:  go, git
BuildRoot:      %{_tmppath}/%{name}-1.2.1-build

%global debug_package %{nil}

%description
kurly is designed to operate in a similar manner to curl, with select features. Notably, kurly is not aiming for feature parity, but common flags and mechanisms particularly within the HTTP(S) realm are to be expected.

%prep
%setup -q

%build
export CGO_ENABLED=0
export GOPATH=/tmp/gopath/
go get -v -d ./...
go build -a -ldflags "-s -w -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" -o kurly

%install
install -D kurly $RPM_BUILD_ROOT/usr/bin/kurly
gzip meta/kurly.man
install -D -m 0644 meta/kurly.man.gz $RPM_BUILD_ROOT/usr/share/man/man1/kurly.1.gz
install -D -m 0644 LICENSE $RPM_BUILD_ROOT/usr/share/licenses/kurly/LICENCE

%post
%postun

%files
%{_bindir}/kurly
/usr/share/man/man1/kurly.1.gz
/usr/share/licenses/kurly/LICENCE
%defattr(-,root,root)

%changelog
* Mon Mar 12 2018 David J Peacock <david.j.peacock@gmail.com> 1.2.1
- Improved verbosity
- TLS Verbosity
- Support for insecure HTTPS
- Added man page
- Behind-the-scenes refactor for future maintenance
- Handle multiple URLs
- Snap installation via desktop UI

* Fri Dec 29 2017 David J Peacock <david.j.peacock@gmail.com> 1.1.0
- Resume transfer from offset
- Cookie and cookie jar support
- Send HTTP multipart post data
