FROM debian:latest

COPY bin/wishlist /wishlist
RUN chmod +x /wishlist

CMD ["/wishlist"]
