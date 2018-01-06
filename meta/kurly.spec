Name:           kurly
Version:        1.1.0
Release:        0
Summary:	alternative to the widely popular curl program
License:        Apache-2.0
Group:          Applications/Internet
Url:            https://github.com/davidjpeacock/kurly
Source:         https://github.com/davidjpeacock/kurly/archive/v1.1.0.tar.gz
%if 0%{?suse_version}
BuildRequires:	go
%else
BuildRequires:	golang
%endif
BuildRoot:      %{_tmppath}/%{name}-1.1.0-build

%description
kurly is designed to operate in a similar manner to curl, with select features. Notably, kurly is not aiming for feature parity, but common flags and mechanisms particularly within the HTTP(S) realm are to be expected.

%prep
%setup -q

%build
export CGO_ENABLED=0
export GOPATH=/tmp/gopath/
go get -v -d ./...
go build -ldflags "-s -w" -o kurly

%install
install -D kurly $RPM_BUILD_ROOT/usr/bin/kurly

%post
%postun

%files
%{_bindir}/kurly
%defattr(-,root,root)

%changelog
* Fri Dec 29 2017 David J Peacock <david.j.peacock@gmail.com> 1.1.0
- Resume transfer from offset
- Cookie and cookie jar support
- Send HTTP multipart post data
