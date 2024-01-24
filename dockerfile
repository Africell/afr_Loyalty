FROM alpine:latest
LABEL maintainer="Maher Sidawy <msidawy@africell.com>"
LABEL decsription="Okapi Outlet service"
RUN mkdir /OkapiOutlet
COPY /afr_OkapiOutlet_d /afr_OkapiOutlet_d
#WORKDIR /brmgw
EXPOSE 9290/tcp
EXPOSE 9291/tcp
ENTRYPOINT ["/afr_OkapiOutlet_d"]
CMD [ "/afr_OkapiOutlet_d" ]
