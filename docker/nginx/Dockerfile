FROM nginx

WORKDIR /

COPY docker/nginx/blogod.conf /etc/nginx/conf.d/default.conf
COPY script/wait-for-it.sh /nginx/script/

CMD ["/nginx/script/wait-for-it.sh", "blogod:8080", "-s", "-t", "60", "--", "nginx", "-g", "daemon off;"]
