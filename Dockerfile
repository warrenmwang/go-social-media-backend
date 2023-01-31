# dockerizing the backend API

# base image
FROM debian:stretch-slim

# set the port to run on
ENV PORT 8080

# copy binary over
COPY goserver /bin/goserver

# run
CMD ["/bin/goserver"]