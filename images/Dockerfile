FROM registry.ci.openshift.org/openshift/centos:stream9

LABEL maintainer="bbarbach@redhat.com"

ADD openshift-installer-ci-plugin /usr/bin/openshift-installer-ci-plugin

ENTRYPOINT ["/usr/bin/openshift-installer-ci-plugin"]