FROM debian:latest

COPY wishlist/wishlist /wishlist
RUN chmod +x /wishlist

CMD ["/wishlist"]
