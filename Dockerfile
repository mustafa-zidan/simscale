FROM scratch
ADD bin/simscale_*_linux_amd64 /simscale
CMD ["/simscale"]
