FROM alpine:latest
LABEL maintainer="Maher Sidawy <msidawy@africell.com>"
LABEL decsription="Lendme service"
RUN mkdir /Lendme
COPY /afr_Lendme_d /afr_Lendme_d
#WORKDIR /brmgw
EXPOSE 9290/tcp
EXPOSE 9291/tcp
ENTRYPOINT ["/afr_Lendme_d"]
CMD [ "/afr_Lendme_d" ]
