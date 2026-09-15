FROM alpine:latest
LABEL maintainer="Maher Sidawy <msidawy@africell.com>"
LABEL decsription="Loyalty service"
RUN mkdir /Loyalty
COPY /afr_Loyalty_d /afr_Loyalty_d
#WORKDIR /brmgw
EXPOSE 9280/tcp
EXPOSE 9281/tcp
ENTRYPOINT ["/afr_Loyalty_d"]
CMD [ "/afr_Loyalty_d" ]
