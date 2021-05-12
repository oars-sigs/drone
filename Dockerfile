FROM library/centos:7
ADD bin/drone-server /bin/
ENTRYPOINT ["/bin/drone-server"]