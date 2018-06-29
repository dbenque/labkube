FROM busybox

RUN ["touch", "/ready"]

COPY ./labkube /labkube
ENTRYPOINT ["/labkube"]
EXPOSE 8080

