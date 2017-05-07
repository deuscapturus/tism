Name:		tism		
Version:	0.0
Release:	1%{?dist}
Summary:	the Immutable Secrets Manager.  Encryption-as-a-service.  Encrypt and Decrypt with PGP/GPG

License:	Apache-2.0
URL:		https://github.com/deuscapturus/tism

Source0: https://github.com/deuscapturus/tism/archive/%{version}.tar.gz
BuildRequires:	golang
%define debug_package %{nil}

%description


%prep
rm -rf %{buildroot}/*
%autosetup
mkdir -p %{_builddir}/go/src/github.com/deuscapturus
ln -s %{_builddir}/%{name}-%{version} %{_builddir}/go/src/github.com/deuscapturus/%{name}

%build
unset GOPATH
echo %{_builddir}
export GOPATH=%{_builddir}/go
cd %{_builddir}/go/src/github.com/deuscapturus/%{name}
go get
go build

%install
#Directories
install -d -m 0755 %{buildroot}%{_bindir}
install -d -m 0755 %{buildroot}%{_sysconfdir}/%{name}
install -d -m 0755 %{buildroot}%{_sysconfdir}/systemd/system
install -d -m 0755 %{buildroot}%{_sysconfdir}/pki/tls/certs/%{name}
install -d -m 0755 %{buildroot}/usr/share/%{name}
install -d -m 0755 %{buildroot}%{_sharedstatedir}/%{name}

#Files
install -p -m 0755 %{name} %{buildroot}/%{_bindir}
install -p -m 0755 config.yaml %{buildroot}%{_sysconfdir}/%{name}/
cp -R -p client/* %{buildroot}/usr/share/%{name}/
install -p -m 0755 %{_builddir}/%{name}.service %{buildroot}%{_sysconfdir}/systemd/system/

%files
%{_bindir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/config.yaml
/usr/share/%{name}/* 
%{_sharedstatedir}/%{name}
%{_sysconfdir}/pki/tls/certs/%{name}
%{_sysconfdir}/systemd/system/%{name}.service

%changelog
